package myRedis

import (
	"context"
	"github.com/redis/go-redis/v9"
)

type WorkWork func(*redis.Message)

type Job struct {
	Channel string
	After   WorkWork
}

func (db *Native) SubscribeWork(ctx context.Context, job Job) {
	pubsub := db.Subscribe(ctx, job.Channel)
	defer pubsub.Close()

	if _, err := pubsub.Receive(ctx); err != nil {
	}

	ch := pubsub.Channel()
	for msg := range ch {
		job.After(msg)
	}
}

func (db *Native) SubscribeWorks(ctx context.Context, jobs ...Job) {
	for _, job := range jobs {
		go func(job Job) {
			db.SubscribeWork(ctx, job)
		}(job)
	}
}
