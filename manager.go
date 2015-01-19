package deps

import (
	"io"
	"log"
	"os"
)

type Manager struct {
	Dir string
}

// Manages microservice dependencies, uses go tool approach
// meaning that it retrieves dependencies by using go-like import paths
func NewManager(dir string) *Manager {
	return &Manager{
		Dir: dir,
	}
}

var stdo = os.Stdout
var stde = os.Stderr

func (m *Manager) captureOutput(w io.Writer) (rerr error) {
	log.SetOutput(w)

	//stdout
	stdo = os.Stdout
	fr, fw, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stdout = fw
	go func() {
		_, rerr = io.Copy(w, fr)
	}()

	//stderr
	stde = os.Stderr
	fre, fwe, err := os.Pipe()
	if err != nil {
		return err
	}
	os.Stderr = fwe
	go func() {
		_, rerr = io.Copy(w, fre)
	}()

	return nil
}

func (m *Manager) restoreOutput() {
	os.Stdout = stdo
	os.Stderr = stde
	log.SetOutput(os.Stdout)
}

// install the dependency but if it already exists, do nothing,
// stdout is replaced by a stream we control, @todo this is not
// very elegant
func (m *Manager) Install(ipath string, w io.Writer) error {
	m.captureOutput(w)
	defer m.restoreOutput()

	//create package
	p := NewPackage(ipath)

	//download it
	return DownloadPackage(p, m.Dir, false)
}

// return the absolate path of a service by its import path
func (m *Manager) Locate(ipath string) (string, error) {

	//create package
	p := NewPackage(ipath)

	//@todo return more info
	_, _, _, root, err := ExpandPackage(p, m.Dir)
	if err != nil {
		return "", err
	}

	return root, nil
}

// update the dependency, if its not installed install
func (m *Manager) Upsert(ipath string, w io.Writer) error {
	m.captureOutput(w)
	defer m.restoreOutput()

	//create package
	p := NewPackage(ipath)

	//download it with update flag enabled
	return DownloadPackage(p, m.Dir, true)
}
