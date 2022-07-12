package mqr

import (
	"context"
	"fmt"
	"time"
)

type Delivery interface {
	Payload() string
	Ack() error
	Reject() error
	Push() error
}

type redisDelivery struct {
	ctx         context.Context
	payload     string
	unackedKey  string
	rejectedKey string
	pushKey     string
	redisClient RedisClient
	errChan     chan<- error
}

func newDelivery(
	ctx context.Context,
	payload string,
	unackedKey string,
	rejectedKey string,
	pushKey string,
	redisClient RedisClient,
	errChan chan<- error,
) *redisDelivery {
	return &redisDelivery{
		ctx:         ctx,
		payload:     payload,
		unackedKey:  unackedKey,
		rejectedKey: rejectedKey,
		pushKey:     pushKey,
		redisClient: redisClient,
		errChan:     errChan,
	}
}

func (delivery *redisDelivery) String() string {
	return fmt.Sprintf("[%s %s]", delivery.payload, delivery.unackedKey)
}

func (delivery *redisDelivery) Payload() string {
	return delivery.payload
}

func (delivery *redisDelivery) Ack() error {
	errorCount := 0
	for {
		count, err := delivery.redisClient.LRem(delivery.unackedKey, 1, delivery.payload)
		if err == nil { // no redis error
			if count == 0 {
				return ErrorNotFound
			}
			return nil
		}

		// redis error

		errorCount++

		select { // try to add error to channel, but don't block
		case delivery.errChan <- &DeliveryError{Delivery: delivery, RedisErr: err, Count: errorCount}:
		default:
		}

		if err := delivery.ctx.Err(); err != nil {
			return ErrorConsumingStopped
		}

		time.Sleep(time.Second)
	}
}

func (delivery *redisDelivery) Reject() error {
	return delivery.move(delivery.rejectedKey)
}

func (delivery *redisDelivery) Push() error {
	if delivery.pushKey == "" {
		return delivery.Reject() // fall back to rejecting
	}

	return delivery.move(delivery.pushKey)
}

func (delivery *redisDelivery) move(key string) error {
	errorCount := 0
	for {
		_, err := delivery.redisClient.LPush(key, delivery.payload)
		if err == nil { // success
			break
		}
		// error

		errorCount++

		select { // try to add error to channel, but don't block
		case delivery.errChan <- &DeliveryError{Delivery: delivery, RedisErr: err, Count: errorCount}:
		default:
		}

		if err := delivery.ctx.Err(); err != nil {
			return ErrorConsumingStopped
		}

		time.Sleep(time.Second)
	}

	return delivery.Ack()
}

// lower level functions which don't retry but just return the first error
