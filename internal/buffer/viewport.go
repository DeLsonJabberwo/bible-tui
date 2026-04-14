package buffer

const MAX_WIDTH int = 60
const PADDING int = 4

type ViewportInfo struct {
	Width	int
}

func NewViewportInfo(width int) ViewportInfo {
	if width == 0 {
		return ViewportInfo{MAX_WIDTH}
	} else {
		return ViewportInfo{width}
	}
}

func (v *ViewportInfo) WordWidthLimit() int {
	if v.Width < MAX_WIDTH + (PADDING * 2) {
		return (v.Width - (PADDING * 2))
	} else {
		return MAX_WIDTH
	}
}

func (v *ViewportInfo) MaxWidth() int {
	return (v.Width - (PADDING * 2))
}
