package menu

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
)

// Param struct: a parameter for a action
type Param struct {
	Name  string `json:"name" binding:"required"`
	Help  string `json:"help"`
	Value string `json:"value"`
}

// Action struct: a thing a user can do
type Action struct {
	Name     string  `json:"name" binding:"required"`
	Help     string  `json:"help"`
	Template string  `json:"template" binding:"required"`
	Params   []Param `json:"params"`
}

// Menu struct: a list of actions
type Menu struct {
	Name    string   `json:"name" binding:"required"`
	Version string   `json:"version"`
	Help    string   `json:"help"`
	Actions []Action `json:"actions" binding:"required"`
	Hash    string   `json:"hash"`
}

// Home retrieves path of kut home e.g. /home/james/.kut
func Home() string {
	kutHome := os.Getenv("KUT_HOME")
	if kutHome == "" {
		if userHome, err := homedir.Dir(); err != nil {
			log.Fatal(err)
		} else {
			kutHome = filepath.Join(userHome, ".kut")
		}
	}

	if err := os.MkdirAll(kutHome, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	return kutHome
}

// Repo retrieves path of local repo e.g. /home/james/.kut/menus
func Repo() string {
	repo := filepath.Join(Home(), "menus")

	if err := os.MkdirAll(repo, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	return repo
}

// Path generates the path to the named menu
func Path(name string) string {
	return filepath.Join(Repo(), name+".yaml")
}

// List generates available menus
func List() []string {
	const suffix = ".yaml"
	var files []os.FileInfo
	var err error
	var ls []string

	if files, err = ioutil.ReadDir(Repo()); err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), suffix) {
			ls = append(ls, strings.TrimSuffix(f.Name(), suffix))
		}
	}

	return ls
}

// Ingest a file into system store
func Ingest(source string) error {
	var data []byte
	var err error
	var parsed *url.URL
	var resp *http.Response
	var menu Menu

	parsed, err = url.Parse(source)
	if strings.HasPrefix(parsed.Scheme, "http") {
		if resp, err = http.Get(source); err != nil {
			return err
		}
		defer resp.Body.Close()
		if data, err = ioutil.ReadAll(resp.Body); err != nil {
			return err
		}
	} else {
		if data, err = ioutil.ReadFile(source); err != nil {
			return err
		}
	}

	if menu, err = Parse(data); err != nil {
		return err
	}

	return ioutil.WriteFile(Path(menu.Name), data, 0644)
}

// Load a menu given a name
func Load(name string) (Menu, error) {
	var data []byte
	var err error

	if data, err = ioutil.ReadFile(Path(name)); err != nil {
		return Menu{}, err
	}

	return Parse(data)
}

// Parse a menu given a blob
func Parse(data []byte) (Menu, error) {
	menu := &Menu{}
	var err error

	if err = yaml.Unmarshal(data, menu); err != nil {
		return *menu, err
	}

	if menu.Actions == nil {
		menu.Actions = []Action{}
	}

	for i, _ := range menu.Actions {
		if menu.Actions[i].Params == nil {
			menu.Actions[i].Params = []Param{}
		}
	}

	hash16 := md5.Sum(data)
	menu.Hash = hex.EncodeToString(hash16[:])[:4]

	return *menu, err
}

// Render an action into a command
func Render(action Action) (string, error) {
	var tpl template.Template
	var err error

	_, err = tpl.Parse(action.Template)
	if err != nil {
		return "", err
	}

	paramMap := make(map[string]string)
	for _, param := range action.Params {
		if param.Value != "" {
			paramMap[param.Name] = param.Value
		}
	}

	out := bytes.Buffer{}
	err = tpl.Execute(&out, paramMap)
	if err != nil {
		return "", err
	}

	return out.String(), nil
}
