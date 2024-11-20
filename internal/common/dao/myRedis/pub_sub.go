package myRedis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"sync"
)

type WorkWork func(*redis.Message) error

type Job struct {
	Channel string
	After   WorkWork
}

func (db *DB) SubscribeWork(ctx context.Context, job Job) error {
	pubsub := db.Subscribe(ctx, job.Channel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
	}

	ch := pubsub.Channel()
	for msg := range ch {
		if err := job.After(msg); err != nil {
			return err
		}
	}
	return nil
}

func (db *DB) SubscribeWorks(ctx context.Context, jobs ...Job) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(jobs))

	for _, job := range jobs {
		wg.Add(1)
		go func(job Job) {
			defer wg.Done()
			if err := db.SubscribeWork(ctx, job); err != nil {
				errChan <- err
			}
		}(job)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var finalErr error
	for err := range errChan {
		finalErr = err
	}
	return finalErr
}
