package chatMessageHistories

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	defaultConnectionString = "postgresql://postgres:mypassword@localhost/chat_history"
	defaultTableName        = "message_store"
)

type PostgresChatMessageHistory struct {
	sessionID      string
	tableName      string
	connectionPool *pgxpool.Pool
}

func NewPostgresChatMessageHistory(sessionID string, connectionString string, tableName string) (*PostgresChatMessageHistory, error) {
	if connectionString == "" {
		connectionString = defaultConnectionString
	}
	if tableName == "" {
		tableName = defaultTableName
	}

	connectionPool, err := pgxpool.Connect(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	postgresHistory := &PostgresChatMessageHistory{
		sessionID:      sessionID,
		tableName:      tableName,
		connectionPool: connectionPool,
	}

	if err := postgresHistory.createTableIfNotExists(); err != nil {
		return nil, err
	}

	return postgresHistory, nil
}

func (p *PostgresChatMessageHistory) createTableIfNotExists() error {
	createTableQuery := `
        CREATE TABLE IF NOT EXISTS %s (
            id SERIAL PRIMARY KEY,
            session_id TEXT NOT NULL,
            message JSONB NOT NULL
        );
    `
	query := fmt.Sprintf(createTableQuery, p.tableName)
	_, err := p.connectionPool.Exec(context.Background(), query)
	return err
}

func (p *PostgresChatMessageHistory) Messages() ([]rootSchema.BaseMessageInterface, error) {
	query := fmt.Sprintf("SELECT message FROM %s WHERE session_id = $1;", p.tableName)
	rows, err := p.connectionPool.Query(context.Background(), query, p.sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := []rootSchema.BaseMessageInterface{}
	for rows.Next() {
		var messageJSON string
		err := rows.Scan(&messageJSON)
		if err != nil {
			return nil, err
		}

		var message rootSchema.BaseMessageInterface
		err = json.Unmarshal([]byte(messageJSON), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (p *PostgresChatMessageHistory) AddUserMessage(message string) error {
	humanMessage := rootSchema.NewHumanMessage(message)
	return p.append(humanMessage)
}

func (p *PostgresChatMessageHistory) AddAIMessage(message string) error {
	aiMessage := rootSchema.NewAIMessage(message)
	return p.append(aiMessage)
}

func (p *PostgresChatMessageHistory) append(message rootSchema.BaseMessageInterface) error {
	query := fmt.Sprintf("INSERT INTO %s (session_id, message) VALUES ($1, $2);", p.tableName)
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}
	_, err = p.connectionPool.Exec(context.Background(), query, p.sessionID, messageJSON)
	return err
}

func (p *PostgresChatMessageHistory) Clear() error {
	query := fmt.Sprintf("DELETE FROM %s WHERE session_id = $1;", p.tableName)
	_, err := p.connectionPool.Exec(context.Background(), query, p.sessionID)
	return err
}
