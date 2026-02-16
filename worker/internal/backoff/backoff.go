package backoff

// func NextDelay(attempt int) time.Duration
// - calculate exponential backoff delay for given attempt number
// - apply jitter to avoid thundering herd
// - cap maximum delay at a reasonable upper bound
