package graphics

import (
	"image"

	"github.com/fogleman/gg"
	"github.com/lucasb-eyer/go-colorful"
	"golang.org/x/image/font"
)

// GGCanvas implements Canvas using the fogleman/gg library
type GGCanvas struct {
	context *gg.Context
	width   int
	height  int
}

// NewGGCanvas creates a new canvas using gg
func NewGGCanvas(width, height int) *GGCanvas {
	ctx := gg.NewContext(width, height)
	return &GGCanvas{
		context: ctx,
		width:   width,
		height:  height,
	}
}

// DrawText renders text at the specified position
func (c *GGCanvas) DrawText(text string, x, y int, fontFace font.Face, textColor colorful.Color) error {
	// Convert colorful.Color to standard color
	r, g, b := textColor.RGB255()
	c.context.SetRGB255(int(r), int(g), int(b))

	// Set font
	c.context.SetFontFace(fontFace)

	// Draw text
	c.context.DrawString(text, float64(x), float64(y))

	return nil
}

// DrawRectangle draws a rectangle with the specified parameters
func (c *GGCanvas) DrawRectangle(x, y, width, height int, rectColor colorful.Color, filled bool) error {
	// Convert colorful.Color to standard color
	r, g, b := rectColor.RGB255()
	c.context.SetRGB255(int(r), int(g), int(b))

	if filled {
		c.context.DrawRectangle(float64(x), float64(y), float64(width), float64(height))
		c.context.Fill()
	} else {
		c.context.DrawRectangle(float64(x), float64(y), float64(width), float64(height))
		c.context.Stroke()
	}

	return nil
}

// DrawCircle draws a circle with the specified parameters
func (c *GGCanvas) DrawCircle(x, y, radius int, circleColor colorful.Color, filled bool) error {
	// Convert colorful.Color to standard color
	r, g, b := circleColor.RGB255()
	c.context.SetRGB255(int(r), int(g), int(b))

	if filled {
		c.context.DrawCircle(float64(x), float64(y), float64(radius))
		c.context.Fill()
	} else {
		c.context.DrawCircle(float64(x), float64(y), float64(radius))
		c.context.Stroke()
	}

	return nil
}

// DrawImage draws an image at the specified position and size
func (c *GGCanvas) DrawImage(img image.Image, x, y, width, height int) error {
	// Scale image if dimensions don't match
	if img.Bounds().Dx() != width || img.Bounds().Dy() != height {
		// Use gg's built-in scaling
		c.context.DrawImageAnchored(img, x+width/2, y+height/2, 0.5, 0.5)
	} else {
		c.context.DrawImage(img, x, y)
	}

	return nil
}

// SetClippingRegion sets a clipping rectangle
func (c *GGCanvas) SetClippingRegion(x, y, width, height int) {
	c.context.DrawRectangle(float64(x), float64(y), float64(width), float64(height))
	c.context.Clip()
}

// ClearClippingRegion removes the current clipping region
func (c *GGCanvas) ClearClippingRegion() {
	c.context.ResetClip()
}

// Clear fills the entire canvas with the specified color
func (c *GGCanvas) Clear(bgColor colorful.Color) error {
	r, g, b := bgColor.RGB255()
	c.context.SetRGB255(int(r), int(g), int(b))
	c.context.Clear()
	return nil
}

// Present displays the rendered content (placeholder for platform-specific implementation)
func (c *GGCanvas) Present() error {
	// This will be implemented in platform-specific renderers
	return nil
}

// GetImage returns the current canvas as an image
func (c *GGCanvas) GetImage() image.Image {
	return c.context.Image()
}

// GetContext returns the underlying gg context for advanced operations
func (c *GGCanvas) GetContext() *gg.Context {
	return c.context
}
