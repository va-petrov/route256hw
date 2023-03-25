package jobs

import (
	"context"
	"fmt"
	"log"
	"time"
)

type Job struct {
	Name    string
	Period  time.Duration
	JobFunc func(ctx context.Context) error
	ticker  *time.Ticker
}

func NewJob(name string, job func(ctx context.Context) error, period time.Duration) *Job {
	return &Job{
		Name:    name,
		Period:  period,
		JobFunc: job,
	}
}

func (job *Job) Run(ctx context.Context) error {
	if job.ticker != nil {
		return fmt.Errorf("this job is already running: %v", job.Name)
	}

	job.ticker = time.NewTicker(job.Period)
	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				job.ticker.Stop()
				return
			case t := <-job.ticker.C:
				log.Printf("Running job %v at %v", job.Name, t.Format("2006-01-02 15:04:05"))
				err := job.JobFunc(ctx)
				if err != nil {
					log.Printf("JobFunc %v funished at %v with error %v", job.Name, time.Now().Format("2006-01-02 15:04:05"), err)
				} else {
					log.Printf("JobFunc %v funished successfuly at %v", job.Name, time.Now().Format("2006-01-02 15:04:05"))
				}
			}
		}
	}(ctx)

	return nil
}

func (job *Job) Stop() error {
	if job.ticker == nil {
		return fmt.Errorf("job is not running: %v", job.Name)
	}

	job.ticker.Stop()
	job.ticker = nil
	return nil
}
