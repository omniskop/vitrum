package parse

import (
	"errors"
	"fmt"
	"io/fs"
	"path"
	"strings"

	"github.com/omniskop/vitrum/vit"
	"github.com/omniskop/vitrum/vit/vpath"
)

// bundles basic information about a vit file
type componentFile struct {
	name string     // name of the file without extension
	path vpath.Path // full file path
}

// The Manager handles everything about loading files and instantiating them into a working component tree
type Manager struct {
	knownComponents   map[string]componentFile
	importPaths       []fs.ReadDirFS
	mainComponentName string
	mainComponent     vit.Component
	globalCtx         vit.GlobalContext
}

func NewManager() *Manager {
	return &Manager{
		knownComponents: make(map[string]componentFile),
		globalCtx: vit.GlobalContext{
			Variables: make(map[string]vit.Value),
		},
	}
}

// AddImportPath adds a folder to the list of folders to search for components
func (m *Manager) AddImportPath(dir fs.ReadDirFS) error {
	entries, err := dir.ReadDir(".")
	if err != nil {
		return err
	}
	m.importPaths = append(m.importPaths, dir)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".vit") {
			// TODO: check if the file name is a valid component name
			name := strings.TrimSuffix(entry.Name(), ".vit")
			// TODO: check for collisions, either here or in Run
			m.knownComponents[name] = componentFile{
				name: name,
				path: vpath.FS(dir, entry.Name()),
			}
		}
	}
	return nil
}

// SetSource sets the primary component that should be instantiated
func (m *Manager) SetSource(filePath vpath.Path) error {
	if m.mainComponentName != "" {
		return fmt.Errorf("main component already set")
	}
	if !strings.HasSuffix(filePath.Path(), ".vit") {
		return fmt.Errorf("not a vit file")
	}
	err := m.AddImportPath(filePath.Dir())
	if err != nil {
		return err
	}
	m.mainComponentName = strings.TrimSuffix(path.Base(filePath.Path()), ".vit")
	return nil
}

// Initialize instantiates the primary component and reports any errors in doing so
func (m *Manager) Initialize(environment vit.ExecutionEnvironment) error {
	var documents = vit.NewComponentContainer()
	var main *VitDocument
	for _, cFile := range m.knownComponents {
		// TODO: maybe change ParseFile to operate on a componentFile?
		doc, err := parseFile(cFile.path, cFile.name)
		if err != nil {
			return err
		}
		documents.Set(cFile.name, &DocumentInstantiator{*doc})
		if cFile.name == m.mainComponentName {
			main = doc
		}
	}

	if main == nil {
		return fmt.Errorf("main component %q not found. Is the source file missing?", m.mainComponentName)
	}

	documents = documents.ToGlobal()

	m.globalCtx.KnownComponents = documents
	m.globalCtx.Environment = environment

	components, err := interpret(*main, "", &m.globalCtx)
	if err != nil {
		return err
	}

	m.mainComponent = components[0]

	err = vit.FinishComponent(m.mainComponent)
	if err != nil {
		return err
	}
	evaluateStaticExpressions(documents)

	return nil
}

func (m *Manager) SetVariable(name string, value interface{}) error {
	return m.globalCtx.SetVariable(name, value)
}

// MainComponent returns the instantiated primary component
func (m *Manager) MainComponent() vit.Component {
	return m.mainComponent
}

// UpdateOnce reevaluates all expressions whose dependencies have changed since the last update.
func (m *Manager) UpdateOnce() (int, vit.ErrorGroup) {
	return m.mainComponent.UpdateExpressions(nil)
}

// UpdateFully reevaluates all expressions whose dependencies have changed since the last update in a loop until no outstanding changes are left.
func (m *Manager) UpdateFully() vit.ErrorGroup {
evaluateExpressions:
	n, errs := m.mainComponent.UpdateExpressions(nil)
	if errs.Failed() {
		return errs
	}
	if n > 0 {
		goto evaluateExpressions
	}
	return errs
}

// FormatError takes an error that has been returned and formats it nicely for printing
func FormatError(err error) string {
	var group vit.ErrorGroup
	if errors.As(err, &group) {
		if len(group.Errors) == 1 {
			return FormatError(group.Errors[0])
		}
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%d Errors:\r\n", len(group.Errors)))
		for _, e := range group.Errors {
			b.WriteString(e.Error())
			b.WriteString("\r\n")
		}
		return b.String()
	}

	var gErr genericError
	if errors.As(err, &gErr) {
		// If the error chain contains a genericError, we will grab it's position and start the line with that.
		// The generic error itself will not print it's position.
		// We will also generate a report to add more information about the error location.
		return fmt.Sprintf("%v: %v\r\n\r\n%s", gErr.position, err, gErr.position.Report())
	}

	var pErr ParseError
	if errors.As(err, &pErr) {
		// If the error chain contains a parseError we will use it to create a report.
		return fmt.Sprintf("%v\r\n\r\n%s", err, pErr.Report())
	}

	var eErr vit.ExpressionError
	if errors.As(err, &eErr) && eErr.Position != nil {
		// same as ParseError
		return fmt.Sprintf("%v\r\n\r\n%s", err, eErr.Position.Report())
	}

	return err.Error()
}
