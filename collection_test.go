package settings

import (
	"os"
	"os/user"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollectionNew(t *testing.T) {
	collection, err := NewCollection("settings", "testapp")
	assert.Nil(t, err)
	assert.NotNil(t, collection)
	assert.Equal(t, "settings", collection.Name())
	assert.Equal(t, defaultExtension, collection.Extension)

	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}

	os := runtime.GOOS
	if os == "windows" {
		assert.Equal(t, path.Join(user.HomeDir, ".config"), strings.Join(collection.SearchPaths, ","))
	} else {
		assert.Equal(t, "/etc,/usr/local/etc,"+path.Join(user.HomeDir, ".config"), strings.Join(collection.SearchPaths, ","))
	}
}

func TestCollectionLoad(t *testing.T) {
	collection, err := NewCollection("settings", "testapp")
	assert.Nil(t, err)
	assert.NotNil(t, collection)

	user, err := user.Current()
	if err != nil {
		t.Fatal(err)
	}
	dirpath := path.Join(user.HomeDir, ".config", "testapp")
	if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.RemoveAll(dirpath)
	}()
	err = os.WriteFile(path.Join(dirpath, "settings.yaml"), []byte("address: 127.0.0.1\nport: 8080"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}

	err = os.WriteFile("settings.yaml", []byte("address: 192.168.100.1\nport: 8181"), os.ModePerm)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove("settings.yaml")
	}()

	model := &struct {
		Address string `yaml:"address"`
		Port    int    `yaml:"port"`
	}{}

	err = collection.Load(model)
	assert.Nil(t, err)
	assert.Equal(t, "192.168.100.1", model.Address)
	assert.Equal(t, 8181, model.Port)
}
