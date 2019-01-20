package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/fsouza/go-dockerclient"

	"github.com/astromechza/md-slides/pkg/slide"
)

const pdfUsage = `Render the slides to pdf.

Usage:
  md-slides pdf [options...] [file.md] [output.pdf]

`

const targetImage = "astromechza/md-slides-pdf-exporter:latest"

func SubcommandPDF(args []string) error {
	fs := flag.NewFlagSet("", flag.ExitOnError)

	noPullFlag := fs.Bool("no-pull", false, "Don't try to pull the Docker image")
	noStaticFlag := fs.Bool("no-statics", false, "Don't serve static files")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, pdfUsage)
		fs.PrintDefaults()
	}
	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() != 2 {
		fs.Usage()
		fmt.Fprintf(os.Stderr, "\n")
		return fmt.Errorf("expected two positional argument")
	}

	sourceFileName, targetFile := fs.Arg(0), fs.Arg(1)
	slideSource := &slide.CachedSource{Inner: &slide.FileSource{Path: sourceFileName}}

	td, err := ioutil.TempDir(os.TempDir(), "md-slides")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory")
	}
	defer os.RemoveAll(td)

	if err := exportSlidesHTML(slideSource, td, *noStaticFlag); err != nil {
		return err
	}

	log.Printf("Attempting to connect to Docker..")
	cli, err := docker.NewClient("unix:///var/run/docker.sock")
	if err != nil {
		return fmt.Errorf("failed to setup client: %s", err)
	}

	if !*noPullFlag {
		log.Printf("Attempting to pull image %s..", targetImage)
		if err := cli.PullImage(docker.PullImageOptions{
			Repository:   targetImage,
			OutputStream: os.Stderr,
		}, docker.AuthConfiguration{}); err != nil {
			return fmt.Errorf("failed to pull image: %s", err)
		}
	}

	container, err := cli.CreateContainer(docker.CreateContainerOptions{
		Config: &docker.Config{
			Image: targetImage,
		},
		HostConfig: &docker.HostConfig{
			Mounts: []docker.HostMount{
				{
					Type:     "bind",
					Source:   td,
					Target:   "/source",
					ReadOnly: true,
				},
			},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to create container: %s", err)
	}
	defer func() {
		log.Printf("Removing container %s", container.ID)
		cli.RemoveContainer(docker.RemoveContainerOptions{
			ID:    container.ID,
			Force: true,
		})
	}()

	log.Printf("Starting container..")
	if err := cli.StartContainer(container.ID, nil); err != nil {
		return fmt.Errorf("failed to start container: %s", err)
	}

	log.Printf("Waiting for container..")
	ctx, _ := context.WithTimeout(context.Background(), time.Second*30)
	status, err := cli.WaitContainerWithContext(container.ID, ctx)
	if err != nil {
		return fmt.Errorf("failed to wait for container: %s", err)
	}
	log.Printf("container exit status: %v", status)

	log.Printf("Attempting to copy out result pdf..")
	f, err := os.Create(targetFile)
	if err != nil {
		return fmt.Errorf("failed to create output file %s: %s", targetFile, err)
	}
	defer f.Close()

	err = cli.DownloadFromContainer(container.ID, docker.DownloadFromContainerOptions{
		OutputStream: f,
		Path:         "/output/index.pdf",
	})
	if err != nil {
		return fmt.Errorf("failed to get output reader: %s", err)
	}

	return nil
}
