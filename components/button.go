package components

import (
	"image"

	"github.com/opd-ai/gui"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
)

// ButtonState represents the visual state of a button
type ButtonState int

const (
	ButtonStateNormal ButtonState = iota
	ButtonStateHover
	ButtonStatePressed
	ButtonStateDisabled
)

// Button represents a clickable button component
type Button struct {
	*gui.Element
	text    string
	icon    image.Image
	font    font.Face
	state   ButtonState
	enabled bool

	// Colors for different states
	normalBgColor   colorful.Color
	hoverBgColor    colorful.Color
	pressedBgColor  colorful.Color
	disabledBgColor colorful.Color

	textColor         colorful.Color
	disabledTextColor colorful.Color
	borderColor       colorful.Color

	// Event callbacks
	onClick   func()
	onHover   func()
	onUnhover func()

	// Visual properties
	borderWidth  int
	cornerRadius int
}

// NewButton creates a new button with the specified text
func NewButton(text string) *Button {
	button := &Button{
		Element: gui.NewElement(0, 0, 100, 30),
		text:    text,
		font:    basicfont.Face7x13,
		state:   ButtonStateNormal,
		enabled: true,

		// Default colors
		normalBgColor:   colorful.Color{R: 0.9, G: 0.9, B: 0.9},    // Light gray
		hoverBgColor:    colorful.Color{R: 0.8, G: 0.8, B: 0.9},    // Light blue
		pressedBgColor:  colorful.Color{R: 0.7, G: 0.7, B: 0.8},    // Darker blue
		disabledBgColor: colorful.Color{R: 0.95, G: 0.95, B: 0.95}, // Very light gray

		textColor:         colorful.Color{R: 0, G: 0, B: 0},       // Black
		disabledTextColor: colorful.Color{R: 0.6, G: 0.6, B: 0.6}, // Gray
		borderColor:       colorful.Color{R: 0.6, G: 0.6, B: 0.6}, // Gray

		borderWidth:  1,
		cornerRadius: 3,
	}

	// Register event handlers
	button.AddEventHandler(gui.EventTypeClick, gui.EventHandlerFunc(button.handleClick))
	button.AddEventHandler(gui.EventTypeMouseMove, gui.EventHandlerFunc(button.handleMouseMove))

	return button
}

// SetText updates the button's text
func (b *Button) SetText(text string) *Button {
	b.text = text
	return b
}

// GetText returns the button's text
func (b *Button) GetText() string {
	return b.text
}

// SetIcon sets an icon image for the button
func (b *Button) SetIcon(icon image.Image) *Button {
	b.icon = icon
	return b
}

// SetFont updates the button's font
func (b *Button) SetFont(font font.Face) *Button {
	b.font = font
	return b
}

// SetEnabled controls whether the button can be clicked
func (b *Button) SetEnabled(enabled bool) *Button {
	b.enabled = enabled
	if !enabled {
		b.state = ButtonStateDisabled
	} else if b.state == ButtonStateDisabled {
		b.state = ButtonStateNormal
	}
	return b
}

// IsEnabled returns whether the button is enabled
func (b *Button) IsEnabled() bool {
	return b.enabled
}

// SetNormalColor sets the normal background color
func (b *Button) SetNormalColor(color colorful.Color) *Button {
	b.normalBgColor = color
	return b
}

// SetHoverColor sets the hover background color
func (b *Button) SetHoverColor(color colorful.Color) *Button {
	b.hoverBgColor = color
	return b
}

// SetPressedColor sets the pressed background color
func (b *Button) SetPressedColor(color colorful.Color) *Button {
	b.pressedBgColor = color
	return b
}

// SetTextColor sets the text color
func (b *Button) SetTextColor(color colorful.Color) *Button {
	b.textColor = color
	return b
}

// SetBorderColor sets the border color
func (b *Button) SetBorderColor(color colorful.Color) *Button {
	b.borderColor = color
	return b
}

// SetOnClick sets the click event callback
func (b *Button) SetOnClick(callback func()) *Button {
	b.onClick = callback
	return b
}

// SetOnHover sets the hover event callback
func (b *Button) SetOnHover(callback func()) *Button {
	b.onHover = callback
	return b
}

// SetOnUnhover sets the unhover event callback
func (b *Button) SetOnUnhover(callback func()) *Button {
	b.onUnhover = callback
	return b
}

// handleClick processes mouse click events
func (b *Button) handleClick(event gui.Event) bool {
	if !b.enabled {
		return false
	}

	clickEvent := event.(*gui.ClickEvent)

	if !b.ContainsPoint(clickEvent.X, clickEvent.Y) {
		return false
	}

	// Visual feedback
	b.state = ButtonStatePressed

	// Execute callback
	if b.onClick != nil {
		b.onClick()
	}

	// Reset state after a brief moment (in a real implementation,
	// this would be handled by the event loop timing)
	b.state = ButtonStateNormal

	return true
}

// handleMouseMove processes mouse movement for hover effects
func (b *Button) handleMouseMove(event gui.Event) bool {
	if !b.enabled {
		return false
	}

	moveEvent := event.(*gui.MouseMoveEvent)
	wasHover := (b.state == ButtonStateHover)
	isHover := b.ContainsPoint(moveEvent.X, moveEvent.Y)

	if isHover && !wasHover {
		b.state = ButtonStateHover
		if b.onHover != nil {
			b.onHover()
		}
		return true
	} else if !isHover && wasHover {
		b.state = ButtonStateNormal
		if b.onUnhover != nil {
			b.onUnhover()
		}
		return true
	}

	return false
}

// getCurrentBackgroundColor returns the appropriate background color for the current state
func (b *Button) getCurrentBackgroundColor() colorful.Color {
	switch b.state {
	case ButtonStateHover:
		return b.hoverBgColor
	case ButtonStatePressed:
		return b.pressedBgColor
	case ButtonStateDisabled:
		return b.disabledBgColor
	default:
		return b.normalBgColor
	}
}

// getCurrentTextColor returns the appropriate text color for the current state
func (b *Button) getCurrentTextColor() colorful.Color {
	if b.state == ButtonStateDisabled {
		return b.disabledTextColor
	}
	return b.textColor
}

// Render draws the button component
func (b *Button) Render(canvas gui.Canvas) error {
	if !b.IsVisible() {
		return nil
	}

	x, y, width, height := b.GetBounds()

	// Draw background with current state color
	bgColor := b.getCurrentBackgroundColor()
	if err := canvas.DrawRectangle(x, y, width, height, bgColor, true); err != nil {
		return err
	}

	// Draw border
	if b.borderWidth > 0 {
		if err := canvas.DrawRectangle(x, y, width, height, b.borderColor, false); err != nil {
			return err
		}
	}

	// Calculate content positioning
	contentX := x
	contentY := y
	contentWidth := width
	contentHeight := height

	// Account for padding
	padding := 8
	contentX += padding
	contentY += padding
	contentWidth -= 2 * padding
	contentHeight -= 2 * padding

	// Draw icon if present
	iconWidth := 0
	if b.icon != nil {
		iconSize := contentHeight
		if iconSize > contentWidth/3 {
			iconSize = contentWidth / 3
		}

		iconY := contentY + (contentHeight-iconSize)/2
		if err := canvas.DrawImage(b.icon, contentX, iconY, iconSize, iconSize); err != nil {
			return err
		}

		iconWidth = iconSize + 4 // Icon + spacing
	}

	// Draw text
	if b.text != "" {
		textColor := b.getCurrentTextColor()

		// Calculate text position (centered)
		textX := contentX + iconWidth
		textWidth := contentWidth - iconWidth

		// Simple text width calculation (approximate)
		if b.font != nil && textWidth > 0 {
			// Center text horizontally in remaining space
			estimatedTextWidth := len(b.text) * 7 // Rough estimate
			if estimatedTextWidth < textWidth {
				textX += (textWidth - estimatedTextWidth) / 2
			}
		}

		// Center text vertically
		textY := contentY + contentHeight/2 + 4 // Adjust for baseline

		if err := canvas.DrawText(b.text, textX, textY, b.font, textColor); err != nil {
			return err
		}
	}

	return nil
}
