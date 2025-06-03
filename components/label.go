package components

import (
	"strings"

	"github.com/opd-ai/gui"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/text/unicode/norm"
)

// TextAlignment defines text alignment options
type TextAlignment int

const (
	AlignLeft TextAlignment = iota
	AlignCenter
	AlignRight
)

// Label represents a text display component
type Label struct {
	*gui.Element
	text      string
	font      font.Face
	color     colorful.Color
	alignment TextAlignment
	wordWrap  bool
	autoSize  bool
}

// NewLabel creates a new label with specified text
func NewLabel(text string) *Label {
	// Normalize text for consistent UTF-8 handling
	normalizedText := norm.NFC.String(text)

	label := &Label{
		Element:   gui.NewElement(0, 0, 100, 20),
		text:      normalizedText,
		font:      basicfont.Face7x13,               // Default font
		color:     colorful.Color{R: 0, G: 0, B: 0}, // Black
		alignment: AlignLeft,
		wordWrap:  false,
		autoSize:  true,
	}

	if label.autoSize {
		label.updateSize()
	}

	return label
}

// SetText updates the label's text content
func (l *Label) SetText(text string) *Label {
	normalizedText := norm.NFC.String(text)
	l.text = normalizedText

	if l.autoSize {
		l.updateSize()
	}

	return l
}

// GetText returns the current text
func (l *Label) GetText() string {
	return l.text
}

// SetFont updates the label's font
func (l *Label) SetFont(font font.Face) *Label {
	l.font = font

	if l.autoSize {
		l.updateSize()
	}

	return l
}

// SetColor updates the text color
func (l *Label) SetColor(color colorful.Color) *Label {
	l.color = color
	return l
}

// SetAlignment updates text alignment
func (l *Label) SetAlignment(alignment TextAlignment) *Label {
	l.alignment = alignment
	return l
}

// SetWordWrap enables or disables word wrapping
func (l *Label) SetWordWrap(enable bool) *Label {
	l.wordWrap = enable
	return l
}

// SetAutoSize enables or disables automatic sizing
func (l *Label) SetAutoSize(enable bool) *Label {
	l.autoSize = enable

	if enable {
		l.updateSize()
	}

	return l
}

// updateSize calculates and sets the optimal size for the text
func (l *Label) updateSize() {
	if l.font == nil || l.text == "" {
		l.SetSize(0, 0)
		return
	}

	metrics := l.font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()

	if !l.wordWrap {
		// Single line - measure width directly
		width := measureTextWidth(l.text, l.font)
		l.SetSize(width, lineHeight)
		return
	}

	// Multi-line with word wrap
	_, _, width, _ := l.GetBounds()
	if width <= 0 {
		width = 200 // Default width for word wrapping
	}

	lines := l.wrapText(l.text, width)
	totalHeight := len(lines) * lineHeight
	maxWidth := 0

	for _, line := range lines {
		lineWidth := measureTextWidth(line, l.font)
		if lineWidth > maxWidth {
			maxWidth = lineWidth
		}
	}

	l.SetSize(maxWidth, totalHeight)
}

// wrapText breaks text into lines that fit within the specified width
func (l *Label) wrapText(text string, maxWidth int) []string {
	if maxWidth <= 0 {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		if measureTextWidth(testLine, l.font) <= maxWidth {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// measureTextWidth calculates the pixel width of text
func measureTextWidth(text string, font font.Face) int {
	if font == nil || text == "" {
		return 0
	}

	width := 0
	for _, r := range text {
		advance, ok := font.GlyphAdvance(r)
		if ok {
			width += advance.Ceil()
		}
	}

	return width
}

// Render draws the label to the canvas
func (l *Label) Render(canvas gui.Canvas) error {
	if !l.IsVisible() || l.text == "" {
		return nil
	}

	x, y, width, _ := l.GetBounds()

	if l.wordWrap && width > 0 {
		return l.renderMultiLine(canvas, x, y, width)
	}

	return l.renderSingleLine(canvas, x, y, width)
}

// renderSingleLine renders text as a single line
func (l *Label) renderSingleLine(canvas gui.Canvas, x, y, width int) error {
	textX := x

	// Apply horizontal alignment
	if l.alignment != AlignLeft && width > 0 {
		textWidth := measureTextWidth(l.text, l.font)
		switch l.alignment {
		case AlignCenter:
			textX = x + (width-textWidth)/2
		case AlignRight:
			textX = x + width - textWidth
		}
	}

	// Position text at baseline
	metrics := l.font.Metrics()
	textY := y + metrics.Ascent.Ceil()

	return canvas.DrawText(l.text, textX, textY, l.font, l.color)
}

// renderMultiLine renders wrapped text across multiple lines
func (l *Label) renderMultiLine(canvas gui.Canvas, x, y, width int) error {
	lines := l.wrapText(l.text, width)
	metrics := l.font.Metrics()
	lineHeight := (metrics.Ascent + metrics.Descent).Ceil()

	for i, line := range lines {
		lineY := y + i*lineHeight + metrics.Ascent.Ceil()
		lineX := x

		// Apply horizontal alignment for each line
		if l.alignment != AlignLeft {
			lineWidth := measureTextWidth(line, l.font)
			switch l.alignment {
			case AlignCenter:
				lineX = x + (width-lineWidth)/2
			case AlignRight:
				lineX = x + width - lineWidth
			}
		}

		if err := canvas.DrawText(line, lineX, lineY, l.font, l.color); err != nil {
			return err
		}
	}

	return nil
}
