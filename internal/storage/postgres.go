package storage

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/krishnaGauss/SoCode/internal/config"
	"github.com/krishnaGauss/SoCode/internal/models"
	_ "github.com/lib/pq"
)

type PostgresStorage struct {
	db *sqlx.DB
}

func NewPostgresStorage(cfg *config.DatabaseConfig) (*PostgresStorage, error) {
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s", cfg.Host, cfg.Port, cfg.Password, cfg.Database, cfg.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		slog.Info("couldn't connect to db.")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("database unreachable ", err)
		return nil, err
	}

	storage := &PostgresStorage{db: db}

	if err := storage.createTables(); err != nil {
		slog.Info("couldn't create table")
		return nil, err
	}

	return storage, nil

}

func (s *PostgresStorage) createTables() error {
	//query here
	query := `
		CREATE TABLE IF NOT EXISTS logs(
			id VARCHAR(255) PRIMARY KEY,
			timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
			level VARCHAR(20) NOT NULL,
			message TEXT NOT NULL,
			source VARCHAR(255) NOT NULL,
			service VARCHAR(255) NOT NULL,
			host VARCHAR(255) NOT NULL,
        	tags JSONB,
        	metadata JSONB,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_logs_timestamp ON logs(timestamp);
		CREATE INDEX IF NOT EXISTS idx_logs_level ON logs(level);
    	CREATE INDEX IF NOT EXISTS idx_logs_source ON logs(source);
    	CREATE INDEX IF NOT EXISTS idx_logs_service ON logs(service);
    	CREATE INDEX IF NOT EXISTS idx_logs_host ON logs(host);
    	CREATE INDEX IF NOT EXISTS idx_logs_tags ON logs USING GIN(tags);
    	CREATE INDEX IF NOT EXISTS idx_logs_message ON logs USING GIN(to_tsvector('english', message));
		`

	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStorage) StoreLogs(logs []models.LogEntry) error {
	if len(logs) == 0 {
		slog.Debug("log length is zero")
		return nil
	}

	query := `
		INSERT INTO logs(id, timestamp, level, message, source, service, host, tags, metadata) 
		VALUES (:id, :timestamp, :level, :message, :source, :service, :host, :tags, :metadata)
		ON CONFLICT (id) DO NOTHING;
	`

	_, err := s.db.NamedExec(query, logs)
	return err
}

func (s *PostgresStorage) QueryLogs(query models.LogQuery) ([]models.LogEntry, error){
	var conditions []string
	var args []interface{}
	argCount:=0

	baseQuery:="SELECT id, timestamp, level, message, source, service, host, tags, metadata FROM logs"

	if query.StartTime!=nil{
		argCount++
		conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", argCount))
		args = append(args, *query.StartTime)
	}

	return nil, nil //to be completed
}