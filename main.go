package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jacobweinstock/ihttp/cmd"
)

func main() {
	exitCode := 0
	defer func() {
		os.Exit(exitCode)
	}()

	ctx, done := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGHUP, syscall.SIGTERM)
	defer done()

	c := cmd.Config{}
	fs := flag.NewFlagSet("ihttp", flag.ExitOnError)
	fs.IntVar(&c.Port, "port", 8080, "Port to listen on")
	fs.StringVar(&c.CFilePath, "config", "config.yaml", "Path to config file")
	fs.StringVar(&c.CBackend, "content", "", "Backend to use for getting the ipxe script content")
	fs.StringVar(&c.LBackend, "locator", "", "Backend to use for getting the ipxe script location")
	if err := fs.Parse(os.Args[1:]); err != nil {
		log.Printf("error parsing flags: %v", err)
		exitCode = 1
		return
	}
	log.Println("Starting server...", c)
	if err := cmd.Run(ctx, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
	}
}
