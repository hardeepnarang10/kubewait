package kubernetes

import (
	"errors"
	"fmt"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func ClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		if errors.Is(err, rest.ErrNotInCluster) {
			return nil, fmt.Errorf("client not running in a cluster: %w", err)
		}
		return nil, fmt.Errorf("unable to create in-cluster config using assigned service account: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("unable to create kubernetes clienset with prepared config: %w", err)
	}

	return clientset, nil
}
