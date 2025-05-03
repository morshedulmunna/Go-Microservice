package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/morshedulmunna/go-microservice/pkg/logger"
	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

// RabbitMQ represents a RabbitMQ connection
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	logger  *zap.Logger
	config  *Config
}

// Config represents RabbitMQ configuration
type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	VHost    string
}

// Message represents a message to be published
type Message struct {
	Type    string      `json:"type"`
	Data    interface{} `json:"data"`
	Time    time.Time   `json:"time"`
	Version string      `json:"version"`
}

// NewRabbitMQ creates a new RabbitMQ instance
func NewRabbitMQ(config *Config) (*RabbitMQ, error) {
	logger := logger.GetLogger()

	// Create connection URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.VHost,
	)

	// Connect to RabbitMQ
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	// Create channel
	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %v", err)
	}

	return &RabbitMQ{
		conn:    conn,
		channel: channel,
		logger:  logger,
		config:  config,
	}, nil
}

// Close closes the RabbitMQ connection
func (r *RabbitMQ) Close() error {
	if err := r.channel.Close(); err != nil {
		return fmt.Errorf("failed to close channel: %v", err)
	}
	if err := r.conn.Close(); err != nil {
		return fmt.Errorf("failed to close connection: %v", err)
	}
	return nil
}

// DeclareExchange declares an exchange
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue declares a queue
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue binds a queue to an exchange
func (r *RabbitMQ) BindQueue(queue, exchange, routingKey string) error {
	return r.channel.QueueBind(
		queue,      // queue name
		routingKey, // routing key
		exchange,   // exchange
		false,      // no-wait
		nil,        // arguments
	)
}

// Publish publishes a message to an exchange
func (r *RabbitMQ) Publish(ctx context.Context, exchange, routingKey string, msg *Message) error {
	// Convert message to JSON
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	// Create publishing
	publishing := amqp.Publishing{
		ContentType:  "application/json",
		Body:         body,
		Timestamp:    time.Now(),
		MessageId:    fmt.Sprintf("%d", time.Now().UnixNano()),
		DeliveryMode: amqp.Persistent,
	}

	// Publish message
	err = r.channel.Publish(
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		publishing, // message
	)
	if err != nil {
		return fmt.Errorf("failed to publish message: %v", err)
	}

	r.logger.Info("Message published",
		zap.String("exchange", exchange),
		zap.String("routing_key", routingKey),
		zap.String("message_id", publishing.MessageId),
	)

	return nil
}

// Consume consumes messages from a queue
func (r *RabbitMQ) Consume(ctx context.Context, queue string, handler func(msg *Message) error) error {
	// Create consumer
	msgs, err := r.channel.Consume(
		queue, // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %v", err)
	}

	// Start consuming messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case delivery := <-msgs:
				// Parse message
				var msg Message
				if err := json.Unmarshal(delivery.Body, &msg); err != nil {
					r.logger.Error("Failed to unmarshal message",
						zap.Error(err),
						zap.String("queue", queue),
					)
					delivery.Reject(false)
					continue
				}

				// Handle message
				if err := handler(&msg); err != nil {
					r.logger.Error("Failed to handle message",
						zap.Error(err),
						zap.String("queue", queue),
					)
					delivery.Reject(true)
					continue
				}

				// Acknowledge message
				if err := delivery.Ack(false); err != nil {
					r.logger.Error("Failed to acknowledge message",
						zap.Error(err),
						zap.String("queue", queue),
					)
				}

				r.logger.Info("Message processed",
					zap.String("queue", queue),
					zap.String("message_id", delivery.MessageId),
				)
			}
		}
	}()

	return nil
}

// CreateDeadLetterQueue creates a dead letter queue and binds it to an exchange
func (r *RabbitMQ) CreateDeadLetterQueue(queue, exchange string) error {
	// Declare dead letter exchange
	dlxName := fmt.Sprintf("%s.dlx", exchange)
	if err := r.DeclareExchange(dlxName, "direct"); err != nil {
		return err
	}

	// Declare dead letter queue
	dlqName := fmt.Sprintf("%s.dlq", queue)
	args := amqp.Table{
		"x-dead-letter-exchange": exchange,
		"x-message-ttl":          int32(24 * 60 * 60 * 1000), // 24 hours
	}

	_, err := r.channel.QueueDeclare(
		dlqName, // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		args,    // arguments
	)
	if err != nil {
		return err
	}

	// Bind dead letter queue to exchange
	return r.BindQueue(dlqName, dlxName, queue)
}

// RetryMessage retries a message from the dead letter queue
func (r *RabbitMQ) RetryMessage(ctx context.Context, queue, exchange string) error {
	dlqName := fmt.Sprintf("%s.dlq", queue)

	// Get message from dead letter queue
	msg, ok, err := r.channel.Get(dlqName, false)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	// Publish message back to original exchange
	err = r.channel.Publish(
		exchange, // exchange
		queue,    // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType:  msg.ContentType,
			Body:         msg.Body,
			Timestamp:    time.Now(),
			MessageId:    msg.MessageId,
			DeliveryMode: amqp.Persistent,
		},
	)
	if err != nil {
		return err
	}

	// Acknowledge message from dead letter queue
	return msg.Ack(false)
}
