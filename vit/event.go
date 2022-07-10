package vit

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
func NewEventListenerFunction[EventType any](code string, position *PositionRange) *EventListenerFunction[EventType] {
	return &EventListenerFunction[EventType]{
		Function: *NewFunction(code, position),
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
	CreateListener(string, *PositionRange) Evaluater
	// Adds an async function as a listener
	AddListenerFunction(*AsyncFunction)
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

func (a *EventAttribute[EventType]) CreateListener(code string, position *PositionRange) Evaluater {
	l := NewEventListenerFunction[EventType](code, position)
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
