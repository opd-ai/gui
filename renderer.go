package gui

// Renderer provides platform-specific window and rendering operations
type Renderer interface {
	// Window management
	Show(title string) error
	Close() error

	// Canvas creation
	CreateCanvas() (Canvas, error)

	// Event handling
	PollEvents() []Event

	// Properties
	Size() (width, height int)
	SetSize(width, height int) error
}

// NewRenderer creates a platform-specific renderer
func NewRenderer(width, height int) (Renderer, error) {
	return newPlatformRenderer(width, height)
}
