package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shimkek/omd-common/api"
	"github.com/shimkek/omd-common/broker"
	"github.com/shimkek/omd-kitchen/gateway"
	"go.opentelemetry.io/otel"
)

type Consumer struct {
	ordersGateway gateway.OrdersGateway
}

func NewConsumer(ordersGateway gateway.OrdersGateway) *Consumer {
	return &Consumer{ordersGateway}
}

func (c *Consumer) Listen(ch *amqp.Channel) {
	q, err := ch.QueueDeclare("", true, false, true, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = ch.QueueBind(q.Name, "", broker.OrderPaidEvent, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan struct{})

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			// Extract the headers
			ctx := broker.ExtractAMQPHeader(context.Background(), d.Headers)

			tr := otel.Tracer("amqp")
			_, messageSpan := tr.Start(ctx, fmt.Sprintf("AMQP - consume - %s", q.Name))

			o := &api.Order{}
			if err := json.Unmarshal(d.Body, o); err != nil {
				d.Nack(false, false)
				log.Printf("failed to unmarshal order: %v", err)
				continue
			}

			if o.Status == "paid" {
				cookOrder(o.OrderID)
				messageSpan.AddEvent(fmt.Sprintf("Order Cooked: %v", o))

				c.ordersGateway.UpdateOrderStatus(ctx, o.OrderID, "ready for pick-up")
				messageSpan.AddEvent(fmt.Sprintf("order.updated :%v", o))

			}
			messageSpan.End()
			d.Ack(false)
		}
	}()

	<-forever
}

func cookOrder(id string) {
	log.Printf("Cooking order %s", id)
	time.Sleep(time.Second * 10)
	log.Printf("Order %s cooked", id)
}
