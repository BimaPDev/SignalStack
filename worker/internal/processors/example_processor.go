package processors

import "context"
type ExampleProcessor struct {
}
func (p *ExampleProcessor) Processor(ctx context.Context, jobID, userID string) ([]byte, error) {
	return nil, nil
}
