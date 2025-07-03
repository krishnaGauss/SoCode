package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/krishnaGauss/SoCode/internal/config"
	"github.com/krishnaGauss/SoCode/internal/models"
	"github.com/redis/go-redis/v9"
)

type RedisQueue struct {
	client *redis.Client
	ctx context.Context
}

func NewRedisQueue(cfg *config.RedisConfig) (*RedisQueue, error){
	client:=redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
        Password: cfg.Password,
        DB:       cfg.DB,
	})

	ctx := context.Background()
    _, err := client.Ping(ctx).Result()
    if err != nil {
        return nil, fmt.Errorf("failed to connect to Redis: %w", err)
    }

    return &RedisQueue{
        client: client,
        ctx:    ctx,
    }, nil
}

func (r *RedisQueue) EnqueueLog(log models.LogEntry) error{
	data,err:=json.Marshal(log)
	if err!=nil{
		slog.Debug("error in marshalling redis enqueue")
		return err
	}

	return r.client.LPush(r.ctx, "log_queue", data).Err()
}

func (r *RedisQueue) DequeueLogs(count int64) ([]models.LogEntry, error){
	results, err := r.client.RPopCount(r.ctx, "log_queue", int(count)).Result()
    if err != nil {
        return nil, err
    }

    logs := make([]models.LogEntry, 0, len(results))
    for _, result := range results {
        var log models.LogEntry
        if err := json.Unmarshal([]byte(result), &log); err != nil {
            continue // Skip malformed logs
        }
        logs = append(logs, log)
    }

    return logs, nil
}

func (r *RedisQueue) QueueLength() (int64, error) {
    return r.client.LLen(r.ctx, "log_queue").Result()
}

func (r *RedisQueue) Close() error {
    return r.client.Close()
}