package handhomeui

import (
	"embed"

	"codeberg.org/mikolajgasior/gocrud/internal/poc/layout"
)

type Handler struct {
	embedHTML embed.FS
	layout    *layout.Layout
}

type HandlerInput struct {
	EmbedHTML embed.FS
	Layout    *layout.Layout
}

func New(input HandlerInput) *Handler {
	return &Handler{
		embedHTML: input.EmbedHTML,
		layout:    input.Layout,
	}
}
