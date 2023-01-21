package std

import (
	"fmt"
	"unicode/utf8"

	vit "github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/vfont"
	"github.com/tdewolff/canvas"
)

func (t *Text) updateFont() {
	familyName := t.font.MustGet("family").GetValue().(string)
	if familyName == "" {
		t.fontData = nil
		t.fontFaceData = nil
		return
	}
	if t.fontData != nil && t.fontData.Name() != familyName {
		t.fontData = canvas.NewFontFamily(familyName)
	}

	var weight Text_FontWeight = Text_FontWeight(t.font.MustGet("weight").GetValue().(int))
	if t.font.MustGet("bold").GetValue().(bool) {
		weight = Text_FontWeight_Bold
	}

	var err error
	t.fontFaceData, err = vfont.LoadFontFace(familyName, vfont.Style{
		Color:     t.color.Color(),
		PointSize: t.font.MustGet("pointSize").GetValue().(float64),
		Italic:    t.font.MustGet("italic").GetValue().(bool),
		Underline: t.font.MustGet("underline").GetValue().(bool),
		Strikeout: t.font.MustGet("strikeout").GetValue().(bool),
		Weight:    vfont.Weight(weight),
	})
	if err != nil {
		fmt.Println("Text: font:", err)
	}

	t.SetContentSize(t.fontFaceData.TextWidth(t.text.String()), t.fontFaceData.LineHeight())
}

func (t *Text) Draw(ctx vit.DrawingContext, area vit.Rect) error {
	if t.fontFaceData == nil {
		return t.Root.DrawChildren(ctx, area)
	}

	compBounds := t.Bounds()
	var drawBounds vit.Rect

	// calculate alignment
	var hAlign canvas.TextAlign
	switch Text_HorizontalAlignment(t.horizontalAlignment.Int()) {
	case Text_HorizontalAlignment_AlignLeft:
		hAlign = canvas.Left
		drawBounds.X1 = compBounds.Left()
	case Text_HorizontalAlignment_AlignHCenter:
		hAlign = canvas.Center
		drawBounds.X1 = compBounds.CenterX()
	case Text_HorizontalAlignment_AlignRight:
		hAlign = canvas.Right
		drawBounds.X1 = compBounds.Right()
	}
	var vAlign canvas.TextAlign
	switch Text_VerticalAlignment(t.verticalAlignment.Int()) {
	case Text_VerticalAlignment_AlignTop:
		vAlign = canvas.Top
		drawBounds.Y1 = compBounds.Top()
	case Text_VerticalAlignment_AlignVCenter:
		vAlign = canvas.Center
		drawBounds.Y1 = compBounds.CenterY()
	case Text_VerticalAlignment_AlignBottom:
		vAlign = canvas.Bottom
		drawBounds.Y1 = compBounds.Bottom()
	}

	drawBounds.X2 = drawBounds.X1
	drawBounds.Y2 = drawBounds.Y1

	var textString = t.text.String()

	// calculate size
	if Text_Elide(t.elide.Int()) != Text_Elide_ElideNone {
		textWidth := t.fontFaceData.TextWidth(textString)
		if textWidth > compBounds.Width() {
			switch Text_Elide(t.elide.Int()) {
			case Text_Elide_ElideLeft:
				textString = t.elideTextLeft(textString, compBounds.Width())
			case Text_Elide_ElideMiddle:
				textString = t.elideTextMiddle(textString, compBounds.Width())
			case Text_Elide_ElideRight:
				textString = t.elideTextRight(textString, compBounds.Width())
			}
		}
	}

	textElement := canvas.NewTextBox(
		t.fontFaceData,
		textString,
		drawBounds.Width(),
		drawBounds.Height(),
		hAlign,
		vAlign,
		0, 0,
	)

	ctx.DrawText(drawBounds.X1, drawBounds.Y1, textElement)

	return t.Root.DrawChildren(ctx, area)
}

func (t *Text) elideTextLeft(textString string, maxWidth float64) string {
	elideWidth := t.fontFaceData.TextWidth("...")
	maxWidth -= elideWidth
	for len(textString) > 0 {
		// remove first rune
		_, size := utf8.DecodeRuneInString(textString)
		textString = textString[size:]

		textWidth := t.fontFaceData.TextWidth(textString)
		if textWidth <= maxWidth {
			return "..." + textString
		}
	}
	return ""
}

func (t *Text) elideTextMiddle(textString string, maxWidth float64) string {
	elideWidth := t.fontFaceData.TextWidth("...")
	maxWidth -= elideWidth
	for len(textString) > 0 {
		// remove center rune
		partA, partB := t.cutStringInMiddle(textString, 1)
		textString = partA + partB
		textWidth := t.fontFaceData.TextWidth(partA + "..." + partB)
		if textWidth <= maxWidth {
			return partA + "..." + partB
		}
	}
	return ""
}

func (t *Text) elideTextRight(textString string, maxWidth float64) string {
	elideWidth := t.fontFaceData.TextWidth("...")
	maxWidth -= elideWidth
	for len(textString) > 0 {
		// remove last rune
		_, size := utf8.DecodeLastRuneInString(textString)
		textString = textString[:len(textString)-size]

		textWidth := t.fontFaceData.TextWidth(textString)
		if textWidth <= maxWidth {
			return textString + "..."
		}
	}
	return ""
}

// cutStringInMiddle removed the given amount of runes from the middle of the string and returns the left and right parts.
// It does this in a utf8-safe way.
func (t *Text) cutStringInMiddle(txt string, amount int) (string, string) {
	halfLen := utf8.RuneCountInString(txt) / 2
	var byteIdx int
	var partAEndByteIndex int
	for runeIdx, r := range txt {
		byteIdx += utf8.RuneLen(r)
		if runeIdx < halfLen {
			partAEndByteIndex = byteIdx
		} else if amount > 0 {
			amount--
		} else {
			return txt[:partAEndByteIndex], txt[byteIdx-1:]
		}
	}
	return txt[:partAEndByteIndex], txt[byteIdx:] // not decrementing byteIdx here
}
