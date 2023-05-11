package chatMessageHistories

import (
	"context"
	"encoding/json"
	"github.com/William-Bohm/langchain-go/langchain-go/rootSchema"
	"github.com/go-redis/redis/v8"
	"time"
)

type RedisChatMessageHistory struct {
	sessionID   string
	redisClient *redis.Client
	keyPrefix   string
	ttl         *time.Duration
}

func NewRedisChatMessageHistory(sessionID string, url string, keyPrefix string, ttl *time.Duration) (*RedisChatMessageHistory, error) {
	if url == "" {
		url = "redis://localhost:6379/0"
	}
	if keyPrefix == "" {
		keyPrefix = "message_store:"
	}

	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}

	redisClient := redis.NewClient(opt)

	return &RedisChatMessageHistory{
		sessionID:   sessionID,
		redisClient: redisClient,
		keyPrefix:   keyPrefix,
		ttl:         ttl,
	}, nil
}

func (r *RedisChatMessageHistory) key() string {
	return r.keyPrefix + r.sessionID
}

func (r *RedisChatMessageHistory) Messages() ([]rootSchema.BaseMessageInterface, error) {
	items, err := r.redisClient.LRange(context.Background(), r.key(), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	messages := []rootSchema.BaseMessageInterface{}
	for _, item := range items {
		var message rootSchema.BaseMessageInterface
		err = json.Unmarshal([]byte(item), &message)
		if err != nil {
			return nil, err
		}
		messages = append(messages, message)
	}

	return messages, nil
}

func (r *RedisChatMessageHistory) AddUserMessage(message string) error {
	humanMessage := rootSchema.NewHumanMessage(message)
	return r.append(humanMessage)
}

func (r *RedisChatMessageHistory) AddAIMessage(message string) error {
	aiMessage := rootSchema.NewAIMessage(message)
	return r.append(aiMessage)
}

func (r *RedisChatMessageHistory) append(message rootSchema.BaseMessageInterface) error {
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return err
	}

	err = r.redisClient.LPush(context.Background(), r.key(), string(messageJSON)).Err()
	if err != nil {
		return err
	}

	if r.ttl != nil {
		err = r.redisClient.Expire(context.Background(), r.key(), *r.ttl).Err()
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RedisChatMessageHistory) Clear() error {
	return r.redisClient.Del(context.Background(), r.key()).Err()
}
