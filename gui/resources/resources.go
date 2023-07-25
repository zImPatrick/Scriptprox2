package resources

import (
	g "github.com/AllenDang/giu"
)

var IconFont *g.FontInfo

func WithIconFont(widgets ...g.Widget) *g.StyleSetter {
	return g.Style().SetFont(IconFont).To(widgets...)
}
