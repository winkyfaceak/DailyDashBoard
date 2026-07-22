package features

import (
	"DailyDashBoard/internal/tui/terminalimage"
)

// imageState implements the shared two-phase image flow: while the raw
// terminal graphics are being flushed the session shows a placeholder, then
// swaps to the final rendered text. Embed it to satisfy ImageSession.
type imageState struct {
	image   terminalimage.Image
	pending bool
}

func newImageState(image terminalimage.Image) imageState {
	return imageState{
		image:   image,
		pending: image.Raw != "",
	}
}

// RawImage returns the raw graphics sequence to flush, if any.
func (s *imageState) RawImage() string {
	return s.image.Raw
}

// ImageShown marks the raw graphics as flushed.
func (s *imageState) ImageShown() {
	s.pending = false
}

// imageText returns the placeholder while pending, or the final image text.
func (s *imageState) imageText(placeholder string) string {
	if s.pending {
		return SubtleText(placeholder)
	}

	return s.image.Text
}
