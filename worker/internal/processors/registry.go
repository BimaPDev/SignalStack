package processors

import "context"

// type Processor interface
// - Process(ctx context.Context, jobID, userID string) (output []byte, err error)

type Processor interface {
	Processor(ctx context.Context, jobID, userID string) (output []byte, err error)
}

// type Registry struct
// - m map[string]Processor

type Registry struct {
	m map[string]Processor
}

// func NewRegistry() *Registry
// - return Registry with empty map

func NewRegistry() *Registry {
	return &Registry {
		
	}
}

// func (r *Registry) Register(jobType string, p Processor)
// - add processor to map keyed by job type

// func (r *Registry) Get(jobType string) (Processor, error)
// - look up processor by job type
// - return error if not found
