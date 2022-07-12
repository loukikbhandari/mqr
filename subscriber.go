package mqr

type Subscriber interface {
	Consume(delivery Delivery)
}

type ConsumerFunc func(Delivery)

func (consumerFunc ConsumerFunc) Consume(delivery Delivery) {
	consumerFunc(delivery)
}
