package main

import (
	"github.com/google/uuid"
	"github.com/loukikbhandari/mqr"
	"log"
	"time"
)

var TestData = []mqr.Message{
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x1", "state":"y2"}`,
	},
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x2", "state":"y2"}`,
	},
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x3", "state":"y3"}`,
	},
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x4", "state":"y4"}`,
	},
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x5", "state":"y5"}`,
	},
	{
		Id:   uuid.NewString(),
		Time: time.Time{},
		Body: `{"city":"x6", "state":"y6"}`,
	},
}

func main() {
	connection, err := mqr.OpenConnection("Publisher", "tcp", "localhost:6379", 2, nil)
	if err != nil {
		panic(err)
	}

	bosch_test_queue, err := connection.OpenQueue("bosch-test-queue")
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(TestData); i++ {

		TestData[i].Time = time.Now()
		delivery, err := TestData[i].MarshalBinary()

		if err != nil {
			log.Fatalf("cqannot marshal into binary")
		}

		if err := bosch_test_queue.PublishBytes(delivery); err != nil {
			log.Printf("failed to publish: %s", err)
		}

		log.Printf("\n Published %v ", string(delivery))

	}
}
