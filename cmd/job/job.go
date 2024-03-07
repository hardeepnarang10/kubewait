package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hardeepnarang10/kubewait/internal/types"
)

func runner(ctx context.Context, ticker *time.Ticker, readiness func(context.Context, types.Service) error) error {
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			return fmt.Errorf("context timeout or cancelled before all services reached ready state: %q", services.String())
		case <-ticker.C:
			completed := make([]*types.Service, 0)
			for key := range services.SVCMap {
				if err := readiness(ctx, *key); err != nil {
					log.Println(err)
					continue
				}
				completed = append(completed, key)
			}

			for _, svc := range completed {
				log.Printf("service %q reached ready state in namespace %q\n", svc.Service, svc.Namespace)
				delete(services.SVCMap, svc)
			}

			if len(services.SVCMap) == 0 {
				return nil
			}
		}
	}
}
