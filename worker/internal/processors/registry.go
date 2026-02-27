package processors

import (
	"context"
	"fmt"
)

type Processor interface {
	Processor(ctx context.Context, jobID, userID string) (output []byte, err error)
}

type Registry struct {
	m map[string]Processor
}

func NewRegistry() *Registry {
	return &Registry{
		m: make(map[string]Processor),
	}
}
func (r *Registry) Register(jobType string, p Processor) {
	r.m[jobType] = p
}

func (r *Registry) Get(jobType string) (Processor, error) {
	value, ok := r.m[jobType]
	if !ok {
		return nil, fmt.Errorf("processor not found for job type: %s", jobType)
	}
	return value, nil
}
