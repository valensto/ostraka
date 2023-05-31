package local

import (
	"context"
	"fmt"
	"github.com/valensto/ostraka/internal/config/static"
	"github.com/valensto/ostraka/internal/logger"
	"github.com/valensto/ostraka/internal/workflow"
	"io"
	"os"
	"path/filepath"
)

type Local struct {
	source string
}

func New(source string) Local {
	return Local{source: source}
}

func (l Local) Extract(_ context.Context) ([]*workflow.Workflow, error) {
	dir, err := os.ReadDir(l.source)
	if err != nil {
		return nil, fmt.Errorf("error reading resources directory: %w", err)
	}

	var workflows [][]byte
	for _, file := range dir {
		ext := filepath.Ext(file.Name())
		if ext != ".yaml" && ext != ".yml" {
			logger.Get().Warn().Msgf(`file (%s) be skipped. No matching with authorized extensions (yaml | yml) found`, file.Name())
			continue
		}

		wf, err := extractBytes(filepath.Join(l.source, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("error extracting workflow: %w", err)
		}
		workflows = append(workflows, wf)
	}

	return static.BuildWorkflows(workflows)
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
