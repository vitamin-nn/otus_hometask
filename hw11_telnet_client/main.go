package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	ErrInvalidArguments = errors.New("invalid arguments")
	ErrTimeout          = errors.New("timeout error")
)

func main() {
	var timeoutIn string
	flag.StringVar(&timeoutIn, "timeout", "10s", "define timeout (default 10s)")
	flag.Parse()
	args := flag.Args()
	timeout, err := time.ParseDuration(timeoutIn)
	if err != nil || len(args) != 2 {
		log.Fatalln(ErrInvalidArguments)
	}
	connAddr := net.JoinHostPort(args[0], args[1])
	c := NewTelnetClient(connAddr, timeout, os.Stdin, os.Stdout)
	err = c.Connect()
	if err != nil {
		log.Fatalln(ErrTimeout)
	}
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", connAddr)

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		// получилось придумать только такое решение блокировки Scan() в методе Receive:
		// досрочное закрытие коннекта
		defer wg.Done()
		<-ctx.Done()
		c.Close()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		waitExitSignal(ctx, cancel)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		sendRoutine(ctx, cancel, c)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		reciveRoutine(ctx, cancel, c)
	}()

	wg.Wait()
}

func reciveRoutine(ctx context.Context, cancel context.CancelFunc, c TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.Receive()
			if err != nil {
				cancel()
				return
			}
		}
	}
}

func sendRoutine(ctx context.Context, cancel context.CancelFunc, c TelnetClient) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := c.Send()
			if err != nil {
				cancel()
				return
			}
		}
	}
}

func waitExitSignal(ctx context.Context, cancel context.CancelFunc) {
	sch := make(chan os.Signal, 1)
	signal.Notify(sch, os.Interrupt, syscall.SIGINT)
	select {
	case <-sch:
		cancel()
	case <-ctx.Done():
		break
	}
}
