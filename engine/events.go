package engine

/*
The eventsystem has ordered handlers, which means that if one handler takes long to execute, the rest will have to wait.
Potentially long event handlers should be ran in goroutines.
*/

import (
	"context"
	"github.com/Sirupsen/logrus"
	"sync"
)

const (
	AddMultipleWildcardWarning = "A Plugin tried to add listener for multiple events when one was wildcard. That's stupid, what are you doing?!"
)

type HandlerFunc func(data *EventData)

var (
	eventHandlers     = make(map[string][]*HandlerFunc)
	eventHandlersLock sync.RWMutex
)

// EmitEvent emits an event
func EmitEvent(name string, data *EventData) {
	if name != "*" {
		go EmitEvent("*", data)
	}

	eventHandlersLock.RLock()
	handlers := eventHandlers[name]
	if len(handlers) < 1 {
		// No handler for event, don't do anything
		eventHandlersLock.RUnlock()
		return
	}

	// Make a copy of the handlers so that we can unlock the main eventhandlerslock
	cop := make([]*HandlerFunc, len(handlers))
	copy(cop, handlers)

	eventHandlersLock.RUnlock()

	for _, v := range handlers {
		(*v)(data)
	}
}

// AddHandler adds one or more handlers and returns a pointer to it
func AddHandler(handler HandlerFunc, evts ...string) *HandlerFunc {
	for _, evt := range evts {
		if evt == "*" && len(evts) > 1 {
			logrus.Warn(AddMultipleWildcardWarning)
			return AddHandler(handler, "*")
		}
	}

	hPtr := &handler

	eventHandlersLock.Lock()
	defer eventHandlersLock.Unlock()

	for _, evt := range evts {
		eventHandlers[evt] = append(eventHandlers[evt], hPtr)
	}

	return hPtr
}

// AddHandlerFirst adds a handler first in the queue and returns a pointer to the handler
func AddHandlerFirst(handler HandlerFunc, evts ...string) *HandlerFunc {
	hPtr := &handler

	for _, evt := range evts {
		if evt == "*" && len(evts) > 1 {
			logrus.Warn(AddMultipleWildcardWarning)
			return AddHandlerFirst(handler, "*")
		}
	}

	eventHandlersLock.Lock()
	defer eventHandlersLock.Unlock()
	for _, v := range evts {
		eventHandlers[v] = append([]*HandlerFunc{hPtr}, eventHandlers[v]...)
	}

	return hPtr
}

// AddHandlerBefore adds a handler to be called before another handler and returns a pointer to the handler
func AddHandlerBefore(handler HandlerFunc, before *HandlerFunc, evts ...string) *HandlerFunc {
	for _, evt := range evts {
		if evt == "*" && len(evts) > 1 {
			logrus.Warn(AddMultipleWildcardWarning)
			return AddHandlerBefore(handler, before, "*")
		}
	}

	hPtr := &handler

	eventHandlersLock.Lock()
	defer eventHandlersLock.Unlock()

OUTER:
	for _, evt := range evts {
		for k, v := range eventHandlers[evt] {
			if v == before {
				// Make a copy with the first half in
				handlersCop := make([]*HandlerFunc, len(eventHandlers[evt])+1)
				copy(handlersCop, eventHandlers[evt][:k])

				// insert the handler
				handlersCop[k] = hPtr

				// add the other half
				for i := k; i < len(eventHandlers[evt]); i++ {
					handlersCop[i+1] = eventHandlers[evt][i]
				}

				eventHandlers[evt] = handlersCop

				continue OUTER
			}
		}

		// Not found, just add to end
		logrus.Error("Unable to add handler before other handler, handler:", handler, ", Before:", before, ", adding to end instead...")
		eventHandlers[evt] = append(eventHandlers[evt], hPtr)
	}

	return hPtr
}

func NumHandlers(evt string) int {

	eventHandlersLock.RLock()
	defer eventHandlersLock.RUnlock()

	if evt == "" {
		total := 0
		for _, v := range eventHandlers {
			total += len(v)
		}
		return total
	}

	return len(eventHandlers[evt])
}

type EventData struct {
	// The raw event data
	Evt interface{}

	// Name of the event
	Name string

	// Context
	ctx context.Context
}

func NewEventData(name string, evt interface{}, ctx context.Context) *EventData {
	return &EventData{
		Name: name,
		Evt:  evt,
		ctx:  ctx,
	}
}

func (e *EventData) Context() context.Context {
	if e.ctx == nil {
		return context.Background()
	}

	return e.ctx
}

func (e *EventData) WithContext(ctx context.Context) *EventData {
	cop := new(EventData)
	*cop = *e
	cop.ctx = ctx
	return cop
}
