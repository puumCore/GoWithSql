package main

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/url"
	"time"
)

const (
	host     = "rabbitmq"
	port     = 5672
	user     = "developer"
	password = "G8ro_Hw8tdhr1Lz]s2~|"
	vhost    = "dev"
	exchange = "tutorials"
	queue    = "go-lang-queue"
	key      = "passthrough"
)

func getMqConnection() (*amqp.Connection, error) {
	connString := fmt.Sprintf("amqp://%s:%s@%s:%d/%s", user, url.QueryEscape(password), host, port, vhost)
	fmt.Println("***** [rabbit-mq] Obtaining connection...")
	connection, err := amqp.Dial(connString)
	checkError(err)
	fmt.Println("**** [rabbit-mq] Opening channel...")
	channel, err := connection.Channel()
	checkError(err)
	defer func(channel *amqp.Channel) {
		err := channel.Close()
		checkError(err)
	}(channel)

	fmt.Println("*** [rabbit-mq] Connected!")

	declare(channel)
	publish(channel, "{\"Id\":4,\"CreatedAt\":\"2024-08-29T12:45:03.463386-07:00\",\"Name\":\"Jane Doe\",\"Phone\":\"254716686433\",\"Email\":\"emandela60@gmail.com\",\"KraPin\":{\"String\":\"\",\"Valid\":false},\"Username\":\"jane\",\"Password\":\"$2a$10$XTGbm.c2wGH7JQvLBJOJD.NfIGaiJMno3nU7rmfbom39QZibFBwpa\"}")
	consume(channel)

	defer func(connection *amqp.Connection) {
		err := connection.Close()
		checkError(err)
	}(connection)

	return connection, err
}

func consume(c *amqp.Channel) {
	deliveries, err := c.Consume(
		queue,
		"go-consumer",
		true,
		false,
		false,
		false,
		nil)
	checkError(err)
	for d := range deliveries {
		fmt.Printf("** [rabbit-mq] Received %s\n", d.Body)
	}
}

func publish(c *amqp.Channel, data string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := c.PublishWithContext(ctx,
		exchange, // exchange
		key,      // routing key
		false,    // mandatory
		false,    // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(data),
		})
	checkError(err)
}

func declare(c *amqp.Channel) {
	err := c.ExchangeDeclare(
		exchange,
		"direct",
		true,
		false,
		false,
		false,
		nil)
	checkError(err)
	_, err = c.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil)
	checkError(err)
	err = c.QueueBind(
		queue,
		key,
		exchange,
		false,
		nil)
	checkError(err)
}
