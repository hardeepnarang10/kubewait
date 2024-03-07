package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/hardeepnarang10/kubewait/client/kubernetes"
)

func main() {
	if len(services.SVCMap) == 0 {
		log.Println("No services specified. Exiting...")
		return
	}

	timeoutCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	ctx, stop := signal.NotifyContext(timeoutCtx, os.Interrupt, os.Kill)
	defer stop()

	cli, err := kubernetes.ClientSet()
	if err != nil {
		log.Panicln(fmt.Errorf("unable to prepare kubernetes clientset with assigned config: %w", err))
	}

	if err := runner(ctx, time.NewTicker(interval), kubernetes.Readiness(cli.CoreV1())); err != nil {
		log.Panicln(fmt.Errorf("runner returned non-ok response status: %w", err))
	}
}
