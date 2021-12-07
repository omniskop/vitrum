package parse

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/omniskop/vitrum/vit"
)

type componentFile struct {
	name     string
	filePath string
}

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

func (m *Manager) Run() error {
	var documents = make(map[string]vit.AbstractComponent)
	var main VitDocument
	for _, cFile := range m.knownComponents {
		// TODO: maybe change ParseFile to operate on a componentFile?
		doc, err := ParseFile(cFile.filePath, cFile.name)
		if err != nil {
			return err
		}
		documents[cFile.name] = &DocumentInstantiator{*doc}
		if cFile.name == m.mainComponentName {
			main = *doc
		}
	}

	components, err := Interpret(main, vit.ComponentResolver{nil, documents})
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

func (m *Manager) MainComponent() vit.Component {
	return m.mainComponent
}

// Update reevaluates all expressions whose dependencies have changed since the last update.
func (m *Manager) Update() (int, error) {
	return m.mainComponent.UpdateExpressions()
}
