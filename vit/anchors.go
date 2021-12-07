package vit

type AnchorLine byte

const (
	AnchorsNone AnchorLine = iota
	AnchorLeft
	AnchorRight
	AnchorTop
	AnchorBottom
	AnchorHorizontalCenter
	AnchorVerticalCenter
)

type Anchors struct {
	AlignWhenCentered      bool
	Baseline               AnchorLine
	BaselineOffset         float64
	Bottom                 AnchorLine
	BottomMargin           float64
	CenterIn               Component
	Fill                   Component
	HorizonzalCenter       AnchorLine
	HorizonzalCenterOffset float64
	Left                   AnchorLine
	LeftMargin             float64
	Margins                float64
	Right                  AnchorLine
	RightMargin            float64
	Top                    AnchorLine
	TopMargin              float64
	VerticalCenter         AnchorLine
	VerticalCenterOffset   float64
}

func NewAnchors() Anchors {
	return Anchors{
		AlignWhenCentered: true,
	}
}

func (a *Anchors) Property(key string) interface{} {
	switch key {
	case "alignWhenCentered":
		return a.AlignWhenCentered
	case "baseline":
		return a.Baseline
	case "baselineOffset":
		return a.BaselineOffset
	case "bottom":
		return a.Bottom
	case "bottomMargin":
		return a.BottomMargin
	case "centerIn":
		return a.CenterIn
	case "fill":
		return a.Fill
	case "horizonzalCenter":
		return a.HorizonzalCenter
	case "horizonzalCenterOffset":
		return a.HorizonzalCenterOffset
	case "left":
		return a.Left
	case "leftMargin":
		return a.LeftMargin
	case "margins":
		return a.Margins
	case "right":
		return a.Right
	case "rightMargin":
		return a.RightMargin
	case "top":
		return a.Top
	case "topMargin":
		return a.TopMargin
	case "verticalCenter":
		return a.VerticalCenter
	case "verticalCenterOffset":
		return a.VerticalCenterOffset
	}
	return nil
}

func (a *Anchors) SetProperty(key string, value interface{}) bool {
	return true
	switch key {
	case "alignWhenCentered":
		a.AlignWhenCentered = value.(bool)
	case "baseline":
		a.Baseline = value.(AnchorLine)
	case "baselineOffset":
		a.BaselineOffset = value.(float64)
	case "bottom":
		a.Bottom = value.(AnchorLine)
	case "bottomMargin":
		a.BottomMargin = value.(float64)
	case "centerIn":
		a.CenterIn = value.(Component)
	case "fill":
		a.Fill = value.(Component)
	case "horizonzalCenter":
		a.HorizonzalCenter = value.(AnchorLine)
	case "horizonzalCenterOffset":
		a.HorizonzalCenterOffset = value.(float64)
	case "left":
		a.Left = value.(AnchorLine)
	case "leftMargin":
		a.LeftMargin = value.(float64)
	case "margins":
		a.Margins = value.(float64)
	case "right":
		a.Right = value.(AnchorLine)
	case "rightMargin":
		a.RightMargin = value.(float64)
	case "top":
		a.Top = value.(AnchorLine)
	case "topMargin":
		a.TopMargin = value.(float64)
	case "verticalCenter":
		a.VerticalCenter = value.(AnchorLine)
	case "verticalCenterOffset":
		a.VerticalCenterOffset = value.(float64)
	default:
		return false
	}
	return true
}
