package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mafredri/cdp/protocol/page"

	"github.com/AstromechZA/md-slides/pdfrenderer"
	"github.com/AstromechZA/md-slides/sliderenderer"
)

func PDF(args []string) error {
	fs := flag.NewFlagSet("pdf", flag.ExitOnError)
	resFlag := fs.String("res", "1600x900", "set render aspect ratio and zoom for rendering")
	listenFlag := fs.String("listen", ":8080", "interface:port to listen on")
	backgroundCSS := fs.String("css-background", "#fffff8", "slide background css")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected a single source file as argument")
	}
	filename := fs.Arg(0)

	xres, yres, err := parseResString(*resFlag)
	if err != nil {
		return fmt.Errorf("bad res string: %s", err)
	}
	mux := http.NewServeMux()

	sr := sliderenderer.SlideRenderer{Filename: filename, XRes: xres, YRes: yres, BGCSS: *backgroundCSS}
	if err = sr.CheckSlides(); err != nil {
		return fmt.Errorf("check failed: %s", err)
	}
	sr.InstallMultiSlideHandler(mux)
	server := &http.Server{Addr: *listenFlag, Handler: mux}
	go func() {
		server.ListenAndServe()
	}()
	defer server.Shutdown(context.Background())

	chrome, err := pdfrenderer.New()
	if err != nil {
		return err
	}
	log.Printf("chrome: %#v", chrome)
	if err = chrome.WaitForPort(context.Background()); err != nil {
		return err
	}
	chromeClient, err := chrome.Client()
	if err != nil {
		return err
	}
	log.Printf("chrome client: %#v", chromeClient)

	// Open a DOMContentEventFired client to buffer this event.
	domContent, err := chromeClient.Page.DOMContentEventFired(context.Background())
	if err != nil {
		return err
	}
	defer domContent.Close()

	log.Printf("navigating to thing")
	_, err = chromeClient.Page.Navigate(context.Background(), page.NewNavigateArgs("http://127.0.0.1:8080/_multislide/"))
	if err != nil {
		return err
	}

	// Wait until we have a DOMContentEventFired event.
	if _, err = domContent.Recv(); err != nil {
		return err
	}

	pdfArgs := page.NewPrintToPDFArgs().SetDisplayHeaderFooter(false).SetLandscape(true).SetPrintBackground(true)
	pdfArgs = pdfArgs.SetMarginBottom(0).SetMarginLeft(0).SetMarginRight(0).SetMarginTop(0)
	repl, err := chromeClient.Page.PrintToPDF(context.Background(), pdfArgs)
	if err != nil {
		return err
	}
	f, err := os.Create("output.pdf")
	if err != nil {
		return err
	}
	f.Write(repl.Data)
	f.Close()

	log.Printf("killing chrome")
	if err = chrome.Kill(); err != nil {
		return err
	}
	log.Printf("done")

	return nil
}
