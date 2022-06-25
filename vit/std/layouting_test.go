package std

import (
	"fmt"
	"math"
	"testing"

	"github.com/omniskop/vitrum/vit"
)

type prop struct {
	name       string
	anchors    string
	value      interface{}
	expression string
}

var layoutingTests = []struct {
	id       string // some name to identify the test
	set      []prop
	expected vit.Rect
}{{
	id:       "0",
	set:      []prop{{name: "width", value: 100}, {name: "height", value: 80}},
	expected: vit.NewRect(0, 0, 100, 80),
}, {
	id:       "10",
	set:      []prop{{name: "width", value: 100}, {name: "height", value: 80}, {name: "x", value: 20}, {name: "y", value: 30}},
	expected: vit.NewRect(20, 30, 100, 80),
}, {
	id:       "20",
	set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "top", expression: "parent.top"}},
	expected: vit.NewRect(0, 10, 100, 100),
},
	{
		id:       "30",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "top", expression: "parent.bottom"}},
		expected: vit.NewRect(0, 1010, 100, 100),
	}, {
		id:       "40",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "top", expression: "parent.top"}, {anchors: "topMargin", value: 20}},
		expected: vit.NewRect(0, 30, 100, 100),
	},
	{
		id:       "50",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "left", expression: "parent.left"}},
		expected: vit.NewRect(10, 0, 100, 100),
	}, {
		id:       "60",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "left", expression: "parent.right"}},
		expected: vit.NewRect(1010, 0, 100, 100),
	}, {
		id:       "70",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "left", expression: "parent.left"}, {anchors: "leftMargin", value: 20}},
		expected: vit.NewRect(30, 0, 100, 100),
	},
	{
		id:       "80",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "right", expression: "parent.right"}},
		expected: vit.NewRect(1010-100, 0, 100, 100),
	}, {
		id:       "90",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "right", expression: "parent.left"}},
		expected: vit.NewRect(-90, 0, 100, 100),
	}, {
		id:       "100",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "right", expression: "parent.right"}, {anchors: "rightMargin", value: 20}},
		expected: vit.NewRect(1010-100-20, 0, 100, 100),
	},
	{
		id:       "110",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "bottom", expression: "parent.bottom"}},
		expected: vit.NewRect(0, 1010-100, 100, 100),
	}, {
		id:       "120",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "bottom", expression: "parent.top"}},
		expected: vit.NewRect(0, 10-100, 100, 100),
	}, {
		id:       "130",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "bottom", expression: "parent.bottom"}, {anchors: "bottomMargin", value: 20}},
		expected: vit.NewRect(0, 1010-100-20, 100, 100),
	},
	{
		id:       "140",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "horizontalCenter", expression: "parent.horizontalCenter"}},
		expected: vit.NewRect((1000-100)/2+10, 0, 100, 100),
	}, {
		id:       "150",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "horizontalCenter", expression: "parent.horizontalCenter"}, {anchors: "horizontalCenterOffset", value: -20}},
		expected: vit.NewRect((1000-100)/2+10-20, 0, 100, 100),
	}, {
		id:       "160",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "verticalCenter", expression: "parent.verticalCenter"}},
		expected: vit.NewRect(0, (1000-100)/2+10, 100, 100),
	}, {
		id:       "170",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "verticalCenter", expression: "parent.verticalCenter"}, {anchors: "verticalCenterOffset", value: -20}},
		expected: vit.NewRect(0, (1000-100)/2+10-20, 100, 100),
	},
	{
		id:       "180",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "centerIn", expression: "parent"}},
		expected: vit.NewRect((1000-100)/2+10, (1000-100)/2+10, 100, 100),
	}, {
		id:       "190",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "centerIn", expression: "parent"}, {anchors: "horizontalCenterOffset", value: 20}, {anchors: "verticalCenterOffset", value: 30}},
		expected: vit.NewRect((1000-100)/2+10+20, (1000-100)/2+10+30, 100, 100),
	},
	{
		id:       "200",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "fill", expression: "parent"}},
		expected: vit.NewRect(10, 10, 1000, 1000),
	}, {
		id:       "210",
		set:      []prop{{name: "width", value: 100}, {name: "height", value: 100}, {anchors: "fill", expression: "parent"}, {anchors: "leftMargin", value: 20}, {anchors: "topMargin", value: 30}, {anchors: "rightMargin", value: 40}, {anchors: "bottomMargin", value: 50}},
		expected: vit.Rect{30, 40, 970, 960},
	},
	{
		id:       "220",
		set:      []prop{{name: "height", value: 100}, {anchors: "left", expression: "parent.left"}, {anchors: "right", expression: "parent.right"}, {anchors: "top", expression: "parent.top"}},
		expected: vit.NewRect(10, 10, 1000, 100),
	},
}

func TestLayouting(t *testing.T) {
testLoop:
	for i, test := range layoutingTests {
		container := NewItem("", vit.NewComponentContainer())
		itm := NewItem("", vit.NewComponentContainer())
		container.AddChild(itm)

		container.SetProperty("x", 10)
		container.SetProperty("y", 10)
		container.SetProperty("width", 1000)
		container.SetProperty("height", 1000)

		for _, prop := range test.set {
			var err error
			if prop.name != "" {
				if prop.value != nil {
					err = itm.SetProperty(prop.name, prop.value)
				} else {
					err = itm.SetPropertyExpression(prop.name, prop.expression, nil)
				}
			} else if prop.anchors != "" {
				if prop.value != nil {
					err = itm.anchors.SetProperty(prop.anchors, prop.value)
				} else {
					ok := itm.anchors.SetPropertyExpression(prop.anchors, prop.expression, nil)
					if !ok {
						err = fmt.Errorf("anchors.%s not found", prop.anchors)
					}
				}
			}
			if err != nil {
				t.Logf("=> test %q (index: %d)", test.id, i)
				t.Error(err)
				continue testLoop
			}
		}

		err := updateAllExpressions(container)
		if err != nil {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expression update failed:")
			if grp, ok := err.(vit.ErrorGroup); ok {
				for _, err := range grp.Errors {
					t.Errorf("\t%v", err)
				}
			}
			continue testLoop
		}

		bounds := itm.Bounds()
		if !valuesEqual(test.expected.X1, bounds.X1) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected left edge to be %v, got %v", test.expected.X1, bounds.X1)
		}
		if !valuesEqual(test.expected.Y1, bounds.Y1) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected top edge to be %v, got %v", test.expected.Y1, bounds.Y1)
		}
		if !valuesEqual(test.expected.X2, bounds.X2) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected right edge to be %v, got %v", test.expected.X2, bounds.X2)
		}
		if !valuesEqual(test.expected.Y2, bounds.Y2) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected bottom edge to be %v, got %v", test.expected.Y2, bounds.Y2)
		}
		if !valuesEqual(test.expected.CenterX(), bounds.CenterX()) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected center X to be %v, got %v", test.expected.CenterX(), bounds.CenterX())
		}
		if !valuesEqual(test.expected.CenterY(), bounds.CenterY()) {
			t.Logf("=> test %q (index: %d)", test.id, i)
			t.Errorf("expected center Y to be %v, got %v", test.expected.CenterY(), bounds.CenterY())
		}
	}
}

func valuesEqual(a, b interface{}) bool {

	fB, ok := b.(float64)
	if ok {
		fA, ok := a.(float64)
		if ok {
			return math.Abs(fA-fB) < 0.000001
		}
		iA, ok := a.(int)
		if ok {
			return math.Abs(float64(iA)-fB) < 0.000001
		}
	}

	return a == b
}

func updateAllExpressions(comp vit.Component) error {
evaluateExpressions:
	n, errs := comp.UpdateExpressions()
	if errs.Failed() {
		return errs
	}
	if n > 0 {
		goto evaluateExpressions
	}
	return nil
}
