package yaml

import (
	"iac/utils/errors"
	"iac/utils/fs"

	"github.com/goccy/go-yaml"
)

func FromFile(path string, output any) error {
	bytes, err := fs.ReadFileAsBytes(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(bytes, output)
	if err != nil {
		return errors.New("Could not unmarshal YAML from file", "path", path, "error", err)
	}

	return nil
}
