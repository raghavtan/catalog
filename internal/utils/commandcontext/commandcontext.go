package commandcontext

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var commandcontext context.Context

func Init() context.Context {
	ctx, cancel := context.WithCancel(context.Background())
	commandcontext = ctx
	HandleSignal(cancel)
	return commandcontext
}

func HandleSignal(cancel context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Println("Received termination signal, shutting down...")
		cancel()
	}()
}
