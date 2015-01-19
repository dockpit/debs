package debs_test

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/dockpit/debs"
)

func TestInstall(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	dir, err := ioutil.TempDir("", "dockpit_test")
	if err != nil {
		t.Fatal(err)
	}

	m := debs.NewManager(dir)

	err = m.Install("github.com/golang/example", buff)
	if err != nil {
		t.Fatal(err)
	}

	//check if we can open it
	fname := filepath.Join(dir, "github.com", "golang", "example", "LICENSE")
	_, err = ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	//corrupt the git repository
	err = ioutil.WriteFile(filepath.Join(dir, "github.com", "golang", "example", ".git", "HEAD"), []byte("a"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	//another call to install
	err = m.Install("github.com/golang/example", buff)
	if err != nil {
		t.Fatal(err)
	}

	// should not throw error since no git interaction is expected
	assert.Equal(t, nil, err)
	assert.Contains(t, buff.String(), "git")
}

func TestLocate(t *testing.T) {

	dir, err := ioutil.TempDir("", "dockpit_test")
	if err != nil {
		t.Fatal(err)
	}

	m := debs.NewManager(dir)

	loc, err := m.Locate("github.com/golang/example")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, filepath.Join(dir, "github.com", "golang", "example"), loc)
}

func TestUpsert(t *testing.T) {
	buff := bytes.NewBuffer(nil)
	dir, err := ioutil.TempDir("", "dockpit_test")
	if err != nil {
		t.Fatal(err)
	}

	m := debs.NewManager(dir)

	err = m.Upsert("github.com/golang/example", buff)
	if err != nil {
		t.Fatal(err)
	}

	//check if we can open it
	fname := filepath.Join(dir, "github.com", "golang", "example", "LICENSE")
	_, err = ioutil.ReadFile(fname)
	if err != nil {
		t.Fatal(err)
	}

	//corrupt the git repository
	err = ioutil.WriteFile(filepath.Join(dir, "github.com", "golang", "example", ".git", "HEAD"), []byte("a"), 0777)
	if err != nil {
		t.Fatal(err)
	}

	//another call to Upsert
	err = m.Upsert("github.com/golang/example", buff)

	// should throw error since git interaction is expected
	assert.NotEqual(t, nil, err)
}
