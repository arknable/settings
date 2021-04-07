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

// Collection is a group of settings which are saved in a single yaml file.
type Collection struct {
	name string

	// Extension is settings file extension, default is 'yaml'.
	Extension string

	// Paths is list of directories to look for settings file.
	SearchPaths []string
}

// Name returns name of this collection.
func (cl *Collection) Name() string {
	return cl.name
}

// Load read settings file and unmarshal the content as given v.
// For example, if a collection created with directory name 'test' and named 'settings',
// with default extension, then this function will look for SEARCH_DIR/test/settings.yaml.
func (cl *Collection) Load(v interface{}) error {
	if len(cl.SearchPaths) == 0 {
		return nil
	}

	var filepath string
	for _, fpath := range cl.SearchPaths {
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

// NewCollection create new collection given its name and directory name.
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
		collection.SearchPaths = []string{
			path.Join(user.HomeDir, dirname, filename),
			filename,
		}
	} else {
		collection.SearchPaths = []string{
			path.Join("/etc", dirname, filename),
			path.Join(user.HomeDir, fmt.Sprintf(".%s", dirname), filename),
			filename,
		}
	}

	return collection, nil
}
