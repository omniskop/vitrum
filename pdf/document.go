package pdf

import (
	"fmt"
	"io"
	"log"

	vit "github.com/omniskop/vitrum/vit"
	parse "github.com/omniskop/vitrum/vit/parse"
	"github.com/omniskop/vitrum/vit/vpath"
	"github.com/tdewolff/canvas"
	"github.com/tdewolff/canvas/renderers/pdf"
)

type componentHandler struct {
	logger *log.Logger
}

func (h componentHandler) RegisterComponent(id string, comp vit.Component) {}

func (h componentHandler) UnregisterComponent(id string, comp vit.Component) {}

func (h componentHandler) RequestFocus(comp vit.FocusableComponent) {}

func (h componentHandler) Logger() *log.Logger {
	return h.logger
}

type Document struct {
	manager *parse.Manager
	handler *componentHandler
}

func NewDocument(path vpath.Path) (*Document, error) {
	doc := &Document{
		manager: parse.NewManager(),
		handler: &componentHandler{log.New(io.Discard, "", 0)},
	}

	doc.manager.SetSource(path)
	err := doc.manager.Initialize(doc.handler)
	if err != nil {
		return nil, err
	}

	return doc, err
}

func (d *Document) AddImportPath(path vpath.Path) {
	d.manager.AddImportPath(path)
}

func (d *Document) SetLogger(log *log.Logger) {
	d.handler.logger = log
}

func (d *Document) Render(out io.Writer) error {
	errs := d.manager.UpdateFully()
	if errs.Failed() {
		return errs
	}

	var c vit.Component = &DocumentComponent{}
	if ok := d.manager.MainComponent().As(&c); !ok {
		return fmt.Errorf("main component is not a Document")
	}
	docComp := c.(*DocumentComponent)

	if len(docComp.Children()) == 0 {
		return fmt.Errorf("document has no pages")
	}

	firstPageBounds := docComp.Children()[0].Bounds()

	p := pdf.New(out, firstPageBounds.Width(), firstPageBounds.Height(), &pdf.Options{
		Compress:      true,
		SubsetFonts:   true,
		ImageEncoding: canvas.Lossless,
	})

	for i, child := range docComp.Children() {
		bounds := child.Bounds()
		if i != 0 {
			p.NewPage(bounds.Width(), bounds.Height())
		}
		canv := canvas.New(bounds.Width(), bounds.Height())
		ctx := canvas.NewContext(canv)
		ctx.SetCoordSystem(canvas.CartesianIV)
		err := child.Draw(
			vit.DrawingContext{ctx},
			bounds,
		)
		if err != nil {
			return fmt.Errorf("page %d: %w", i, err)
		}
		canv.Render(p) // This has been changed to RenderTo in newer versions of the canvas library
	}

	err := p.Close()
	if err != nil {
		return fmt.Errorf("pdf: %v", err)
	}

	return nil
}
