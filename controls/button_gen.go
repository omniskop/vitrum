// Code generated by vitrum gencmd. DO NOT EDIT.

package controls

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
	std "github.com/omniskop/vitrum/vit/std"
	vpath "github.com/omniskop/vitrum/vit/vpath"
)

func newFileContextForButton(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	fileCtx := vit.NewFileContext(globalCtx)

	var lib parse.Library
	var err error
	lib, err = parse.ResolveLibrary([]string{"Vit"})
	if err != nil {
		// The file used to generate the "Button" component imported a library called "Vit".
		// If this error occurs that imported failed. Probably because the library is not known.
		return nil, fmt.Errorf("unable to create file context for generated \"Button\" component: %w", err)
	}
	parse.AddLibraryToContainer(lib, &fileCtx.KnownComponents)

	return fileCtx, nil
}

type Button struct {
	*std.Item
	id string

	text          vit.StringValue
	pressed       vit.AliasValue
	fontPointSize vit.IntValue

	onClicked vit.EventAttribute[std.MouseEvent]

	clicked vit.Method
}

// newButtonInGlobal creates an appropriate file context for the component and then returns a new Button instance.
// The returned error will only be set if a library import that is required by the component fails.
func newButtonInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Button, error) {
	fileCtx, err := newFileContextForButton(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewButton(id, fileCtx), nil
}
func NewButton(id string, context *vit.FileContext) *Button {
	b := &Button{
		Item:          std.NewItem("", context),
		id:            id,
		text:          *vit.NewEmptyStringValue(),
		pressed:       *vit.NewAliasValueFromCode(vit.Code{FileCtx: context, Code: "mouseArea.pressed", Position: nil}),
		fontPointSize: *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "40", Position: nil}),
		onClicked:     *vit.NewEventAttribute[std.MouseEvent](),
		clicked:       vit.NewMethod("clicked", vit.Code{FileCtx: context, Code: "function(e) {\n        onClicked.Fire(e)\n    }", Position: nil}),
	}
	// property assignments on embedded components
	b.Item.SetPropertyCode("width", vit.Code{FileCtx: context, Code: "150", Position: nil})
	b.Item.SetPropertyCode("height", vit.Code{FileCtx: context, Code: "25", Position: nil})
	// register listeners for when a property changes
	// register event listeners
	var event vit.Listenable
	var listener vit.Evaluater
	event, _ = b.Root.Event("onCompleted")
	listener = event.CreateListener(vit.Code{FileCtx: context, Code: "function() {\n        mouseArea.onClicked.AddEventListener(clicked)\n    }", Position: nil})
	b.AddListenerFunction(listener)
	// register enumerations
	// add child components
	var child vit.Component
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "MouseArea", ID: "mouseArea", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 22, StartColumn: 9, EndLine: 22, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 22, StartColumn: 23, EndLine: 22, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 23, StartColumn: 9, EndLine: 23, EndColumn: 45}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 23, StartColumn: 26, EndLine: 23, EndColumn: 45}, Identifier: []string{"acceptedButtons"}, Expression: "MouseArea.leftButton", Tags: map[string]string{}}}}, context)
	b.AddChild(child)
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "Rectangle", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 27, StartColumn: 9, EndLine: 27, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 27, StartColumn: 23, EndLine: 27, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 28, StartColumn: 9, EndLine: 28, EndColumn: 82}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 28, StartColumn: 16, EndLine: 28, EndColumn: 82}, Identifier: []string{"color"}, Expression: "mouseArea.pressed ? Vit.rgb(180, 180, 180) : Vit.rgb(200, 200, 200)", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 29, StartColumn: 9, EndLine: 29, EndColumn: 17}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 29, StartColumn: 17, EndLine: 29, EndColumn: 17}, Identifier: []string{"radius"}, Expression: "5", Tags: map[string]string{}}}}, context)
	b.AddChild(child)
	child, _ = parse.InstantiateComponent(&vit.ComponentDefinition{BaseName: "Text", ID: "text", Properties: []vit.PropertyDefinition{vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 34, StartColumn: 9, EndLine: 34, EndColumn: 28}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 34, StartColumn: 23, EndLine: 34, EndColumn: 28}, Identifier: []string{"anchors", "fill"}, Expression: "parent", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 35, StartColumn: 9, EndLine: 35, EndColumn: 25}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 35, StartColumn: 15, EndLine: 35, EndColumn: 25}, Identifier: []string{"text"}, Expression: "parent.text", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 36, StartColumn: 9, EndLine: 36, EndColumn: 44}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 36, StartColumn: 25, EndLine: 36, EndColumn: 44}, Identifier: []string{"font", "pointSize"}, Expression: "parent.fontPointSize", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 37, StartColumn: 9, EndLine: 37, EndColumn: 33}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 37, StartColumn: 22, EndLine: 37, EndColumn: 33}, Identifier: []string{"font", "family"}, Expression: "\"Montserrat\"", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 38, StartColumn: 9, EndLine: 38, EndColumn: 32}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 38, StartColumn: 22, EndLine: 38, EndColumn: 32}, Identifier: []string{"font", "weight"}, Expression: "Text.Medium", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 39, StartColumn: 9, EndLine: 39, EndColumn: 44}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 39, StartColumn: 28, EndLine: 39, EndColumn: 44}, Identifier: []string{"verticalAlignment"}, Expression: "Text.AlignVCenter", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 40, StartColumn: 9, EndLine: 40, EndColumn: 46}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 40, StartColumn: 30, EndLine: 40, EndColumn: 46}, Identifier: []string{"horizontalAlignment"}, Expression: "Text.AlignHCenter", Tags: map[string]string{}}, vit.PropertyDefinition{Pos: vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 41, StartColumn: 9, EndLine: 41, EndColumn: 31}, ValuePos: &vit.PositionRange{FilePath: vpath.Virtual("Button.vit"), StartLine: 41, StartColumn: 16, EndLine: 41, EndColumn: 31}, Identifier: []string{"elide"}, Expression: "Text.ElideMiddle", Tags: map[string]string{}}}}, context)
	b.AddChild(child)

	context.RegisterComponent("", b)

	return b
}

func (b *Button) String() string {
	return fmt.Sprintf("Button(%s)", b.id)
}

func (b *Button) Property(key string) (vit.Value, bool) {
	switch key {
	case "text":
		return &b.text, true
	case "pressed":
		return &b.pressed, true
	case "fontPointSize":
		return &b.fontPointSize, true
	default:
		return b.Item.Property(key)
	}
}

func (b *Button) MustProperty(key string) vit.Value {
	v, ok := b.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (b *Button) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "text":
		err = b.text.SetValue(value)
	case "pressed":
		err = b.pressed.SetValue(value)
	case "fontPointSize":
		err = b.fontPointSize.SetValue(value)
	default:
		return b.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Button", key, b.id, err)
	}
	return nil
}

func (b *Button) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "text":
		b.text.SetCode(code)
	case "pressed":
		b.pressed.SetCode(code)
	case "fontPointSize":
		b.fontPointSize.SetCode(code)
	default:
		return b.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (b *Button) Event(name string) (vit.Listenable, bool) {
	switch name {
	case "onClicked":
		return &b.onClicked, true
	default:
		return b.Item.Event(name)
	}
}

func (b *Button) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "text":
		return &b.text, true
	case "pressed":
		return &b.pressed, true
	case "fontPointSize":
		return &b.fontPointSize, true
	case "clicked":
		return &b.clicked, true
	case "onClicked":
		return &b.onClicked, true
	default:
		return b.Item.ResolveVariable(key)
	}
}

func (b *Button) AddChild(child vit.Component) {
	child.SetParent(b)
	b.AddChildButKeepParent(child)
}

func (b *Button) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range b.Children() {
		if child.As(&targetType) {
			addThis.SetParent(b)
			b.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	b.AddChild(addThis)
}

func (b *Button) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = b
	}
	// properties
	if changed, err := b.text.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Button", "text", b.id, err))
		}
	}
	if changed, err := b.pressed.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Button", "pressed", b.id, err))
		}
	}
	if changed, err := b.fontPointSize.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Button", "fontPointSize", b.id, err))
		}
	}

	// methods
	if b.clicked.ShouldEvaluate() {
		_, err := b.clicked.Evaluate(b)
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Button", "clicked", b.id, err))
		}
	}

	n, err := b.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (b *Button) As(target *vit.Component) bool {
	if _, ok := (*target).(*Button); ok {
		*target = b
		return true
	}
	return b.Item.As(target)
}

func (b *Button) ID() string {
	return b.id
}

func (b *Button) Finish() error {
	return b.RootC().FinishInContext(b)
}
