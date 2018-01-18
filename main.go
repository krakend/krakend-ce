package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/devopsfaith/krakend-cobra"
	"github.com/devopsfaith/krakend-viper"
)

func main() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		select {
		case sig := <-sigs:
			log.Println("Signal intercepted:", sig)
			cancel()
		case <-ctx.Done():
		}
	}()

	RegisterEncoders()

	cmd.Execute(viper.New(), NewExecutor(ctx))
}
