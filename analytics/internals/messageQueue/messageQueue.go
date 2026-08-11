package messagequeue

import (
	"context"
	"fmt"
	"time"

	"github.com/rabbitmq/amqp091-go"
)

type MessageQueue interface {
	Push(string) error
	Pull()
}

type RabbmitMQ struct {
	connection *amqp091.Connection
	channel    *amqp091.Channel
	queue      amqp091.Queue
}

func NewRabbmitMQ() (*RabbmitMQ, error) {
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		"url-shortner",
		true,
		false,
		false,
		false,
		amqp091.Table{
			amqp091.QueueTypeArg: amqp091.QueueTypeQuorum,
		},
	)
	if err != nil {
		return nil, err
	}

	return &RabbmitMQ{
		connection: conn,
		channel:    ch,
		queue:      q,
	}, nil
}

func (r *RabbmitMQ) Push(urlID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := r.channel.PublishWithContext(ctx,
		"",
		r.queue.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType: "text/plain",
			Body:        []byte(urlID),
		},
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *RabbmitMQ) Pull() {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			fmt.Println("Message recieved : ", string(d.Body))
		}
		close(forever)
	}()

	<-forever
}

func (r *RabbmitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.connection != nil {
		r.connection.Close()
	}
}
