// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForImage(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type Image_FillMode uint

const (
	Image_FillMode_Fill            Image_FillMode = 0
	Image_FillMode_Fit             Image_FillMode = 1
	Image_FillMode_PreferUnchanged Image_FillMode = 2
)

func (enum Image_FillMode) String() string {
	switch enum {
	case Image_FillMode_Fill:
		return "Fill"
	case Image_FillMode_Fit:
		return "Fit"
	case Image_FillMode_PreferUnchanged:
		return "PreferUnchanged"
	default:
		return "<unknownFillMode>"
	}
}

type Image struct {
	*Item
	id string

	path      vit.StringValue
	fillMode  vit.IntValue
	imageData *img
}

// newImageInGlobal creates an appropriate file context for the component and then returns a new Image instance.
// The returned error will only be set if a library import that is required by the component fails.
func newImageInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*Image, error) {
	fileCtx, err := newFileContextForImage(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewImage(id, fileCtx), nil
}
func NewImage(id string, context *vit.FileContext) *Image {
	i := &Image{
		Item:      NewItem("", context),
		id:        id,
		path:      *vit.NewEmptyStringValue(),
		fillMode:  *vit.NewIntValueFromCode(vit.Code{FileCtx: context, Code: "FillMode.Fit", Position: nil}),
		imageData: nil,
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	i.path.AddDependent(vit.FuncDep(i.reloadImage))
	// register event listeners
	// register enumerations
	i.DefineEnum(vit.Enumeration{
		Embedded: true,
		Name:     "FillMode",
		Position: nil,
		Values:   map[string]int{"Fill": 0, "Fit": 1, "PreferUnchanged": 2},
	})
	// add child components

	context.RegisterComponent("", i)

	return i
}

func (i *Image) String() string {
	return fmt.Sprintf("Image(%s)", i.id)
}

func (i *Image) Property(key string) (vit.Value, bool) {
	switch key {
	case "path":
		return &i.path, true
	case "fillMode":
		return &i.fillMode, true
	default:
		return i.Item.Property(key)
	}
}

func (i *Image) MustProperty(key string) vit.Value {
	v, ok := i.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (i *Image) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "path":
		err = i.path.SetValue(value)
	case "fillMode":
		err = i.fillMode.SetValue(value)
	default:
		return i.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("Image", key, i.id, err)
	}
	return nil
}

func (i *Image) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "path":
		i.path.SetCode(code)
	case "fillMode":
		i.fillMode.SetCode(code)
	default:
		return i.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (i *Image) Event(name string) (vit.Listenable, bool) {
	switch name {
	default:
		return i.Item.Event(name)
	}
}

func (i *Image) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "path":
		return &i.path, true
	case "fillMode":
		return &i.fillMode, true
	default:
		return i.Item.ResolveVariable(key)
	}
}

func (i *Image) AddChild(child vit.Component) {
	child.SetParent(i)
	i.AddChildButKeepParent(child)
}

func (i *Image) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range i.Children() {
		if child.As(&targetType) {
			addThis.SetParent(i)
			i.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	i.AddChild(addThis)
}

func (i *Image) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = i
	}
	// properties
	if changed, err := i.path.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Image", "path", i.id, err))
		}
	}
	if changed, err := i.fillMode.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("Image", "fillMode", i.id, err))
		}
	}

	// methods

	n, err := i.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (i *Image) As(target *vit.Component) bool {
	if _, ok := (*target).(*Image); ok {
		*target = i
		return true
	}
	return i.Item.As(target)
}

func (i *Image) ID() string {
	return i.id
}

func (i *Image) Finish() error {
	return i.RootC().FinishInContext(i)
}

func (i *Image) staticAttribute(name string) (interface{}, bool) {
	switch name {
	case "Fill":
		return uint(Image_FillMode_Fill), true
	case "Fit":
		return uint(Image_FillMode_Fit), true
	case "PreferUnchanged":
		return uint(Image_FillMode_PreferUnchanged), true
	default:
		return nil, false
	}
}
