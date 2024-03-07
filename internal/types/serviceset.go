package types

import (
	"fmt"
	"strings"
)

type Service struct {
	Service   string
	Namespace string
}

const (
	nsInferPath string = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
)

func NewService(value string) (*Service, error) {
	values := strings.Split(strings.TrimSpace(value), ",")
	if len(values) < 1 || len(values) > 2 {
		return nil, fmt.Errorf("bad service value specified: %q", value)
	}

	return &Service{
		Service: values[0],
		Namespace: func() string {
			if len(values) == 1 {
				return ns.String()
			}
			return values[1]
		}(),
	}, nil
}

type ServiceSet struct {
	SVCMap map[*Service]struct{}
}

func (s *ServiceSet) String() string {
	keys := make([]string, 0, len(s.SVCMap))
	for k := range s.SVCMap {
		keys = append(keys, k.Service)
	}
	return strings.Join(keys, ",")
}

func (s *ServiceSet) Set(value string) error {
	if s.SVCMap == nil {
		s.SVCMap = make(map[*Service]struct{})
	}

	svc, err := NewService(value)
	if err != nil {
		return err
	}

	s.SVCMap[svc] = struct{}{}
	return nil
}
