package extractor

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
)

type (
	Extension string
	Files     map[Extension][]byte
)

var Extensions = map[Extension]struct{}{
	".json": {},
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
	default:
		return fmt.Errorf("extension %s not supported", e)
	}
}
