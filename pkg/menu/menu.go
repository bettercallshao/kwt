package menu

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/ghodss/yaml"
	"github.com/mitchellh/go-homedir"
	"github.com/yargevad/filepathx"
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
	Path    string   `json:"path"`
}

// Home retrieves path of kwt home e.g. /home/james/.kwt
func Home() string {
	kwtHome := os.Getenv("KWT_HOME")
	if kwtHome == "" {
		if userHome, err := homedir.Dir(); err != nil {
			log.Fatal(err)
		} else {
			kwtHome = filepath.Join(userHome, ".kwt")
		}
	}

	if err := os.MkdirAll(kwtHome, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	return kwtHome
}

// Pwd path used to find menus i.e. .kwt
func Pwd() string {
	return ".kwt"
}

// Ingest a file into system store
func Ingest(source string) error {
	var data []byte
	var err error
	var parsed *url.URL
	var resp *http.Response
	var menu Menu
	var path string

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

	// we can't handle '/' character in file name yet
	path = filepath.Join(Home(), strings.ReplaceAll(menu.Name, "/", "+")+".yaml")

	return ioutil.WriteFile(path, data, 0644)
}

// MapRepo loads available menus from repo
func MapRepo(repo string) map[string]Menu {
	var pattern = repo + "/**/*.yaml"
	var paths []string
	var err error
	var menuMap map[string]Menu
	var data []byte
	var menu Menu

	menuMap = make(map[string]Menu)

	if paths, err = filepathx.Glob(pattern); err != nil {
		log.Fatal(err)
	}

	for _, path := range paths {
		if data, err = ioutil.ReadFile(path); err != nil {
			// silently ignore bad files
			continue
		}
		if menu, err = Parse(data); err != nil {
			continue
		}
		// ignore duplicate keys for now
		menu.Path = path
		menuMap[menu.Name] = menu
	}

	return menuMap
}

// Map loads all available menus
func Map() map[string]Menu {
	homeMap := MapRepo(Home())
	pwdMap := MapRepo(Pwd())
	for name, menu := range pwdMap {
		homeMap[name] = menu
	}
	return homeMap
}

func (m *Menu) isPwd() bool {
	return strings.HasPrefix(m.Path, Pwd())
}

// List generates list of menu names, home first, pwd after
func List() []string {
	menuMap := Map()
	names := make([]string, len(menuMap))

	i := 0
	for name := range menuMap {
		names[i] = name
		i++
	}

	sort.Slice(names, func(a, b int) bool {
		aMenu := menuMap[names[a]]
		bMenu := menuMap[names[b]]
		if aMenu.isPwd() != bMenu.isPwd() {
			return bMenu.isPwd()
		}
		return aMenu.Name < bMenu.Name
	})
	return names
}

// Load loads the menu identified by name
func Load(name string) (Menu, error) {
	var err error

	menuMap := Map()
	menu, found := menuMap[name]
	if !found {
		err = errors.New("")
	}

	return menu, err
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
