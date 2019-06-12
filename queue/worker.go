package queue

// New takes a processor and returns both
// a client and a worker. The client allows
// pushing jobs to the queue (with CreateJob)
// and the worker can run those jobs using
// the given Processor
func New(p Processor) (Client, Worker) {
	return nil, nil
}
