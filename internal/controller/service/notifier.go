package service

import (
	"context"
	"distributed-configuration/pkg/utils"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisNotifier struct {
	rdb        *redis.Client
	log        *utils.Logger
	mu         sync.RWMutex
	listeners  []chan struct{}
	channelKey string
}

func NewRedisNotifier(rds *redis.Client, channelKey string, log *utils.Logger) *RedisNotifier {
	rn := &RedisNotifier{
		rdb:        rds,
		channelKey: channelKey,
		listeners:  make([]chan struct{}, 0),
		log:        log,
	}

	go rn.listenGlobalUpdates()

	return rn
}

func (r *RedisNotifier) listenGlobalUpdates() {
	r.log.Info("redis pubsub worker started", zap.String("channel", r.channelKey))

	for {
		ctx := context.Background()
		pubsub := r.rdb.Subscribe(ctx, r.channelKey)

		_, err := pubsub.Receive(ctx)
		if err != nil {
			r.log.Error("redis pubsub subscribe failed, retrying", zap.Error(err))
			pubsub.Close()
			time.Sleep(2 * time.Second)
			continue
		}

		ch := pubsub.Channel()
		for msg := range ch {
			r.log.Info("redis received update signal", zap.String("message", msg.String()))
			r.broadcastToLocal()
		}

		pubsub.Close()
		r.log.Warn("redis pubsub connection lost, attempting to reconnect")
		time.Sleep(1 * time.Second)
	}
}

func (r *RedisNotifier) broadcastToLocal() {
	r.mu.Lock()
	currentListeners := r.listeners
	r.listeners = make([]chan struct{}, 0)
	r.mu.Unlock()

	for _, ch := range currentListeners {
		select {
		case ch <- struct{}{}:
		default:
		}

		close(ch)
	}
}

func (r *RedisNotifier) PublishUpdate(ctx context.Context) error {
	r.log.Info("publishing global update signal", zap.String("channel", r.channelKey))
	return r.rdb.Publish(ctx, r.channelKey, "updated").Err()
}

func (r *RedisNotifier) Subscribe() chan struct{} {
	r.mu.Lock()
	defer r.mu.Unlock()
	ch := make(chan struct{}, 1)
	r.listeners = append(r.listeners, ch)
	return ch
}
