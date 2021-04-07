package settings

import (
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"

	"gopkg.in/yaml.v2"
)

const defaultExtension = "yaml"

// ErrCollectionNotFound occurred when settings file cannot be found in search paths.
var ErrCollectionNotFound = errors.New("collection not found")

// Collection is a group of settings which are saved in a single yaml file.
type Collection struct {
	name    string
	dirname string

	// Extension is settings file extension, default is 'yaml'.
	Extension string

	// Paths is list of directories to look for settings file.
	// Default values are as follows,
	// - $HOME/.config
	// - /etc				(Non Windows)
	// - /usr/local/etc		(Non Windows)
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

	var (
		filepath string
		found    = false
		pathList = cl.SearchPaths
	)
	pathList = append(pathList, "")

	for _, fpath := range pathList {
		if fpath == "" {
			filepath = fmt.Sprintf("%s.%s", cl.name, cl.Extension)
		} else {
			filepath = path.Join(fpath, cl.dirname, fmt.Sprintf("%s.%s", cl.name, cl.Extension))
		}

		if _, err := os.Stat(fpath); os.IsNotExist(err) {
			continue
		}
		found = true
	}

	if !found {
		return ErrCollectionNotFound
	}

	content, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}

	return yaml.Unmarshal(content, v)
}

// NewCollection create new collection given its name and directory name.
func NewCollection(name, dirname string) (*Collection, error) {
	collection := &Collection{
		name:      name,
		dirname:   dirname,
		Extension: defaultExtension,
	}

	user, err := user.Current()
	if err != nil {
		return nil, err
	}
	os := runtime.GOOS

	if os != "windows" {
		collection.SearchPaths = append(collection.SearchPaths, []string{
			"/etc",
			"/usr/local/etc",
		}...)
	}
	collection.SearchPaths = append(collection.SearchPaths, path.Join(user.HomeDir, ".config"))

	return collection, nil
}
