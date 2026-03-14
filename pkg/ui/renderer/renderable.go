package renderer

// Renderable is an interface that objects can implement to control their own rendering
type Renderable interface {
	RenderString() string
}
