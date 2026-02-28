package processors

import "context"

type ExportProcessor struct {
}

func (p *ExportProcessor) Processor(ctx context.Context, jobID, userID string) ([]byte, error) {
	return nil, nil
}
