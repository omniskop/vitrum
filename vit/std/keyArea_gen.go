// Code generated by vitrum gencmd. DO NOT EDIT.

package std

import (
	"fmt"
	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
)

func newFileContextForKeyArea(globalCtx *vit.GlobalContext) (*vit.FileContext, error) {
	return vit.NewFileContext(globalCtx), nil
}

type KeyArea struct {
	*Item
	id string

	enabled vit.BoolValue
	pressed vit.BoolValue

	onKeyDown vit.EventAttribute[KeyEvent]
	onKeyUp   vit.EventAttribute[KeyEvent]
}

// newKeyAreaInGlobal creates an appropriate file context for the component and then returns a new KeyArea instance.
// The returned error will only be set if a library import that is required by the component fails.
func newKeyAreaInGlobal(id string, globalCtx *vit.GlobalContext, thisLibrary parse.Library) (*KeyArea, error) {
	fileCtx, err := newFileContextForKeyArea(globalCtx)
	if err != nil {
		return nil, err
	}
	parse.AddLibraryToContainer(thisLibrary, &fileCtx.KnownComponents)
	return NewKeyArea(id, fileCtx), nil
}
func NewKeyArea(id string, context *vit.FileContext) *KeyArea {
	k := &KeyArea{
		Item:      NewItem("", context),
		id:        id,
		enabled:   *vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "true", Position: nil}),
		pressed:   *vit.NewBoolValueFromCode(vit.Code{FileCtx: context, Code: "false", Position: nil}),
		onKeyDown: *vit.NewEventAttribute[KeyEvent](),
		onKeyUp:   *vit.NewEventAttribute[KeyEvent](),
	}
	// property assignments on embedded components
	// register listeners for when a property changes
	k.enabled.AddDependent(vit.FuncDep(k.enableDisable))
	// register event listeners
	// register enumerations
	// add child components

	context.RegisterComponent("", k)

	return k
}

func (k *KeyArea) String() string {
	return fmt.Sprintf("KeyArea(%s)", k.id)
}

func (k *KeyArea) Property(key string) (vit.Value, bool) {
	switch key {
	case "enabled":
		return &k.enabled, true
	case "pressed":
		return &k.pressed, true
	default:
		return k.Item.Property(key)
	}
}

func (k *KeyArea) MustProperty(key string) vit.Value {
	v, ok := k.Property(key)
	if !ok {
		panic(fmt.Errorf("MustProperty called with unknown key %q", key))
	}
	return v
}

func (k *KeyArea) SetProperty(key string, value interface{}) error {
	var err error
	switch key {
	case "enabled":
		err = k.enabled.SetValue(value)
	case "pressed":
		err = k.pressed.SetValue(value)
	default:
		return k.Item.SetProperty(key, value)
	}
	if err != nil {
		return vit.NewPropertyError("KeyArea", key, k.id, err)
	}
	return nil
}

func (k *KeyArea) SetPropertyCode(key string, code vit.Code) error {
	switch key {
	case "enabled":
		k.enabled.SetCode(code)
	case "pressed":
		k.pressed.SetCode(code)
	default:
		return k.Item.SetPropertyCode(key, code)
	}
	return nil
}

func (k *KeyArea) Event(name string) (vit.Listenable, bool) {
	switch name {
	case "onKeyDown":
		return &k.onKeyDown, true
	case "onKeyUp":
		return &k.onKeyUp, true
	default:
		return k.Item.Event(name)
	}
}

func (k *KeyArea) ResolveVariable(key string) (interface{}, bool) {
	switch key {
	case "enabled":
		return &k.enabled, true
	case "pressed":
		return &k.pressed, true
	case "onKeyDown":
		return &k.onKeyDown, true
	case "onKeyUp":
		return &k.onKeyUp, true
	default:
		return k.Item.ResolveVariable(key)
	}
}

func (k *KeyArea) AddChild(child vit.Component) {
	child.SetParent(k)
	k.AddChildButKeepParent(child)
}

func (k *KeyArea) AddChildAfter(afterThis vit.Component, addThis vit.Component) {
	var targetType vit.Component = afterThis

	for ind, child := range k.Children() {
		if child.As(&targetType) {
			addThis.SetParent(k)
			k.AddChildAtButKeepParent(addThis, ind+1)
			return
		}
	}
	k.AddChild(addThis)
}

func (k *KeyArea) UpdateExpressions(context vit.Component) (int, vit.ErrorGroup) {
	var sum int
	var errs vit.ErrorGroup

	if context == nil {
		context = k
	}
	// properties
	if changed, err := k.enabled.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("KeyArea", "enabled", k.id, err))
		}
	}
	if changed, err := k.pressed.Update(context); changed || err != nil {
		sum++
		if err != nil {
			errs.Add(vit.NewPropertyError("KeyArea", "pressed", k.id, err))
		}
	}

	// methods

	n, err := k.Item.UpdateExpressions(context)
	sum += n
	errs.AddGroup(err)
	return sum, errs
}

func (k *KeyArea) As(target *vit.Component) bool {
	if _, ok := (*target).(*KeyArea); ok {
		*target = k
		return true
	}
	return k.Item.As(target)
}

func (k *KeyArea) ID() string {
	return k.id
}

func (k *KeyArea) Finish() error {
	return k.RootC().FinishInContext(k)
}
