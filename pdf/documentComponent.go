package pdf

import (
	vit "github.com/omniskop/vitrum/vit"
)

func (d *DocumentComponent) recalculateLayout() {
	var x, y float64
	width, height := pageSize(PageComponent_Format(d.format.Int()), PageComponent_Orientation(d.orientation.Int()))
	for _, child := range d.Children() {
		var localW, localH float64 = width, height
		layout := d.childLayouts[child]
		if w, ok := layout.GetTargetWidth(); ok && w != 0 {
			localW = w
		}
		if h, ok := layout.GetTargetHeight(); ok && h != 0 {
			localH = h
		}

		d.childLayouts[child].SetPosition(&x, &y)
		d.childLayouts[child].SetSize(&localW, &localH)
	}
}

func (d *DocumentComponent) createNewChildLayout(child vit.Component) *vit.Layout {
	l := vit.NewLayout()
	d.childLayouts[child] = l
	l.SetTargetSize(nil, nil)
	l.AddDependent(vit.FuncDep(d.recalculateLayout))
	return l
}

func (d *DocumentComponent) childWasAdded(child vit.Component) {
	child.ApplyLayout(d.createNewChildLayout(child))
	d.recalculateLayout()
}
