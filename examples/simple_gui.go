package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gui"
	"github.com/gui/components"
	"github.com/lucasb-eyer/go-colorful"
)

func main() {
	// Create main window
	window, err := gui.NewWindow("Simple GUI Demo", 400, 300)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	defer window.Close()

	// Create and configure components
	setupComponents(window)

	// Show window
	if err := window.Show(); err != nil {
		log.Fatalf("Failed to show window: %v", err)
	}

	// Main event loop
	ticker := time.NewTicker(16 * time.Millisecond) // ~60 FPS
	defer ticker.Stop()

	for window.IsRunning() {
		select {
		case <-ticker.C:
			// Process events
			events := window.PollEvents()
			for _, event := range events {
				if !window.HandleEvent(event) {
					// Handle global events (like window close)
					handleGlobalEvent(window, event)
				}
			}

			// Update display
			if err := window.Update(); err != nil {
				log.Printf("Render error: %v", err)
			}
		}
	}
}

func setupComponents(window *gui.Window) {
	// Title label
	titleLabel := components.NewLabel("GUI Library Demo").
		SetColor(colorful.Color{R: 0.2, G: 0.2, B: 0.8}).
		SetAlignment(components.AlignCenter)
	titleLabel.SetPosition(50, 20)
	titleLabel.SetSize(300, 30)
	window.AddChild(titleLabel)

	// Input field with label
	inputLabel := components.NewLabel("Enter your name:")
	inputLabel.SetPosition(50, 70)
	window.AddChild(inputLabel)

	nameInput := components.NewInput().
		SetPlaceholder("Type here...").
		SetMaxLength(50)
	nameInput.SetPosition(50, 95)
	nameInput.SetSize(200, 25)

	// Result label (initially empty)
	resultLabel := components.NewLabel("")
	resultLabel.SetPosition(50, 140)
	resultLabel.SetSize(300, 25)
	window.AddChild(resultLabel)

	// Configure input callbacks
	nameInput.SetOnChange(func(text string) {
		if text != "" {
			resultLabel.SetText(fmt.Sprintf("Hello, %s!", text))
		} else {
			resultLabel.SetText("")
		}
	})

	nameInput.SetOnSubmit(func(text string) {
		fmt.Printf("Submitted: %s\n", text)
	})

	window.AddChild(nameInput)

	// Buttons demonstration
	clickCountLabel := components.NewLabel("Button clicks: 0")
	clickCountLabel.SetPosition(50, 180)
	window.AddChild(clickCountLabel)

	clickCount := 0

	// Click me button
	clickButton := components.NewButton("Click Me!").
		SetNormalColor(colorful.Color{R: 0.8, G: 0.9, B: 0.8}).
		SetHoverColor(colorful.Color{R: 0.7, G: 0.9, B: 0.7}).
		SetTextColor(colorful.Color{R: 0, G: 0.4, B: 0})

	clickButton.SetPosition(50, 210)
	clickButton.SetSize(100, 30)

	clickButton.SetOnClick(func() {
		clickCount++
		clickCountLabel.SetText(fmt.Sprintf("Button clicks: %d", clickCount))
		fmt.Printf("Button clicked! Count: %d\n", clickCount)
	})

	window.AddChild(clickButton)

	// Reset button
	resetButton := components.NewButton("Reset").
		SetNormalColor(colorful.Color{R: 0.9, G: 0.8, B: 0.8}).
		SetHoverColor(colorful.Color{R: 0.9, G: 0.7, B: 0.7}).
		SetTextColor(colorful.Color{R: 0.4, G: 0, B: 0})

	resetButton.SetPosition(160, 210)
	resetButton.SetSize(80, 30)

	resetButton.SetOnClick(func() {
		clickCount = 0
		clickCountLabel.SetText("Button clicks: 0")
		nameInput.SetText("")
		resultLabel.SetText("")
		fmt.Println("Reset clicked!")
	})

	window.AddChild(resetButton)

	// Status label
	statusLabel := components.NewLabel("Ready").
		SetColor(colorful.Color{R: 0.5, G: 0.5, B: 0.5})
	statusLabel.SetPosition(50, 260)
	window.AddChild(statusLabel)

	// Update status based on input focus
	nameInput.SetOnFocus(func() {
		statusLabel.SetText("Typing in name field...")
	})

	nameInput.SetOnBlur(func() {
		statusLabel.SetText("Ready")
	})
}

func handleGlobalEvent(window *gui.Window, event gui.Event) {
	switch event.Type() {
	case gui.EventTypeKeyPress:
		keyEvent := event.(*gui.KeyPressEvent)
		if keyEvent.Key == gui.KeyEscape {
			fmt.Println("Escape pressed, closing window...")
			window.Close()
		}
	}
}
