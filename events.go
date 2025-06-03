package gui

import "time"

// EventType represents different types of GUI events
type EventType int

const (
	EventTypeClick EventType = iota
	EventTypeKeyPress
	EventTypeTextInput
	EventTypeFocus
	EventTypeBlur
	EventTypeMouseMove
	EventTypeResize
)

// Event defines the interface for all GUI events
type Event interface {
	Type() EventType
	Timestamp() time.Time
}

// EventHandler processes GUI events
type EventHandler interface {
	Handle(event Event) bool
}

// EventHandlerFunc adapts functions to EventHandler interface
type EventHandlerFunc func(event Event) bool

func (f EventHandlerFunc) Handle(event Event) bool {
	return f(event)
}

// BaseEvent provides common event functionality
type BaseEvent struct {
	eventType EventType
	timestamp time.Time
}

func NewBaseEvent(eventType EventType) BaseEvent {
	return BaseEvent{
		eventType: eventType,
		timestamp: time.Now(),
	}
}

func (e BaseEvent) Type() EventType {
	return e.eventType
}

func (e BaseEvent) Timestamp() time.Time {
	return e.timestamp
}

// ClickEvent represents mouse click events
type ClickEvent struct {
	BaseEvent
	X, Y   int
	Button MouseButton
}

type MouseButton int

const (
	MouseButtonLeft MouseButton = iota
	MouseButtonRight
	MouseButtonMiddle
)

func NewClickEvent(x, y int, button MouseButton) *ClickEvent {
	return &ClickEvent{
		BaseEvent: NewBaseEvent(EventTypeClick),
		X:         x,
		Y:         y,
		Button:    button,
	}
}

// KeyPressEvent represents keyboard input
type KeyPressEvent struct {
	BaseEvent
	Key       Key
	Modifiers KeyModifiers
}

type Key int

const (
	KeyUnknown Key = iota
	KeyA
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	Key0
	Key1
	Key2
	Key3
	Key4
	Key5
	Key6
	Key7
	Key8
	Key9
	KeySpace
	KeyEnter
	KeyTab
	KeyBackspace
	KeyDelete
	KeyEscape
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight
)

type KeyModifiers int

const (
	ModifierNone  KeyModifiers = 0
	ModifierShift KeyModifiers = 1 << iota
	ModifierCtrl
	ModifierAlt
	ModifierSuper
)

func NewKeyPressEvent(key Key, modifiers KeyModifiers) *KeyPressEvent {
	return &KeyPressEvent{
		BaseEvent: NewBaseEvent(EventTypeKeyPress),
		Key:       key,
		Modifiers: modifiers,
	}
}

// TextInputEvent represents text input
type TextInputEvent struct {
	BaseEvent
	Text string
}

func NewTextInputEvent(text string) *TextInputEvent {
	return &TextInputEvent{
		BaseEvent: NewBaseEvent(EventTypeTextInput),
		Text:      text,
	}
}

// FocusEvent represents focus gain
type FocusEvent struct {
	BaseEvent
}

func NewFocusEvent() *FocusEvent {
	return &FocusEvent{
		BaseEvent: NewBaseEvent(EventTypeFocus),
	}
}

// BlurEvent represents focus loss
type BlurEvent struct {
	BaseEvent
}

func NewBlurEvent() *BlurEvent {
	return &BlurEvent{
		BaseEvent: NewBaseEvent(EventTypeBlur),
	}
}

// MouseMoveEvent represents mouse movement
type MouseMoveEvent struct {
	BaseEvent
	X, Y int
}

func NewMouseMoveEvent(x, y int) *MouseMoveEvent {
	return &MouseMoveEvent{
		BaseEvent: NewBaseEvent(EventTypeMouseMove),
		X:         x,
		Y:         y,
	}
}

// ResizeEvent represents window resize
type ResizeEvent struct {
	BaseEvent
	Width, Height int
}

func NewResizeEvent(width, height int) *ResizeEvent {
	return &ResizeEvent{
		BaseEvent: NewBaseEvent(EventTypeResize),
		Width:     width,
		Height:    height,
	}
}
