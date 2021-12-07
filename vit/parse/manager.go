package parse

import (
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
	var documents = make(map[string]vit.AbstractComponent)
	var main VitDocument
	for _, cFile := range m.knownComponents {
		// TODO: maybe change ParseFile to operate on a componentFile?
		doc, err := parseFile(cFile.filePath, cFile.name)
		if err != nil {
			return err
		}
		documents[cFile.name] = &DocumentInstantiator{*doc}
		if cFile.name == m.mainComponentName {
			main = *doc
		}
	}

	components, err := interpret(main, vit.ComponentResolver{Parent: nil, Components: documents})
	if err != nil {
		return err
	}

	m.mainComponent = components[0]

evaluateExpressions:
	n, err := m.Update()
	if err != nil {
		return err
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
func (m *Manager) Update() (int, error) {
	return m.mainComponent.UpdateExpressions()
}
