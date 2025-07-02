package storage

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"strings"

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

func (s *PostgresStorage) QueryLogs(query models.LogQuery) ([]models.LogEntry, error) {
	var conditions []string
	var args []interface{}
	argCount := 0

	baseQuery := "SELECT id, timestamp, level, message, source, service, host, tags, metadata FROM logs"

	if query.StartTime != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("timestamp >= $%d", argCount))
		args = append(args, *query.StartTime)
	}

	if query.EndTime != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("timestamp <= $%d", argCount))
		args = append(args, *query.EndTime)
	}

	if len(query.Level) > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("level = ANY($%d)", argCount))
		levels := make([]string, len(query.Level))
		for i, level := range query.Level {
			levels[i] = string(level)
		}

		args = append(args, levels)
	}

	if len(query.Source) > 0 {
		argCount++
		conditions = append(conditions, fmt.Sprintf("source = ANY($%d)", argCount))
		args = append(args, query.Source)
	}

	if query.Search != "" {
		argCount++
		conditions = append(conditions, fmt.Sprintf("message ILIKE $%d", argCount))
		args = append(args, "%"+query.Search+"%")
	}

	if len(conditions) > 0 {
		baseQuery += "WHERE" + strings.Join(conditions, " AND ")
	}

	baseQuery += "ORDER BY timestamp DESC" //remember to put ;

	if query.Limit > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, query.Limit)
	}

	if query.Offset > 0 {
		argCount++
		baseQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, query.Offset)
	}

	rows, err := s.db.Query(baseQuery, args...) //using ... to help in using argument as an interface slice individually

	if err != nil {
		slog.Debug("cannot execute query in postgres")
		return nil, err
	}

	defer rows.Close()

	var logs []models.LogEntry
	for rows.Next() {
		var log models.LogEntry
		var tagsJSON, metadataJSON sql.NullString

		err := rows.Scan(
			&log.ID, &log.Timestamp, &log.Level, &log.Message,
			&log.Source, &log.Service, &log.Host, &tagsJSON, &metadataJSON,
		)

		if err != nil {
			return nil, err
		}

		 if tagsJSON.Valid {
            json.Unmarshal([]byte(tagsJSON.String), &log.Tags)
        }
        if metadataJSON.Valid {
            log.Metadata = json.RawMessage(metadataJSON.String)
        }

        logs = append(logs, log)
	}

	return logs, nil

}


func (s *PostgresStorage) Close() error{
	return s.db.Close()
}
