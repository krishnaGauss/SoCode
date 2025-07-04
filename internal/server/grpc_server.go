package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/krishnaGauss/SoCode/internal/models"
	"github.com/krishnaGauss/SoCode/internal/storage"
	"github.com/krishnaGauss/SoCode/proto/SoCode/proto"

	"github.com/google/uuid"
	// "google.golang.org/grpc"
    "google.golang.org/protobuf/types/known/timestamppb"
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

func (s *LogServer) QueryLogs(ctx context.Context, req *proto.QueryRequest) (*proto.QueryResponse, error){
	query:=models.LogQuery{
		Search: req.Search,
        Limit:  int(req.Limit),
        Offset: int(req.Offset),
	}

	if req.StartTime!=nil{
		startTime:=req.StartTime.AsTime()
		query.StartTime=&startTime
	}

	if req.EndTime != nil {
        endTime := req.EndTime.AsTime()
        query.EndTime = &endTime
    }

	for _, level:=range req.Levels{
		query.Level = append(query.Level, models.LogLevel(level))
	}

	 query.Source = req.Sources
    query.Service = req.Services
    query.Host = req.Hosts
    query.Tags = req.Tags

	 logs, err := s.storage.QueryLogs(query)
    if err != nil {
        return nil, err
    }

    response := &proto.QueryResponse{
        Total: int32(len(logs)),
	}

	 for _, log := range logs {
        response.Logs = append(response.Logs, s.modelToProto(log))
    }

    return response, nil
}

func (s *LogServer) protoToModel(req *proto.LogRequest) models.LogEntry {
	log:= models.LogEntry{
		ID:      req.Id,
        Level:   models.LogLevel(req.Level),
        Message: req.Message,
        Source:  req.Source,
        Service: req.Service,
        Host:    req.Host,
        Tags:    req.Tags,
	}

	if req.Id==""{
		log.ID=uuid.New().String()
	}

	 if req.Timestamp != nil {
        log.Timestamp = req.Timestamp.AsTime()
    } else {
        log.Timestamp = time.Now()
    }

    if req.Metadata != "" {
        log.Metadata = []byte(req.Metadata)
    }

    return log
}

func (s *LogServer) modelToProto(log models.LogEntry) *proto.LogRequest {
    return &proto.LogRequest{
        Id:        log.ID,
        Timestamp: timestamppb.New(log.Timestamp),
        Level:     string(log.Level),
        Message:   log.Message,
        Source:    log.Source,
        Service:   log.Service,
        Host:      log.Host,
        Tags:      log.Tags,
        Metadata:  string(log.Metadata),
    }
}

