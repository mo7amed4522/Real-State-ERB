package pubsub

import (
	"sync"
)

type PubSub struct {
	mu   sync.Mutex
	subs map[string][]chan interface{}
}

var (
	ps *PubSub
	once sync.Once
)

func GetInstance() *PubSub {
	once.Do(func() {
		ps = &PubSub{
			subs: make(map[string][]chan interface{}),
		}
	})
	return ps
}

func (ps *PubSub) Subscribe(topic string) <-chan interface{} {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	ch := make(chan interface{}, 1)
	ps.subs[topic] = append(ps.subs[topic], ch)
	return ch
}

func (ps *PubSub) Publish(topic string, data interface{}) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if subs, ok := ps.subs[topic]; ok {
		for _, ch := range subs {
			ch <- data
		}
	}
} 