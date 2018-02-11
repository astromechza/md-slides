package pdfrenderer

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/url"
	"os/exec"
	"syscall"
	"time"

	"github.com/mafredri/cdp"
	"github.com/mafredri/cdp/devtool"
	"github.com/mafredri/cdp/rpcc"
	// import CDP for now
	_ "github.com/mafredri/cdp"
)

type ChromeProcess struct {
	cmd       *exec.Cmd
	port      int
	listening bool
}

func New() (*ChromeProcess, error) {
	binary := DetectPathToChrome()
	if binary == "" {
		return nil, fmt.Errorf("failed to detect chrome binary or app")
	}

	p := &ChromeProcess{}
	p.port = 10000 + rand.Intn(10000)
	p.cmd = exec.CommandContext(context.Background(), binary,
		"--headless",
		"--disable-gpu",
		fmt.Sprintf("--remote-debugging-port=%d", p.port),
	)
	if err := p.cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start chrome: %s", err)
	}
	log.Printf("Spawned chrome pid=%d", p.cmd.Process.Pid)
	return p, nil
}

func (cp *ChromeProcess) WaitForPort(ctx context.Context) error {
	if cp.port <= 0 {
		return fmt.Errorf("cannot wait for port <= 0")
	}
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("wait for port expired")
		default:
		}
		conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", cp.port), time.Second)
		if conn != nil {
			conn.Close()
			cp.listening = true
			return nil
		}
		log.Printf("wait for port error: %s", err)
		time.Sleep(time.Millisecond * 100)
	}
}

func (cp *ChromeProcess) hostPort() string {
	return fmt.Sprintf("127.0.0.1:%d", cp.port)
}

func (cp *ChromeProcess) url() string {
	return (&url.URL{Scheme: "http", Host: cp.hostPort()}).String()
}

func (cp *ChromeProcess) Client() (*cdp.Client, error) {
	if !cp.listening {
		return nil, fmt.Errorf("no idea if chrome is listening, please call WaitForPort")
	}
	devt := devtool.New(cp.url())
	pt, err := devt.Get(context.Background(), devtool.Page)
	if err != nil {
		pt, err = devt.Create(context.Background())
		if err != nil {
			return nil, err
		}
	}

	conn, err := rpcc.DialContext(context.Background(), pt.WebSocketDebuggerURL)
	if err != nil {
		return nil, err
	}
	c := cdp.NewClient(conn)
	c.Page.Enable(context.Background())
	c.LayerTree.Enable(context.Background())
	return c, nil
}

func (cp *ChromeProcess) Kill() error {
	return cp.cmd.Process.Signal(syscall.SIGTERM)
}
