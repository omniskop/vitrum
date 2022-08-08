package vit

import "fmt"

// EventDefinition from a vit file
type EventDefinition struct {
	Name       string
	Parameters []PropertyDefinition
	Position   *PositionRange
}

// Evaluaters are structures that can be evaluated inside of a component context.
type Evaluater interface {
	ShouldEvaluate() bool                    // Returns true if the structure needs to be reevaluated. (For example because a dependency has changed.)
	Evaluate(Component) (interface{}, error) // Evaluates the structure and returns the result and an optional error.
}

// Listener provides an interface specific to an EventType.
type Listener[EventType any] interface {
	Notify(*EventType) // Notify the listener that the event has been triggered.
}

// EventListenerFunction is a function that takes a single specific argument.
type EventListenerFunction[EventType any] struct {
	Function
	event *EventType
	dirty bool
}

// Creates a new EventListenerFunction using the specific code.
func NewEventListenerFunction[EventType any](code Code) *EventListenerFunction[EventType] {
	return &EventListenerFunction[EventType]{
		Function: *NewFunctionFromCode(code),
		event:    nil,
		dirty:    false, // async functions start clean
	}
}

// Notify the function that the event it is listening to has been triggered.
// It satisfies the Listener interface for the specific EventType.
func (f *EventListenerFunction[EventType]) Notify(event *EventType) {
	f.dirty = true
	f.event = event
}

// Evaluate the function.
// Satisfies the Evaluater interface.
func (f *EventListenerFunction[EventType]) Evaluate(context Component) (interface{}, error) {
	f.dirty = false
	return f.Call(context, f.event)
}

// ShouldEvaluate returns true if the event this function is listening to has been trigged since the last
// time the function was executed.
// Satisfies the Evaluater interface.
func (f *EventListenerFunction[EventType]) ShouldEvaluate() bool {
	return f.dirty
}

// ListenerCallback is an adapter that implements the Dependent interface
// and simply calls the callback function when the dependency changes.
type ListenerCallback[EventType any] struct {
	Callback *func(*EventType)
}

func ListenerCB[EventType any](cb func(*EventType)) ListenerCallback[EventType] {
	return ListenerCallback[EventType]{
		Callback: &cb,
	}
}

func (d ListenerCallback[EventType]) Notify(e *EventType) {
	(*d.Callback)(e)
}

type JSListener[EventType any] struct {
	function *AsyncFunction
}

func (l *JSListener[EventType]) Notify(event *EventType) {
	l.function.dirty = true
	l.function.arguments = []interface{}{event}
}

// Listenable provides a universal interface for all events that can be listened for, no matter the EventType.
type Listenable interface {
	// Creates a new listener for the event using the code. The listener is registered with this event and is returned in the form of an Evaluater.
	CreateListener(Code) Evaluater
	// Adds an async function as a listener
	AddListenerFunction(*AsyncFunction)
}

// Can be implemented by Events to enable vitrum to automatically create an event from a javascript value.
type MaybeSetable interface {
	MaybeSet(any) error // Tries to set the implementing struct based on the passed value. Returns an error if the value is not valid.
}

// EventSource provides a universal interface for EventAttributes without the generic event type.
type EventSource interface {
	Listenable
	MaybeFire(any) error // Tries to fire an event using the passed value. Returns an error if the value can't be converted to the correct event type.
}

// an EventAttribute that is defined on a component.
type EventAttribute[EventType any] struct {
	Listeners map[Listener[EventType]]bool
}

func NewEventAttribute[EventType any]() *EventAttribute[EventType] {
	return &EventAttribute[EventType]{
		Listeners: make(map[Listener[EventType]]bool),
	}
}

func (a *EventAttribute[EventType]) AddListener(l Listener[EventType]) {
	a.Listeners[l] = true
}

func (a *EventAttribute[EventType]) CreateListener(code Code) Evaluater {
	l := NewEventListenerFunction[EventType](code)
	a.Listeners[l] = true
	return l
}

func (a *EventAttribute[EventType]) AddListenerFunction(f *AsyncFunction) {
	l := &JSListener[EventType]{f}
	a.Listeners[l] = true
}

func (a *EventAttribute[EventType]) RemoveListener(l Listener[EventType]) {
	delete(a.Listeners, l)
}

func (a *EventAttribute[EventType]) Fire(e *EventType) {
	for l := range a.Listeners {
		l.Notify(e)
	}
}

// MaybeFire tries to convert the value to the correct EventType and if that's possible fires the event.
// The event type has to implement the MaybeSetable interface.
// If the value can't be converted the event will not be fired and an error will be returned.
func (a *EventAttribute[EventType]) MaybeFire(v any) error {
	// TODO: try to convert ordinary values via goja?
	var event = new(EventType)
	if setable, ok := any(event).(MaybeSetable); ok {
		err := setable.MaybeSet(v)
		if err != nil {
			return err
		}
		a.Fire(event)
		return nil
	}
	return fmt.Errorf("event %T can't be set from interface{}", new(EventType))
}

func (a *EventAttribute[EventType]) eventType() any {
	return new(EventType)
}
