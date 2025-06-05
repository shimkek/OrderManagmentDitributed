package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-common/broker"
	"go.opentelemetry.io/otel"
)

type consumer struct {
	service PaymentsService
}

func NewConsumer(service PaymentsService) *consumer {
	return &consumer{
		service: service,
	}
}

func (c *consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare(broker.OrderCreatedEvent, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register consumer: %v", err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {

			ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			log.Printf("Received a message: %s", d.Body)

			o := &api.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("Failed to unmarshal order: %v", err)
				continue
			}

			paymentLink, err := c.service.CreatePayment(context.Background(), o)
			if err == nil {
				log.Printf("Failed to create payment: %v", err)

				if err := broker.HandleRetry(ch, &d); err != nil {
					log.Printf("Error handling retry: %v", err)
				}
				d.Nack(false, false)
				continue
			}

			messageSpan.AddEvent(fmt.Sprintf("payment.created: %s", paymentLink))
			messageSpan.End()

			d.Ack(false)
			log.Printf("Payment created with link: %s for Order ID: %s", paymentLink, o.OrderID)
		}
	}()

	<-forever
}

func (c *consumer) Consume(ctx context.Context, order *api.Order) error {
	// Call the service to create a payment
	paymentLink, err := c.service.CreatePayment(ctx, order)
	if err != nil {
		return err
	}

	// Log or handle the payment ID as needed
	log.Printf("Payment created with link: %s for Order ID: %s", paymentLink, order.OrderID)
	return nil
}
