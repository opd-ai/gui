# GUI - Pure Go Cross-Platform GUI Library

A lightweight, cross-platform GUI drawing library implemented in pure Go, providing essential components for building graphical user interfaces with concurrent-safe operations and event handling.

---

## Installation

This package requires Go 1.19 or later. Install using:

```bash
go get github.com/opd-ai/gui
```

### Dependencies

The package relies on the following external libraries:

```bash
go get github.com/fogleman/gg@v1.3.0
go get github.com/lucasb-eyer/go-colorful@v1.2.0
go get golang.org/x/image@v0.13.0
go get golang.org/x/text@v0.13.0
```

---

## Usage

### Basic Window Creation

```go
package main

import (
    "github.com/opd-ai/gui"
    "github.com/lucasb-eyer/go-colorful"
)

func main() {
    // Create a new window
    window, err := gui.NewWindow("My Application", 800, 600)
    if err != nil {
        panic(err)
    }
    defer window.Close()

    // Show the window
    if err := window.Show(); err != nil {
        panic(err)
    }

    // Main event loop
    for window.IsRunning() {
        events := window.PollEvents()
        for _, event := range events {
            window.HandleEvent(event)
        }
        
        if err := window.Update(); err != nil {
            panic(err)
        }
    }
}
```

### Creating GUI Elements

```go
// Create a base element
element := gui.NewElement(10, 10, 200, 100)
element.SetVisible(true)

// Add to window
window.AddChild(element)

// Position and size management
element.SetPosition(50, 50)
element.SetSize(300, 150)

// Check bounds
x, y, width, height := element.GetBounds()
```

### Event Handling

```go
// Add event handler to element
element.AddEventHandler(gui.ClickEventType, &CustomClickHandler{})

// Check if point is within element
if element.ContainsPoint(mouseX, mouseY) {
    // Handle mouse interaction
}
```

---

## Features

- **Cross-Platform Compatibility** - Pure Go implementation for platform independence
- **Concurrent-Safe Operations** - Thread-safe element management with mutex protection
- **Hierarchical Element System** - Parent-child relationships with automatic rendering
- **Event Handling Framework** - Extensible event system with custom handler support
- **Canvas Drawing Operations** - Text rendering, shapes, images, and clipping regions
- **Window Management** - Complete window lifecycle with show/hide/close operations
- **Visibility Control** - Individual element visibility management
- **Bounds Detection** - Point containment and collision detection
- **Color Support** - Integration with go-colorful for advanced color operations
- **Font Rendering** - Text drawing with customizable fonts and positioning

### Core Interfaces

- **GUIElement** - Base interface for all GUI components with rendering and event handling
- **Canvas** - Drawing abstraction supporting text, shapes, images, and transformations
- **Element** - Concrete implementation with child management and event propagation
- **Window** - Main application window with renderer integration

### Drawing Capabilities

- Rectangle and circle primitives (filled and outlined)
- Text rendering with font face support
- Image drawing with scaling and positioning
- Clipping region management
- Canvas clearing and presentation

---

## Requirements

- Go 1.19 or later
- Compatible with golang.org/x/image and golang.org/x/text packages
- Requires fogleman/gg for graphics operations
- Uses go-colorful for color management
