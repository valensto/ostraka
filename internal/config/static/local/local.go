package local

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/valensto/ostraka/internal/config"
	"github.com/valensto/ostraka/internal/logger"
	"gopkg.in/yaml.v3"
	"io"
	"os"
	"path/filepath"
)

func Extract(source string) ([]config.Workflow, error) {
	dir, err := os.ReadDir(source)
	if err != nil {
		return nil, fmt.Errorf("error reading resources directory: %w", err)
	}

	workflows := make([]config.Workflow, 0, len(dir))
	for _, file := range dir {
		fn := file.Name()

		ext := filepath.Ext(fn)
		if ext != ".yaml" && ext != ".yml" {
			logger.Get().Warn().Msgf(`file (%s) be skipped. No matching with authorized extensions (yaml | yml) found`, fn)
			continue
		}

		wf, err := extractBytes(filepath.Join(source, fn))
		if err != nil {
			return nil, fmt.Errorf("error extracting workflow: %w", err)
		}

		var cw config.Workflow
		err = yaml.Unmarshal(wf, &cw)
		if err != nil {
			return nil, fmt.Errorf("error parsing YAML wf: %w in file %s", err, fn)
		}

		err = validator.New().Struct(cw)
		if err != nil {
			return nil, fmt.Errorf("error: validating wf: %w in file %s", err, fn)
		}
	}

	return workflows, nil
}

func extractBytes(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filename, err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Get().Error().Msgf("error closing file %s: %v", filename, err)
		}
	}(f)

	b, err := io.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("error reading workflow file wf: %w", err)
	}

	return b, nil
}
