package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/streadway/amqp"
)

func main() {
	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	if len(amqpServerURL) <= 0 {
		amqpServerURL = "amqp://test:test@localhost:5672/"
	}

	isConnected := false
	var serverRMQ *amqp.Connection
	for !isConnected {
		var err error
		serverRMQ, err = amqp.Dial(amqpServerURL)
		if err == nil {
			isConnected = true
		} else {
			fmt.Println("Failed to connect RabbitMQ, waiting 5 seconds...")
			time.Sleep(5 * time.Second)
		}
	}
	defer serverRMQ.Close()

	channelRMQ, err := serverRMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRMQ.Close()

	queueNames := []string{"queue-1", "queue-2"}
	for _, qName := range queueNames {
		_, err = channelRMQ.QueueDeclare(
			qName,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			panic(err)
		}
	}

	for i := 1; i <= 10; i++ {
		targetQueue := "queue-" + strconv.Itoa(((i-1)%2)+1)
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("message no. " + strconv.Itoa(i)),
		}
		err := channelRMQ.Publish(
			"",
			targetQueue,
			false,
			false,
			message,
		)
		if err != nil {
			panic(err)
		}
	}

}
