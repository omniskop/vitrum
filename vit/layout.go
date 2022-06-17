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
	l.x = x
	l.y = y
	l.positionChanged = true
}

func (l *Layout) SetPreferredPosition(x, y *float64) {
	l.preferredX = x
	l.preferredY = y
	l.positionChanged = true
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
	l.targetWidth = width
	l.targetHeight = height
	l.targetSizeChanged = true
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
	l.width = width
	l.height = height
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
