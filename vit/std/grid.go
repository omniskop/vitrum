package std

import (
	"math"

	vit "github.com/omniskop/vitrum/vit"
)

func (g *Grid) getTopPadding() float64 {
	if g.topPadding.IsSet() {
		return g.topPadding.Value().Float64()
	}
	return g.padding.Float64()
}

func (g *Grid) getRightPadding() float64 {
	if g.rightPadding.IsSet() {
		return g.rightPadding.Value().Float64()
	}
	return g.padding.Float64()
}

func (g *Grid) getBottomPadding() float64 {
	if g.bottomPadding.IsSet() {
		return g.bottomPadding.Value().Float64()
	}
	return g.padding.Float64()
}

func (g *Grid) getLeftPadding() float64 {
	if g.leftPadding.IsSet() {
		return g.leftPadding.Value().Float64()
	}
	return g.padding.Float64()
}

// Recalculate Layout of all child components.
func (g *Grid) recalculateLayout() {
	// First of this could probably be done more efficiently but it works for now.
	// We are trying to place each child in the grid according to the flow and alignment.

	// First we will create a list of all children that are actually visible. All others will be ignored.
	var children = make([]vit.Component, 0, len(g.Children()))
	for _, child := range g.Children() {
		bounds := child.Bounds()
		if bounds.Width() == 0 || bounds.Height() == 0 {
			continue
		}
		children = append(children, child)
	}

	// Now we calculate the correct number of rows and columns based on the properties and the number of children.
	var rowCount uint
	var columnCount uint
	if g.rows.IsSet() {
		// the rows are set
		rowCount = uint(g.rows.Value().Int())
		if g.columns.IsSet() {
			// the columns are set
			columnCount = uint(g.columns.Value().Int())
		} else {
			// calculate columns based on number of children and rows
			if rowCount != 0 { // prevent division by zero
				columnCount = uint(math.Ceil(float64(len(children)) / float64(rowCount)))
			}
		}
	} else {
		// we need to calculate rows, after we figured out the columns
		if g.columns.IsSet() {
			// the columns are set
			columnCount = uint(g.columns.Value().Int())

		} else {
			// default to 4 columns
			columnCount = 4
		}
		// now calculate the rows based on the number of children and columns
		if columnCount != 0 { // prevent division by zero
			rowCount = uint(math.Ceil(float64(len(children)) / float64(columnCount)))
		}
	}

	// Now we calculate the size of each column and row based on the maximum size of their respective children.
	columnSizes := make([]float64, columnCount)
	rowSizes := make([]float64, rowCount)
	for row := uint(0); row < rowCount; row++ {
		for column := uint(0); column < columnCount; column++ {
			child, ok := g.getChildInGrid(children, row, column, rowCount, columnCount)
			if !ok {
				continue // in LeftToRight flow we could break here but in TopToBottom that would be wrong
			}
			bounds := child.Bounds()
			if bounds.Height() > rowSizes[row] {
				rowSizes[row] = bounds.Height()
			}
			if bounds.Width() > columnSizes[column] {
				columnSizes[column] = bounds.Width()
			}
		}
	}

	// Now we begin to actually position the components.

	// these will contain the current coordinates
	x := g.left.Float64() + g.getLeftPadding()
	y := g.top.Float64() + g.getTopPadding()

	// Set the position of all children that don't fit into the grid to the top left corner (just as a sane default).
	// This will only be necessary if rows and columns are both set by the user and the number of cells is too low to fit all children.
	for i := int(rowCount * columnCount); i < len(children); i++ {
		child := children[i]
		g.childLayouts[child].SetPosition(&x, &y)
		x += columnSizes[0] + g.spacing.Float64()
	}

	// iterate through all rows and columns and set the child's position
	for row := uint(0); row < rowCount; row++ {
		for column := uint(0); column < columnCount; column++ {
			child, ok := g.getChildInGrid(children, row, column, rowCount, columnCount)
			if !ok {
				continue // in LeftToRight flow we could break here but in TopToBottom that would be wrong
			}

			offsetX, offsetY := g.calculateOffsetInCell(columnSizes[column], rowSizes[row], child)
			offsetX += x
			offsetY += y

			g.childLayouts[child].SetPosition(&offsetX, &offsetY)
			x += columnSizes[column] + g.spacing.Float64() // advance the coordinates in the x axis
		}
		x = g.left.Float64() + g.getLeftPadding() // reset x to the left
		y += rowSizes[row] + g.spacing.Float64()  // advance the coordinates in the y axis
	}

	// Add all column and rows sizes (including spacing) to calculate the width and height of the grid component.
	width := g.getLeftPadding() + g.getRightPadding()
	height := g.getTopPadding() + g.getBottomPadding()
	for _, size := range columnSizes {
		width += size + g.spacing.Float64()
	}
	if columnCount > 0 {
		width -= g.spacing.Float64() // subtract last spacing again
	}
	for _, size := range rowSizes {
		height += size + g.spacing.Float64()
	}
	if rowCount > 0 {
		height -= g.spacing.Float64() // subtract last spacing again
	}

	// as we now know how large the grid is we use that for our own layout
	g.contentWidth = width
	g.contentHeight = height
	g.layouting()
}

// getChildInGrid is a helper method that returns the child that should be visible at a specific row and column.
// It takes the grid's flow into account.
func (g *Grid) getChildInGrid(children []vit.Component, row, column, rowCount, columnCount uint) (vit.Component, bool) {
	var index uint
	if Grid_Flow(g.flow.Int()) == Grid_Flow_LeftToRight {
		index = row*columnCount + column
	} else if Grid_Flow(g.flow.Int()) == Grid_Flow_TopToBottom {
		index = column*rowCount + row
	}
	if index >= uint(len(children)) {
		return nil, false
	}
	return children[index], true
}

func (g *Grid) calculateOffsetInCell(cellWidth, cellHeight float64, comp vit.Component) (float64, float64) {
	bounds := comp.Bounds()
	var offsetX, offsetY float64

	// horizontal alignment
	switch Grid_HorizontalItemAlignment(g.horizontalItemAlignment.Int()) {
	case Grid_HorizontalItemAlignment_AlignLeft:
	case Grid_HorizontalItemAlignment_AlignHCenter:
		offsetX = (cellWidth - bounds.Width()) / 2
	case Grid_HorizontalItemAlignment_AlignRight:
		offsetX = cellWidth - bounds.Width()
	}

	// vertical alignment
	switch Grid_VerticalItemAlignment(g.verticalItemAlignment.Int()) {
	case Grid_VerticalItemAlignment_AlignTop:
	case Grid_VerticalItemAlignment_AlignVCenter:
		offsetY = (cellHeight - bounds.Height()) / 2
	case Grid_VerticalItemAlignment_AlignBottom:
		offsetY = cellHeight - bounds.Height()
	}

	return offsetX, offsetY
}

func (g *Grid) createNewChildLayout(child vit.Component) *vit.Layout {
	l := vit.NewLayout()
	g.childLayouts[child] = l
	l.SetTargetSize(nil, nil)
	l.AddDependent(vit.FuncDep(g.recalculateLayout))
	return l
}

func (g *Grid) childWasAdded(child vit.Component) {
	child.ApplyLayout(g.createNewChildLayout(child))
	g.recalculateLayout()
}
