package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

const MaxRetryCount = 3
const DLQ = "dlq_main"

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", user, password, host, port)
	conn, err := amqp.Dial(address)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}

	err = ch.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare an exchange:", err)
	}

	err = ch.ExchangeDeclare(OrderPaidEvent, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare an exchange:", err)
	}

	err = createDLQAndDLX(ch)
	if err != nil {
		log.Fatal("Failed to create DLQ:", err)
	}

	return ch, conn.Close
}

func HandleRetry(ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	retryCount, ok := d.Headers["x-retry-count"].(int64)
	if !ok {
		retryCount = 0
	}

	retryCount++
	d.Headers["x-retry-count"] = retryCount

	log.Printf("Retrying message %s, retry count: %d", d.Body, retryCount)

	if retryCount >= MaxRetryCount {
		// DLQ
		log.Printf("Moving message to DLQ %s", DLQ)

		return ch.PublishWithContext(context.Background(), "", DLQ, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}
	time.Sleep(time.Second * time.Duration(retryCount))

	return ch.PublishWithContext(context.Background(), d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Headers:      d.Headers,
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

func createDLQAndDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"main_queue",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLX
	dlx := "dlx_main"
	err = ch.ExchangeDeclare(
		dlx, "fanout", true, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	// Bind main queue to DLX
	err = ch.QueueBind(
		q.Name,
		"",
		dlx,
		false,
		nil,
	)
	if err != nil {
		return nil
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(
		DLQ, true, false, false, false, nil,
	)
	if err != nil {
		return err
	}
	return err
}

type AmqpHeaderCarrier map[string]interface{}

func (a AmqpHeaderCarrier) Get(k string) string {
	value, ok := a[k]
	if !ok {
		return ""
	}
	return value.(string)
}
func (a AmqpHeaderCarrier) Set(k string, v string) {
	a[k] = v
}

func (a AmqpHeaderCarrier) Keys() []string {
	keys := make([]string, len(a))
	i := 0
	for k := range a {
		keys[i] = k
		i++
	}
	return keys
}

func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(AmqpHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return carrier
}

func ExtractAMQPHeader(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, AmqpHeaderCarrier(headers))
}
