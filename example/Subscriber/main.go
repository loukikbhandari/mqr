package main

import (
	"encoding/json"
	"fmt"
	"github.com/loukikbhandari/mqr"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	prefetchLimit = 1000
	pollDuration  = 100 * time.Millisecond
	numConsumers  = 5

	//reportBatchSize = 10000
	consumeDuration = time.Millisecond
	shouldLog       = true
)

func main() {
	errChan := make(chan error, 10)
	go logErrors(errChan)

	testStrArr := []string{"x1", "x2", "x3", "x4", "x5"}

	connection, err := mqr.OpenConnection("Subscriber", "tcp", "localhost:6379", 2, errChan)
	if err != nil {
		panic(err)
	}

	queue, err := connection.OpenQueue("bosch-test-queue")
	if err != nil {
		panic(err)
	}

	if err := queue.StartConsuming(prefetchLimit, pollDuration); err != nil {
		panic(err)
	}

	for i := 0; i < numConsumers; i++ {
		name := fmt.Sprintf("Subscriber %d", i)
		if _, err := queue.AddConsumer(name, NewConsumer(i, testStrArr[i])); err != nil {
			panic(err)
		}
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	defer signal.Stop(signals)

	<-signals // wait for signal
	go func() {
		<-signals // hard exit on second signal (in case shutdown gets stuck)
		os.Exit(1)
	}()

	<-connection.StopAllConsuming() // wait for all Consume() calls to finish
}

type Consumer struct {
	name    string
	count   int
	before  time.Time
	testStr string
}

func NewConsumer(tag int, testStr string) *Consumer {
	return &Consumer{
		name:    fmt.Sprintf("Subscriber%d", tag),
		count:   0,
		before:  time.Now(),
		testStr: testStr,
	}
}

func (consumer *Consumer) Consume(delivery mqr.Delivery) {
	payload := delivery.Payload()
	debugf("start consume %s", payload)
	time.Sleep(consumeDuration)

	consumer.count++
	var msg mqr.Message

	if err := msg.UnmarshalBinary([]byte(payload)); err != nil {
		// handle json error

		if err := delivery.Reject(); err != nil {
			// handle reject error
		}
		return
	}

	var jsonBody mqr.JsonBody

	err := json.Unmarshal([]byte(msg.Body), &jsonBody)
	if err != nil {
		log.Printf("jsonBody = %v", jsonBody)
	}

	if jsonBody.City == consumer.testStr {

		if err := delivery.Ack(); err != nil {
			// handle ack error
		}
	} else {
		log.Printf("Rejected message with body jsonBody = %v by subscriber %s ", jsonBody, consumer.name)
		if err := delivery.Reject(); err != nil {
			// handle reject error
		}
	}

}

func logErrors(errChan <-chan error) {
	for err := range errChan {
		switch err := err.(type) {
		case *mqr.HeartbeatError:
			if err.Count == 10 {
				log.Print("heartbeat error (limit): ", err)
			} else {
				log.Print("heartbeat error: ", err)
			}
		case *mqr.ConsumeError:
			log.Print("consume error: ", err)
		case *mqr.DeliveryError:
			log.Print("delivery error: ", err.Delivery, err)
		default:
			log.Print("other error: ", err)
		}
	}
}

func debugf(format string, args ...interface{}) {
	if shouldLog {
		log.Printf(format, args...)
	}
}
