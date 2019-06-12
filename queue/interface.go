package queue

import (
	"context"
)

// State represents the state of a job
type State string

// MarshalUnmarshaler wraps the Marshal and Unmarshal methods.
// Job payloads must implement this interface that allows
// converting them into a slice of bytes and back.
//
// The Marshal method marshals the object into a stream of bytes
// and returns the slice of bytes of an error.
//
// The Unmarshal method unmarshals an object from a stream of bytes
// into itself and may return an error.
//
// Implementations should ensure the unmarshaling of the slice of bytes
// provided by the marshaling of an object provides essentially the same
// object.
//
// This interface has already an implementation by us in the main.go file.
type MarshalUnmarshaler interface {
	Marshal() ([]byte, error)
	Unmarshal([]byte) error
}

// Job is the interface that wraps the methods
// use to read data from a job obtained from a client.
//
// The ID methods provide the ID of the job.
//
// The GetData method unmarshals into data the payload of the job.
//
// The State method returns the State of the job.
//
// The Error method returns a string describing the error with which a job failed.
type Job interface {
	ID() string
	GetData(data MarshalUnmarshaler) error
	State() State
	Error() string
}

// JobProcessingAccess is just the Job interface with an extra method
// that allows setting the data payload of the job. It is meant to be
// used by job processors launched by workers, which will need to store
// their results there.
//
// The SetData method takes some data and updates the
// job's payload data with it. Implementations should use the context
// argument to allow canceling or expiration of the SetData operation,
// and return an error in that case or if the payload data update cannot
// be performed.
type JobProcessingAccess interface {
	Job
	SetData(ctx context.Context, data MarshalUnmarshaler) error
}

// A Processor defines the worker's job execution.
// It returns an error:
//  * If the error is not nil, the job is marked as Failed.
//  * If the error is nil the job is marked as finished
//    (successfully).
//
// This interface has already an implementation by us in the main.go file.
type Processor interface {
	Process(ctx context.Context, j JobProcessingAccess) error
}

// Worker is an interface that wraps the Run method, which
// allows processing jobs in a queue.
//
// Implementations should attempt to process as many jobs as
// the given worker integer simultaneously using a Processor.
// They should also use the given context to allow users to timeout
// or cancel processing, returning only after all workers have stopped.
type Worker interface {
	Run(ctx context.Context, workers int) error
}

// Client is an interface that allows pushing jobs into a queue
// and querying their state and results.
//
// Implementations of GetJob should return
//  * a nil job and a nil error when the job is not found
//  * the job and a nil error when the job is found
//  * a nil job and an error, when some error prevents the retrieval
//    of the job
type Client interface {
	CreateJob(ctx context.Context, id string, initialData MarshalUnmarshaler) error
	GetJob(ctx context.Context, id string) (Job, error)
}

const (
	// Queued but not processing yet
	Queued State = "queued"
	// Processing - The worker got active
	Processing State = "processing"
	// Failed by worker or from outside (timeout) - See Jobs's Error field
	Failed State = "failed"
	// Finished successfully
	Finished State = "finished"
)
