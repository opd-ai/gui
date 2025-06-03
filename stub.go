//go:build !windows && !darwin && !linux
// +build !windows,!darwin,!linux

package gui

import (
    "fmt"
    "image/png"
    "os"
    "time"

    "github.com/opd-ai/gui/graphics"
)

// StubRenderer provides a basic file-based renderer for testing and development
// This implementation saves rendered frames as PNG files and simulates basic events
type StubRenderer struct {
    width      int
    height     int
    canvas     *graphics.GGCanvas
    running    bool
    frameCount int
    lastFrame  time.Time
    events     []Event
}

// Show displays the window (creates output directory for stub renderer)
func (r *StubRenderer) Show(title string) error {
    if r.running {
        return fmt.Errorf("window already shown")
    }

    // Create output directory for rendered frames
    if err := os.MkdirAll("gui_output", 0755); err != nil && !os.IsExist(err) {
        return fmt.Errorf("failed to create output directory: %w", err)
    }

    r.running = true
    fmt.Printf("GUI Window '%s' opened (%dx%d) - frames saved to gui_output/\n", 
        title, r.width, r.height)
    
    // Generate initial test events
    r.generateTestEvents()
    
    return nil
}

// Close closes the window
func (r *StubRenderer) Close() error {
    if !r.running {
        return nil
    }

    r.running = false
    fmt.Printf("GUI Window closed after %d frames\n", r.frameCount)
    return nil
}

// PollEvents returns pending events (simulated for stub renderer)
func (r *StubRenderer) PollEvents() []Event {
    if !r.running {
        return nil
    }

    // Return queued events and clear the queue
    events := make([]Event, len(r.events))
    copy(events, r.events)
    r.events = r.events[:0] // Clear events

    // Occasionally generate new test events
    if time.Since(r.lastFrame) > 2*time.Second {
        r.generateTestEvents()
        r.lastFrame = time.Now()
    }

    return events
}

// Size returns the window dimensions
func (r *StubRenderer) Size() (width, height int) {
    return r.width, r.height
}

// SetSize updates the window dimensions
func (r *StubRenderer) SetSize(width, height int) error {
    if width <= 0 || height <= 0 {
        return fmt.Errorf("invalid dimensions: width=%d, height=%d", width, height)
    }

    r.width = width
    r.height = height

    // Recreate canvas with new dimensions
    if r.canvas != nil {
        r.canvas = graphics.NewGGCanvas(width, height)
        if r.canvas == nil {
            return fmt.Errorf("failed to recreate canvas with new dimensions")
        }
    }

    return nil
}

// generateTestEvents creates simulated user interactions for testing
func (r *StubRenderer) generateTestEvents() {
    // Simulate some mouse clicks at different positions
    r.events = append(r.events, 
        NewClickEvent(100, 50, MouseButtonLeft),
        NewClickEvent(200, 150, MouseButtonLeft),
        NewMouseMoveEvent(150, 100),
    )

    // Simulate some key presses
    r.events = append(r.events,
        NewKeyPressEvent(KeyA, ModifierNone),
        NewTextInputEvent("Hello"),
    )
}

// Enhanced GGCanvas with frame saving capability
type EnhancedGGCanvas struct {
    *graphics.GGCanvas
    renderer *StubRenderer
}

// Present saves the current frame as a PNG file
func (c *EnhancedGGCanvas) Present() error {
    if c.renderer == nil {
        return fmt.Errorf("no renderer associated with canvas")
    }

    if !c.renderer.running {
        return nil // Don't save frames when window is closed
    }

    // Get the rendered image
    img := c.GetImage()
    if img == nil {
        return fmt.Errorf("no image to present")
    }

    // Save frame as PNG file
    filename := fmt.Sprintf("gui_output/frame_%04d.png", c.renderer.frameCount)
    file, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("failed to create frame file: %w", err)
    }
    defer file.Close()

    if err := png.Encode(file, img); err != nil {
        return fmt.Errorf("failed to encode frame: %w", err)
    }

    c.renderer.frameCount++
    
    // Print progress every 10 frames
    if c.renderer.frameCount%10 == 0 {
        fmt.Printf("Rendered %d frames\n", c.renderer.frameCount)
    }

    return nil
}

// Update CreateCanvas to return enhanced canvas
func (r *StubRenderer) CreateCanvas() (Canvas, error) {
    baseCanvas := graphics.NewGGCanvas(r.width, r.height)
    if baseCanvas == nil {
        return nil, fmt.Errorf("failed to create base canvas")
    }

    enhancedCanvas := &EnhancedGGCanvas{
        GGCanvas: baseCanvas,
        renderer: r,
    }

    r.canvas = baseCanvas
    return enhancedCanvas, nil
}