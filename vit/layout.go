package vit

type Layout struct {
	// set by the parent
	x               *float64
	y               *float64
	preferredX      *float64
	preferredY      *float64
	positionChanged bool
	// set by the child
	targetWidth       *float64
	targetHeight      *float64
	targetSizeChanged bool
	// set by the parent
	width       *float64
	height      *float64
	sizeChanged bool
}

func (l *Layout) SetPosition(x, y *float64) {
	if x == nil {
		l.x = nil
	} else if l.x == nil || *x != *l.x {
		var xCopy float64 = *x
		l.x = &xCopy
		l.positionChanged = true
	}
	if y == nil {
		l.y = nil
	} else if l.y == nil || *y != *l.y {
		var yCopy float64 = *y
		l.y = &yCopy
		l.positionChanged = true
	}
}

func (l *Layout) SetPreferredPosition(x, y *float64) {
	if x == nil {
		l.preferredX = nil
	} else if l.preferredX == nil || *x != *l.preferredX {
		var xCopy float64 = *x
		l.preferredX = &xCopy
		l.positionChanged = true
	}
	if y == nil {
		l.preferredY = nil
	} else if l.preferredY == nil || *y != *l.preferredY {
		var yCopy float64 = *y
		l.preferredY = &yCopy
		l.positionChanged = true
	}
}

func (l *Layout) PositionChanged() bool {
	if l == nil {
		return false
	}
	return l.positionChanged
}

func (l *Layout) AckPositionChange() {
	l.positionChanged = false
}

func (l *Layout) GetX() (float64, bool) {
	if l.x == nil {
		return 0, false
	}
	return *l.x, true
}

func (l *Layout) GetY() (float64, bool) {
	if l.y == nil {
		return 0, false
	}
	return *l.y, true
}

func (l *Layout) GetPreferredX() (float64, bool) {
	if l.preferredX == nil {
		return 0, false
	}
	return *l.preferredX, true
}

func (l *Layout) GetPreferredY() (float64, bool) {
	if l.preferredY == nil {
		return 0, false
	}
	return *l.preferredY, true
}

func (l *Layout) SetTargetSize(width, height *float64) {
	if l == nil {
		return
	}
	if width == nil {
		l.targetWidth = nil
	} else if l.targetWidth == nil || *width != *l.targetWidth {
		var widthCopy float64 = *width
		l.targetWidth = &widthCopy
		l.targetSizeChanged = true
	}
	if height == nil {
		l.targetHeight = nil
	} else if l.targetHeight == nil || *height != *l.targetHeight {
		var heightCopy float64 = *height
		l.targetHeight = &heightCopy
		l.targetSizeChanged = true
	}
}

func (l *Layout) TargetSizeChanged() bool {
	if l == nil {
		return false
	}
	return l.targetSizeChanged
}

func (l *Layout) AckTargetSizeChange() {
	l.targetSizeChanged = false
}

func (l *Layout) GetTargetWidth() (float64, bool) {
	if l.targetWidth == nil {
		return 0, false
	}
	return *l.targetWidth, true
}

func (l *Layout) GetTargetHeight() (float64, bool) {
	if l.targetHeight == nil {
		return 0, false
	}
	return *l.targetHeight, true
}

func (l *Layout) SetSize(width, height *float64) {
	if width == nil {
		l.width = nil
	} else if l.width == nil || *width != *l.width {
		var widthCopy float64 = *width
		l.width = &widthCopy
		l.sizeChanged = true
	}
	if height == nil {
		l.height = nil
	} else if l.height == nil || *height != *l.height {
		var heightCopy float64 = *height
		l.height = &heightCopy
		l.sizeChanged = true
	}
}

func (l *Layout) SizeChanged() bool {
	if l == nil {
		return false
	}
	return l.sizeChanged
}

func (l *Layout) AckSizeChange() {
	l.sizeChanged = false
}

func (l *Layout) GetWidth() (float64, bool) {
	if l.width == nil {
		return 0, false
	}
	return *l.width, true
}

func (l *Layout) GetHeight() (float64, bool) {
	if l.height == nil {
		return 0, false
	}
	return *l.height, true
}

type LayoutList map[Component]*Layout

// ShouldEvaluate returns true if one of the contained layouts had it's target size changed.
// It adds compatibility with other Value types.
func (l LayoutList) ShouldEvaluate() bool {
	for _, layout := range l {
		if layout.TargetSizeChanged() {
			return true
		}
	}
	return false
}

// Update acknowledges the change of the target size on all contained layouts.
// It adds compatibility with other Value types.
func (l LayoutList) Update(Component) error {
	for _, layout := range l {
		layout.AckTargetSizeChange()
	}
	return nil
}

// GetExpression adds compatibility with other Value types but does not serve a particular purpose here.
func (l LayoutList) GetExpression() *Expression {
	return NewExpression("", nil)
}
