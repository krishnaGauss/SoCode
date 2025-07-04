package server

import (
	"log/slog"
	"time"

	"github.com/krishnaGauss/SoCode/internal/storage"
)

type LogProcessor struct {
	queue *storage.RedisQueue
	storage *storage.PostgresStorage
	batchSize int 
	interval time.Duration
	stopChan chan struct{}
}

func NewLogProcessor(queue *storage.RedisQueue, storage *storage.PostgresStorage) *LogProcessor{
	return &LogProcessor{
        queue:     queue,
        storage:   storage,
        batchSize: 100,
        interval:  5 * time.Second,
        stopChan:  make(chan struct{}),
    }
}

func (p *LogProcessor) Start(){
	ticker:=time.NewTicker(p.interval)
	defer ticker.Stop()

	for{
		select{
		case <-ticker.C:
			p.processLogs()
		case <-p.stopChan:
			return
		}
	}
}


func (p *LogProcessor) Stop(){
	close(p.stopChan)
}

func (p *LogProcessor) processLogs(){
	logs,err:=p.queue.DequeueLogs(int64(p.batchSize))
	if err != nil {
        slog.Debug("failed to dequeue logs")
        return
    }

	if len(logs) == 0 {
        return
    }

	 if err := p.storage.StoreLogs(logs); err != nil {
        slog.Debug("failed to store logs")
        // Re-queue logs on failure
        for _, logEntry := range logs {
            p.queue.EnqueueLog(logEntry)
        }
        return
    }

    slog.Info("Processed logs")

}