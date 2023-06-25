package config

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

type (
	Extension string
	Files     map[Extension][]byte

	Provider interface {
		Extract() (Files, error)
	}
)

var Extensions = map[Extension]struct{}{
	".json": {},
	".toml": {},
	".yaml": {},
	".yml":  {},
}

func LookupExtension(ext string) (Extension, bool) {
	extension := Extension(ext)
	_, ok := Extensions[extension]
	return extension, ok
}

func (e Extension) String() string {
	return string(e)
}

func (e Extension) Unmarshal(in []byte, out interface{}) error {
	switch e {
	case ".json":
		return json.Unmarshal(in, out)
	case ".yaml", ".yml":
		return yaml.Unmarshal(in, out)
	case ".toml":
		return toml.Unmarshal(in, out)
	default:
		return fmt.Errorf("extension %s not supported", e)
	}
}
