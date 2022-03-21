package core

import bus "github.com/asaskevich/EventBus"

const (
	EventExit = "exit"
)

type EventBus = bus.Bus
type BusSubscriber = bus.BusSubscriber

var eventBus EventBus

func PublishEvent(topic string, args ...interface{}) {
	eventBus.Publish(topic, args...)
}

func SubscribeEvent(topic string, fn interface{}) {
	err := eventBus.Subscribe(topic, fn)
	if err != nil {
		panic(err)
	}
}

func SubscribeEventOnce(topic string, fn interface{}) {
	err := eventBus.SubscribeOnce(topic, fn)
	if err != nil {
		panic(err)
	}
}

func init() {
	eventBus = bus.New()
}
