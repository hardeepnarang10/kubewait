package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/hardeepnarang10/kubewait/internal/types"
)

var services types.ServiceSet
var interval, timeout time.Duration

func init() {
	if err := types.SetNamespace(); err != nil {
		log.Panicln(fmt.Errorf("Unable to infer pod namespace from (default) Kubernetes secret mountpath. Is the process running inside Kubernetes?\n%w", err))
	}

	flag.Var(&services, "service", "Service, which pods to wait for. Can be specified multiple times")
	flag.DurationVar(&interval, "interval", time.Duration(5)*time.Second, "Duration to wait between check attempts")
	flag.DurationVar(&timeout, "timeout", time.Duration(5)*time.Minute, "Total duration to wait for services to be ready")
	flag.Parse()
}
