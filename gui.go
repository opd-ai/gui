// Package gui provides a cross-platform GUI drawing library in pure Go
package gui

import (
	"image"
	"sync"

	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
)

// GUIElement defines the core interface that all GUI components must implement
type GUIElement interface {
	// Render draws the element to the provided canvas
	Render(canvas Canvas) error

	// GetBounds returns the element's position and dimensions
	GetBounds() (x, y, width, height int)

	// SetPosition updates the element's position
	SetPosition(x, y int)

	// IsVisible returns the element's visibility state
	IsVisible() bool

	// SetVisible controls element visibility
	SetVisible(visible bool)

	// HandleEvent processes input events
	HandleEvent(event Event) bool
}

// Canvas provides drawing operations abstraction
type Canvas interface {
	// Text rendering
	DrawText(text string, x, y int, font font.Face, color colorful.Color) error

	// Shape primitives
	DrawRectangle(x, y, width, height int, color colorful.Color, filled bool) error
	DrawCircle(x, y, radius int, color colorful.Color, filled bool) error

	// Image operations
	DrawImage(img image.Image, x, y, width, height int) error

	// Clipping and transformations
	SetClippingRegion(x, y, width, height int)
	ClearClippingRegion()

	// Canvas management
	Clear(color colorful.Color) error
	Present() error
}

// Element provides base implementation for GUI elements
type Element struct {
	mu       sync.RWMutex
	x, y     int
	width    int
	height   int
	visible  bool
	parent   GUIElement
	children []GUIElement
	handlers map[EventType][]EventHandler
}

// NewElement creates a new base element
func NewElement(x, y, width, height int) *Element {
	return &Element{
		x:        x,
		y:        y,
		width:    width,
		height:   height,
		visible:  true,
		children: make([]GUIElement, 0),
		handlers: make(map[EventType][]EventHandler),
	}
}

// GetBounds returns the element's bounds
func (e *Element) GetBounds() (x, y, width, height int) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.x, e.y, e.width, e.height
}

// SetPosition updates the element's position
func (e *Element) SetPosition(x, y int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.x, e.y = x, y
}

// SetSize updates the element's dimensions
func (e *Element) SetSize(width, height int) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.width, e.height = width, height
}

// IsVisible returns visibility state
func (e *Element) IsVisible() bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.visible
}

// SetVisible controls visibility
func (e *Element) SetVisible(visible bool) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.visible = visible
}

// AddChild adds a child element
func (e *Element) AddChild(child GUIElement) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.children = append(e.children, child)
}

// RemoveChild removes a child element
func (e *Element) RemoveChild(child GUIElement) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for i, c := range e.children {
		if c == child {
			e.children = append(e.children[:i], e.children[i+1:]...)
			break
		}
	}
}

// ContainsPoint checks if a point is within the element's bounds
func (e *Element) ContainsPoint(x, y int) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return x >= e.x && x < e.x+e.width && y >= e.y && y < e.y+e.height
}

// Render provides default rendering for child elements
func (e *Element) Render(canvas Canvas) error {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if !e.visible {
		return nil
	}

	// Render all children
	for _, child := range e.children {
		if err := child.Render(canvas); err != nil {
			return err
		}
	}

	return nil
}

// HandleEvent processes events and propagates to children
func (e *Element) HandleEvent(event Event) bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// Check if event is within bounds for position-based events
	if posEvent, ok := event.(*ClickEvent); ok {
		if !e.ContainsPoint(posEvent.X, posEvent.Y) {
			return false
		}
	}

	// Try children first (reverse order for proper z-order)
	for i := len(e.children) - 1; i >= 0; i-- {
		if e.children[i].HandleEvent(event) {
			return true
		}
	}

	// Handle at this level
	if handlers, exists := e.handlers[event.Type()]; exists {
		for _, handler := range handlers {
			if handler.Handle(event) {
				return true
			}
		}
	}

	return false
}

// AddEventHandler registers an event handler
func (e *Element) AddEventHandler(eventType EventType, handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// Window represents the main application window
type Window struct {
	*Element
	title    string
	canvas   Canvas
	renderer Renderer
	running  bool
	mu       sync.RWMutex
}

// NewWindow creates a new application window
func NewWindow(title string, width, height int) (*Window, error) {
	renderer, err := NewRenderer(width, height)
	if err != nil {
		return nil, err
	}

	canvas, err := renderer.CreateCanvas()
	if err != nil {
		return nil, err
	}

	return &Window{
		Element:  NewElement(0, 0, width, height),
		title:    title,
		canvas:   canvas,
		renderer: renderer,
		running:  false,
	}, nil
}

// Show displays the window
func (w *Window) Show() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.running = true
	return w.renderer.Show(w.title)
}

// Close closes the window
func (w *Window) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.running = false
	return w.renderer.Close()
}

// IsRunning returns whether the window is active
func (w *Window) IsRunning() bool {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return w.running
}

// Update renders the window contents
func (w *Window) Update() error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	if !w.running {
		return nil
	}

	// Clear canvas
	w.canvas.Clear(colorful.Color{R: 1.0, G: 1.0, B: 1.0})

	// Render all elements
	if err := w.Element.Render(w.canvas); err != nil {
		return err
	}

	// Present to screen
	return w.canvas.Present()
}

// PollEvents processes pending events
func (w *Window) PollEvents() []Event {
	return w.renderer.PollEvents()
}
