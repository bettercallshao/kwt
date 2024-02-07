package menu

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ghodss/yaml"
	"gotest.tools/assert"
)

// Save a menu to a named file
func Save(name string, menu *Menu) error {
	var data []byte
	var err error

	if data, err = yaml.Marshal(menu); err != nil {
		return err
	}

	return os.WriteFile(name, data, 0644)
}

func TestIngestListLoadHappyPath(t *testing.T) {
	home, _ := os.MkdirTemp("", "")
	os.Setenv("KWT_HOME", home)
	defer os.RemoveAll(home)

	name := "test"
	menu := &Menu{
		Name: name,
	}

	source := filepath.Join(home, "source.yaml")
	Save(source, menu)

	Ingest(source)

	ls := List()
	if ls[0] != name {
		t.Fail()
	}

	other, _ := Load(name)
	if other.Name != menu.Name {
		t.Fail()
	}
	if other.Actions == nil {
		t.Fail()
	}
	if other.Hash != "7d9f" {
		t.Fail()
	}
}

func TestRender(t *testing.T) {
	action := Action{
		Name:     "jim",
		Template: "kubectl {{if .namespace}}-n {{.namespace}} {{end}}get all",
		Params: []Param{
			{
				Name: "namespace",
			},
		},
	}

	out1, _ := Render(action)
	assert.Equal(t, out1, "kubectl get all")

	action.Params[0].Value = "jim"
	out2, _ := Render(action)
	assert.Equal(t, out2, "kubectl -n jim get all")
}
