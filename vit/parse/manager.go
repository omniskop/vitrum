package parse

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

// bundles basic information about a vit file
type componentFile struct {
	name     string // name of the file without extension
	filePath string // full file path
}

// The Manager handles everything about loading files and instantiating them into a working component tree
type Manager struct {
	knownComponents   map[string]componentFile
	importPaths       []string
	mainComponentName string
	mainComponent     vit.Component
}

func NewManager() *Manager {
	return &Manager{
		knownComponents: make(map[string]componentFile),
	}
}

// AddImportPath adds a folder to the list of folders to search for components
func (m *Manager) AddImportPath(filePath string) error {
	entries, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}
	m.importPaths = append(m.importPaths, filePath)
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), ".vit") {
			// TODO: check if the file name is a valid component name
			name := strings.TrimSuffix(entry.Name(), ".vit")
			// TODO: check for collisions, either here or in Run
			m.knownComponents[name] = componentFile{
				name:     name,
				filePath: path.Join(filePath, entry.Name()),
			}
		}
	}
	return nil
}

// SetSource sets the primary component that should be instantiated
func (m *Manager) SetSource(filePath string) error {
	if m.mainComponentName != "" {
		return fmt.Errorf("main component already set")
	}
	if !strings.HasSuffix(filePath, ".vit") {
		return fmt.Errorf("not a vit file")
	}
	err := m.AddImportPath(path.Dir(filePath))
	if err != nil {
		return err
	}
	m.mainComponentName = strings.TrimSuffix(path.Base(filePath), ".vit")
	return nil
}

// Run instantiates the primary component and reports any errors in doing so
func (m *Manager) Run() error {
	var documents = vit.NewComponentContainer()
	var main VitDocument
	for _, cFile := range m.knownComponents {
		// TODO: maybe change ParseFile to operate on a componentFile?
		doc, err := parseFile(cFile.filePath, cFile.name)
		if err != nil {
			return err
		}
		documents.Set(cFile.name, &DocumentInstantiator{*doc})
		if cFile.name == m.mainComponentName {
			main = *doc
		}
	}

	documents = documents.ToGlobal()

	components, err := interpret(main, "", documents)
	if err != nil {
		return err
	}

	m.mainComponent = components[0]

	err = vit.FinishComponent(m.mainComponent)
	if err != nil {
		return err
	}
	evaluateStaticExpressions(documents)

evaluateExpressions:
	n, errs := m.Update()
	if errs.Failed() {
		return errs
	}
	if n > 0 {
		goto evaluateExpressions
	}

	return nil
}

// MainComponent returns the instantiated primary component
func (m *Manager) MainComponent() vit.Component {
	return m.mainComponent
}

// Update reevaluates all expressions whose dependencies have changed since the last update.
func (m *Manager) Update() (int, vit.ErrorGroup) {
	return m.mainComponent.UpdateExpressions()
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
