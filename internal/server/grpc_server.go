package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/krishnaGauss/SoCode/internal/models"
	"github.com/krishnaGauss/SoCode/internal/storage"
	"github.com/krishnaGauss/SoCode/proto/SoCode/proto"
)

type LogServer struct {
	proto.UnimplementedLogServiceServer
	queue   *storage.RedisQueue
	storage *storage.PostgresStorage
}

func NewLogServer(queue *storage.RedisQueue, storage *storage.PostgresStorage) *LogServer {
	return &LogServer{
		queue:   queue,
		storage: storage,
	}
}

func (s *LogServer) SendLog(ctx context.Context, req *proto.LogRequest) (*proto.LogResponse, error) {
	log := s.protoToModel(req)

	if err := s.queue.EnqueueLog(log); err != nil {
		return &proto.LogResponse{
			Success: false,
			Message: fmt.Sprintf("failed to enqueue log: %v", err),
		}, nil

	}

	return &proto.LogResponse{
        Success: true,
        Message: "Log received successfully",
    }, nil
}

func (s *LogServer) SendLogStream(stream proto.LogService_SendLogStreamServer) error {
	count:=0

	for{
		req, err:=stream.Recv()
		if err!=nil{
			break
		}

		log:=s.protoToModel(req)
		if err:=s.queue.EnqueueLog(log); err!=nil{
			slog.Info("failed to enqueue log:", slog.String(" ", err.Error()))
            continue
		}
		count++
	}

	return stream.SendAndClose(&proto.LogResponse{
        Success: true,
        Message: fmt.Sprintf("Processed %d logs", count),
    })
}

func (s *LogServer) protoToModel(req *proto.LogRequest) models.LogEntry {

}
