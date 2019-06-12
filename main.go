package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"time"

	"github.com/ingrammicro/backend-test/queue"
)

// piComputeData holds both the input and output of a pi processing job
// Total (both input and output) is the number of random points to pick
// in a [0,1)x[0,1) square.
// InCircle (only output) is the number of the randomly picked points that
// where inside the circle of radius 1 cented in (0,0).
type piComputeData struct {
	InCircle uint64 `json:"i"`
	Total    uint64 `json:"t"`
}

// piProcessor is a processor that can work out pi processing jobs
type piProcessor struct{}

// Process processes a pi processing job. To do so, it extracts piComputeData from
// the given job, computes it and stores it back into the job. It returns an error
// if any of the three operations fail.
func (pp piProcessor) Process(ctx context.Context, j queue.JobProcessingAccess) error {
	pcd := &piComputeData{}
	err := j.GetData(pcd)
	if err != nil {
		return err
	}
	err = pcd.Compute(ctx)
	if err != nil {
		return err
	}
	err = j.SetData(ctx, pcd)
	if err != nil {
		return err
	}
	return nil
}

// Compute picks a Total number of points in the [0,1)x[0,1) square
// and checks for each of one if they are inside the circle of radius 1 cented in (0,0).
// Specifically, given a (x,y) point, it checks whether x²+y² <= 1.
// It updates InCircle with the number of points that were inside.
func (pcd *piComputeData) Compute(ctx context.Context) error {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	for i := uint64(0); i < pcd.Total; i++ {
		x, y := r.Float64(), r.Float64()
		if (x*x)+(y*y) <= 1 {
			pcd.InCircle++
		}
	}
	return nil
}

func (pcd *piComputeData) String() string {
	return fmt.Sprintf("%d/%d", pcd.InCircle, pcd.Total)
}

// Marshal encodes the piComputeData into the returned byte slice
// as JSON.
func (pcd *piComputeData) Marshal() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(pcd)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Unmarshal decodes the JSON in the given byte slice
// and fills the piComputeData with it.
func (pcd *piComputeData) Unmarshal(b []byte) error {
	return json.NewDecoder(bytes.NewReader(b)).Decode(pcd)
}

// main pushes numberOfJobs pi processing jobs (each computing a million points),
// starts 10 workers, waits for all the jobs to be processed and then aggregates
// the results of the jobs to approximate pi. Finally, it prints the approximation
// and exits orderly.
func main() {
	const numberOfJobs = 10000
	ctx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()
	client, worker := queue.New(piProcessor{})
	log.Printf("Pushing %d pi processing jobs...", numberOfJobs)
	for i := 0; i < numberOfJobs; i++ {
		client.CreateJob(ctx, fmt.Sprintf("j-%d", i), &piComputeData{Total: 1000000})
	}
	log.Print("Starting 10 workers...")
	workerStopped := make(chan struct{})
	go func() {
		worker.Run(ctx, 10)
		close(workerStopped)
	}()
	log.Print("Waiting for results and aggregating them...")
	result := &big.Rat{}
	for i := 0; i < numberOfJobs; i++ {
		jobID := fmt.Sprintf("j-%d", i)
		for {
			job, err := client.GetJob(ctx, jobID)
			if err != nil {
				log.Fatal(err)
			}
			if job == nil {
				log.Fatalf("Job %q could not be found", jobID)
			}
			state := job.State()
			if state == queue.Failed {
				log.Fatal(job.Error())
			}
			if state == queue.Finished {
				var partialResult piComputeData
				err = job.GetData(&partialResult)
				if err != nil {
					log.Fatal(err)
				}
				result = result.Add(result, big.NewRat(4*int64(partialResult.InCircle), int64(partialResult.Total)))
				break
			}
			time.Sleep(5 * time.Second) // Wait a bit for the job to finish
		}
	}
	cancelCtx()
	result = result.Mul(result, big.NewRat(1, numberOfJobs))
	log.Printf("Result is %+v = %s", result, result.FloatString(20))
	log.Printf("Preparing to exit...")
	<-workerStopped
	log.Printf("Exiting")
}
