package settings

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

const defaultExtension = "yaml"

type Collection struct {
	name string

	Extension string
	Paths     []string
}

func (cl *Collection) Name() string {
	return cl.name
}

func (cl *Collection) Load(v interface{}) error {
	if len(cl.Paths) == 0 {
		return nil
	}

	var filepath string
	for _, fpath := range cl.Paths {
		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			continue
		}
		filepath = fpath
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, v)
}

func NewCollection(dirname, name string) (*Collection, error) {
	collection := &Collection{
		name:      name,
		Extension: defaultExtension,
	}
	filename := fmt.Sprintf("%s.%s", name, collection.Extension)

	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	os := runtime.GOOS

	if os == "windows" {
		collection.Paths = []string{
			path.Join(user.HomeDir, dirname, filename),
			filename,
		}
	} else {
		collection.Paths = []string{
			path.Join("/etc", dirname, filename),
			path.Join(user.HomeDir, fmt.Sprintf(".%s", dirname), filename),
			filename,
		}
	}

	return collection, nil
}
