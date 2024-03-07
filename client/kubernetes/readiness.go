package kubernetes

import (
	"context"
	"fmt"

	"github.com/hardeepnarang10/kubewait/internal/types"
	corev1 "k8s.io/api/core/v1"
	k8serror "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
)

const (
	readyCondition corev1.PodConditionType = "Ready"
	readyStatus    corev1.ConditionStatus  = "True"
)

func Readiness(cli clientv1.CoreV1Interface) func(context.Context, types.Service) error {
	return func(ctx context.Context, service types.Service) error {
		svc, err := cli.Services(service.Namespace).Get(ctx, service.Service, metav1.GetOptions{})
		if err != nil {
			if k8serror.IsNotFound(err) {
				return fmt.Errorf("service %q not found in namespace %q: %w", service.Service, service.Namespace, err)
			}
			return fmt.Errorf("unable to get service %q in namespace %q: %w", service.Service, service.Namespace, err)
		}

		ls := labels.Set(svc.Spec.Selector).AsSelector().String()
		pods, err := cli.Pods(service.Namespace).List(ctx, metav1.ListOptions{LabelSelector: ls})
		if err != nil {
			return fmt.Errorf("unable to list associated pods with labels %q in namespace %q: %w", ls, service.Namespace, err)
		}

		if len(pods.Items) == 0 {
			return fmt.Errorf("no pods found with labels %q in namespace %q", ls, service.Namespace)
		}

		for _, pod := range pods.Items {
			for _, condition := range pod.Status.Conditions {
				if condition.Type == readyCondition && condition.Status == readyStatus {
					return nil
				}
			}
		}
		return fmt.Errorf("pods with labels %q in namespace %q are not ready", ls, service.Namespace)
	}
}
