package vit

/*
	Layout hierarchy: (from most important to least)
	- Fill through 'fill'
	- Centering through 'centerIn'
	- Positioning through 'left', 'right', 'top', 'bottom', 'horizontalCenter', 'verticalCenter'
	- Sizing through 'width', 'height'

	Margins:
	- are applicable to 'fill', but not 'centerIn'
	- A specific margin (through 'leftMargin', 'rightMargin', 'topMargin', 'bottomMargin') overrides the general 'margins'
*/

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

type AnchorsValue struct {
	Expression
	AlignWhenCentered      BoolValue
	Baseline               FloatValue
	BaselineOffset         FloatValue
	Bottom                 FloatValue
	BottomMargin           FloatValue
	CenterIn               ComponentValue
	Fill                   ComponentValue
	HorizontalCenter       FloatValue
	HorizontalCenterOffset FloatValue
	Left                   FloatValue
	LeftMargin             FloatValue
	Margins                FloatValue
	Right                  FloatValue
	RightMargin            FloatValue
	Top                    FloatValue
	TopMargin              FloatValue
	VerticalCenter         FloatValue
	VerticalCenterOffset   FloatValue
}

func NewAnchors(expression string, position *PositionRange) AnchorsValue {
	return AnchorsValue{
		AlignWhenCentered: *NewBoolValue("true", nil),
	}
}

func (a *AnchorsValue) Property(key string) (interface{}, bool) {
	switch key {
	case "alignWhenCentered":
		return a.AlignWhenCentered, true
	case "baseline":
		return a.Baseline, true
	case "baselineOffset":
		return a.BaselineOffset, true
	case "bottom":
		return a.Bottom, true
	case "bottomMargin":
		return a.BottomMargin, true
	case "centerIn":
		return a.CenterIn, true
	case "fill":
		return a.Fill, true
	case "horizonzalCenter":
		return a.HorizontalCenter, true
	case "horizonzalCenterOffset":
		return a.HorizontalCenterOffset, true
	case "left":
		return a.Left, true
	case "leftMargin":
		return a.LeftMargin, true
	case "margins":
		return a.Margins, true
	case "right":
		return a.Right, true
	case "rightMargin":
		return a.RightMargin, true
	case "top":
		return a.Top, true
	case "topMargin":
		return a.TopMargin, true
	case "verticalCenter":
		return a.VerticalCenter, true
	case "verticalCenterOffset":
		return a.VerticalCenterOffset, true
	}
	return nil, false
}

func (a *AnchorsValue) SetProperty(key string, expression string, position *PositionRange) bool {
	switch key {
	case "alignWhenCentered":
		a.AlignWhenCentered.ChangeCode(expression, position)
	case "baseline":
		a.Baseline.ChangeCode(expression, position)
	case "baselineOffset":
		a.BaselineOffset.ChangeCode(expression, position)
	case "bottom":
		a.Bottom.ChangeCode(expression, position)
	case "bottomMargin":
		a.BottomMargin.ChangeCode(expression, position)
	// case "centerIn":
	// 	a.CenterIn.ChangeCode(expression, position)
	// case "fill":
	// 	a.Fill.ChangeCode(expression, position)
	case "horizonzalCenter":
		a.HorizontalCenter.ChangeCode(expression, position)
	case "horizonzalCenterOffset":
		a.HorizontalCenterOffset.ChangeCode(expression, position)
	case "left":
		a.Left.ChangeCode(expression, position)
	case "leftMargin":
		a.LeftMargin.ChangeCode(expression, position)
	case "margins":
		a.Margins.ChangeCode(expression, position)
	case "right":
		a.Right.ChangeCode(expression, position)
	case "rightMargin":
		a.RightMargin.ChangeCode(expression, position)
	case "top":
		a.Top.ChangeCode(expression, position)
	case "topMargin":
		a.TopMargin.ChangeCode(expression, position)
	case "verticalCenter":
		a.VerticalCenter.ChangeCode(expression, position)
	case "verticalCenterOffset":
		a.VerticalCenterOffset.ChangeCode(expression, position)
	default:
		return false
	}
	return true
}

func (v *AnchorsValue) SetFromProperty(def PropertyDefinition) {
	panic("not implemented")
}

func (v *AnchorsValue) Update(context Component) error {
	panic("not implemented")
}

func (v *AnchorsValue) GetValue() interface{} {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) MakeDirty([]*Expression) {
	panic("not implemented")
}

func (v *AnchorsValue) GetExpression() *Expression {
	panic("not implemented")
	return nil
}

func (v *AnchorsValue) AddDependent(dep *Expression) {
	panic("not implemented")
}

func (v *AnchorsValue) RemoveDependent(dep *Expression) {
	panic("not implemented")
}

func (v *AnchorsValue) ShouldEvaluate() bool {
	panic("not implemented")
	return false
}

func (v *AnchorsValue) Err() error {
	panic("not implemented")
	return nil
}

func parseAnchorLine() {

}
