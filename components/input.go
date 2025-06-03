package components

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/opd-ai/gui"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/text/unicode/norm"
)

// InputValidatorFunc validates input text
type InputValidatorFunc func(text string) bool

// Input represents a text input component
type Input struct {
	*gui.Element
	text             string
	placeholder      string
	font             font.Face
	textColor        colorful.Color
	bgColor          colorful.Color
	borderColor      colorful.Color
	placeholderColor colorful.Color
	cursorPos        int
	selectionStart   int
	selectionEnd     int
	focused          bool
	maxLength        int
	validator        InputValidatorFunc
	onChange         func(text string)
	onSubmit         func(text string)
	onFocus          func()
	onBlur           func()
}

// NewInput creates a new text input component
func NewInput() *Input {
	input := &Input{
		Element:          gui.NewElement(0, 0, 150, 25),
		text:             "",
		placeholder:      "",
		font:             basicfont.Face7x13,
		textColor:        colorful.Color{R: 0, G: 0, B: 0},       // Black
		bgColor:          colorful.Color{R: 1, G: 1, B: 1},       // White
		borderColor:      colorful.Color{R: 0.7, G: 0.7, B: 0.7}, // Gray
		placeholderColor: colorful.Color{R: 0.6, G: 0.6, B: 0.6}, // Light gray
		cursorPos:        0,
		selectionStart:   -1,
		selectionEnd:     -1,
		focused:          false,
		maxLength:        -1, // No limit
	}

	// Register event handlers
	input.AddEventHandler(gui.EventTypeClick, gui.EventHandlerFunc(input.handleClick))
	input.AddEventHandler(gui.EventTypeKeyPress, gui.EventHandlerFunc(input.handleKeyPress))
	input.AddEventHandler(gui.EventTypeTextInput, gui.EventHandlerFunc(input.handleTextInput))
	input.AddEventHandler(gui.EventTypeFocus, gui.EventHandlerFunc(input.handleFocus))
	input.AddEventHandler(gui.EventTypeBlur, gui.EventHandlerFunc(input.handleBlur))

	return input
}

// SetText updates the input's text content
func (i *Input) SetText(text string) *Input {
	// Normalize text for consistent UTF-8 handling
	normalizedText := norm.NFC.String(text)

	// Apply length limit
	if i.maxLength > 0 && utf8.RuneCountInString(normalizedText) > i.maxLength {
		runes := []rune(normalizedText)
		normalizedText = string(runes[:i.maxLength])
	}

	// Apply validation
	if i.validator != nil && !i.validator(normalizedText) {
		return i
	}

	i.text = normalizedText
	i.cursorPos = utf8.RuneCountInString(i.text)
	i.clearSelection()

	if i.onChange != nil {
		i.onChange(i.text)
	}

	return i
}

// GetText returns the current text
func (i *Input) GetText() string {
	return i.text
}

// SetPlaceholder sets placeholder text
func (i *Input) SetPlaceholder(placeholder string) *Input {
	i.placeholder = placeholder
	return i
}

// SetFont updates the input's font
func (i *Input) SetFont(font font.Face) *Input {
	i.font = font
	return i
}

// SetTextColor updates the text color
func (i *Input) SetTextColor(color colorful.Color) *Input {
	i.textColor = color
	return i
}

// SetBackgroundColor updates the background color
func (i *Input) SetBackgroundColor(color colorful.Color) *Input {
	i.bgColor = color
	return i
}

// SetBorderColor updates the border color
func (i *Input) SetBorderColor(color colorful.Color) *Input {
	i.borderColor = color
	return i
}

// SetMaxLength sets maximum character limit
func (i *Input) SetMaxLength(length int) *Input {
	i.maxLength = length

	// Truncate existing text if necessary
	if length > 0 && utf8.RuneCountInString(i.text) > length {
		runes := []rune(i.text)
		i.text = string(runes[:length])
		if i.cursorPos > length {
			i.cursorPos = length
		}
	}

	return i
}

// SetValidator sets input validation function
func (i *Input) SetValidator(validator InputValidatorFunc) *Input {
	i.validator = validator
	return i
}

// SetOnChange sets the onChange callback
func (i *Input) SetOnChange(callback func(text string)) *Input {
	i.onChange = callback
	return i
}

// SetOnSubmit sets the onSubmit callback (Enter key)
func (i *Input) SetOnSubmit(callback func(text string)) *Input {
	i.onSubmit = callback
	return i
}

// SetOnFocus sets the onFocus callback
func (i *Input) SetOnFocus(callback func()) *Input {
	i.onFocus = callback
	return i
}

// SetOnBlur sets the onBlur callback
func (i *Input) SetOnBlur(callback func()) *Input {
	i.onBlur = callback
	return i
}

// Focus gives focus to the input
func (i *Input) Focus() {
	i.focused = true
	if i.onFocus != nil {
		i.onFocus()
	}
}

// Blur removes focus from the input
func (i *Input) Blur() {
	i.focused = false
	i.clearSelection()
	if i.onBlur != nil {
		i.onBlur()
	}
}

// IsFocused returns whether the input has focus
func (i *Input) IsFocused() bool {
	return i.focused
}

// clearSelection removes text selection
func (i *Input) clearSelection() {
	i.selectionStart = -1
	i.selectionEnd = -1
}

// hasSelection returns whether text is selected
func (i *Input) hasSelection() bool {
	return i.selectionStart >= 0 && i.selectionEnd >= 0 && i.selectionStart != i.selectionEnd
}

// getSelectedText returns the currently selected text
func (i *Input) getSelectedText() string {
	if !i.hasSelection() {
		return ""
	}

	start := i.selectionStart
	end := i.selectionEnd
	if start > end {
		start, end = end, start
	}

	runes := []rune(i.text)
	if start < 0 || start >= len(runes) || end < 0 || end > len(runes) {
		return ""
	}

	return string(runes[start:end])
}

// deleteSelection removes selected text
func (i *Input) deleteSelection() {
	if !i.hasSelection() {
		return
	}

	start := i.selectionStart
	end := i.selectionEnd
	if start > end {
		start, end = end, start
	}

	runes := []rune(i.text)
	newRunes := append(runes[:start], runes[end:]...)
	i.text = string(newRunes)
	i.cursorPos = start
	i.clearSelection()
}

// insertText inserts text at the cursor position
func (i *Input) insertText(text string) {
	// Normalize input text
	normalizedText := norm.NFC.String(text)

	// Remove any existing selection
	if i.hasSelection() {
		i.deleteSelection()
	}

	// Insert new text
	runes := []rune(i.text)
	newRunes := []rune(normalizedText)

	// Check length limit
	totalLength := len(runes) + len(newRunes)
	if i.maxLength > 0 && totalLength > i.maxLength {
		maxNewRunes := i.maxLength - len(runes)
		if maxNewRunes <= 0 {
			return
		}
		newRunes = newRunes[:maxNewRunes]
	}

	// Build new text
	result := append(runes[:i.cursorPos], append(newRunes, runes[i.cursorPos:]...)...)
	newText := string(result)

	// Apply validation
	if i.validator != nil && !i.validator(newText) {
		return
	}

	i.text = newText
	i.cursorPos += len(newRunes)

	if i.onChange != nil {
		i.onChange(i.text)
	}
}

// handleClick processes mouse click events
func (i *Input) handleClick(event gui.Event) bool {
	clickEvent := event.(*gui.ClickEvent)

	if !i.ContainsPoint(clickEvent.X, clickEvent.Y) {
		if i.focused {
			i.Blur()
		}
		return false
	}

	if !i.focused {
		i.Focus()
	}

	// Calculate cursor position from click location
	x, y, _, _ := i.GetBounds()
	relativeX := clickEvent.X - x - 5 // Account for padding

	if relativeX <= 0 {
		i.cursorPos = 0
	} else {
		// Find closest character position
		runes := []rune(i.text)
		width := 0

		for pos, r := range runes {
			charWidth := getCharWidth(r, i.font)
			if width+charWidth/2 > relativeX {
				i.cursorPos = pos
				break
			}
			width += charWidth
			i.cursorPos = pos + 1
		}
	}

	i.clearSelection()
	return true
}

// handleKeyPress processes keyboard input
func (i *Input) handleKeyPress(event gui.Event) bool {
	if !i.focused {
		return false
	}

	keyEvent := event.(*gui.KeyPressEvent)

	switch keyEvent.Key {
	case gui.KeyBackspace:
		if i.hasSelection() {
			i.deleteSelection()
		} else if i.cursorPos > 0 {
			runes := []rune(i.text)
			i.text = string(append(runes[:i.cursorPos-1], runes[i.cursorPos:]...))
			i.cursorPos--
		}

		if i.onChange != nil {
			i.onChange(i.text)
		}
		return true

	case gui.KeyDelete:
		if i.hasSelection() {
			i.deleteSelection()
		} else {
			runes := []rune(i.text)
			if i.cursorPos < len(runes) {
				i.text = string(append(runes[:i.cursorPos], runes[i.cursorPos+1:]...))
			}
		}

		if i.onChange != nil {
			i.onChange(i.text)
		}
		return true

	case gui.KeyArrowLeft:
		if i.cursorPos > 0 {
			i.cursorPos--
		}
		i.clearSelection()
		return true

	case gui.KeyArrowRight:
		if i.cursorPos < utf8.RuneCountInString(i.text) {
			i.cursorPos++
		}
		i.clearSelection()
		return true

	case gui.KeyEnter:
		if i.onSubmit != nil {
			i.onSubmit(i.text)
		}
		return true
	}

	return false
}

// handleTextInput processes text input events
func (i *Input) handleTextInput(event gui.Event) bool {
	if !i.focused {
		return false
	}

	textEvent := event.(*gui.TextInputEvent)

	// Filter out control characters
	filtered := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\t' {
			return -1
		}
		return r
	}, textEvent.Text)

	if filtered != "" {
		i.insertText(filtered)
	}

	return true
}

// handleFocus processes focus events
func (i *Input) handleFocus(event gui.Event) bool {
	i.Focus()
	return true
}

// handleBlur processes blur events
func (i *Input) handleBlur(event gui.Event) bool {
	i.Blur()
	return true
}

// getCharWidth calculates the width of a character
func getCharWidth(r rune, font font.Face) int {
	if font == nil {
		return 7 // Default width
	}

	advance, ok := font.GlyphAdvance(r)
	if !ok {
		return 7 // Default width
	}

	return advance.Ceil()
}

// Render draws the input component
func (i *Input) Render(canvas gui.Canvas) error {
	if !i.IsVisible() {
		return nil
	}

	x, y, width, height := i.GetBounds()

	// Draw background
	if err := canvas.DrawRectangle(x, y, width, height, i.bgColor, true); err != nil {
		return err
	}

	// Draw border
	if err := canvas.DrawRectangle(x, y, width, height, i.borderColor, false); err != nil {
		return err
	}

	// Draw text or placeholder
	textX := x + 5            // Padding
	textY := y + height/2 + 4 // Center vertically

	if i.text != "" {
		// Render actual text
		if err := canvas.DrawText(i.text, textX, textY, i.font, i.textColor); err != nil {
			return err
		}

		// Draw cursor if focused
		if i.focused {
			cursorX := textX + i.getCursorPixelPosition()
			cursorColor := colorful.Color{R: 0, G: 0, B: 0} // Black cursor
			if err := canvas.DrawRectangle(cursorX, y+2, 1, height-4, cursorColor, true); err != nil {
				return err
			}
		}
	} else if i.placeholder != "" {
		// Render placeholder text
		if err := canvas.DrawText(i.placeholder, textX, textY, i.font, i.placeholderColor); err != nil {
			return err
		}
	}

	return nil
}

// getCursorPixelPosition calculates the pixel position of the cursor
func (i *Input) getCursorPixelPosition() int {
	if i.cursorPos <= 0 {
		return 0
	}

	runes := []rune(i.text)
	if i.cursorPos > len(runes) {
		i.cursorPos = len(runes)
	}

	width := 0
	for j := 0; j < i.cursorPos && j < len(runes); j++ {
		width += getCharWidth(runes[j], i.font)
	}

	return width
}
