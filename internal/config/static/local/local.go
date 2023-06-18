package local

import (
	"fmt"
	"github.com/valensto/ostraka/internal/config/static"
	"github.com/valensto/ostraka/internal/logger"
	"io"
	"os"
	"path/filepath"
)

func Extract(source string) (static.ContentFile, error) {
	dir, err := os.ReadDir(source)
	if err != nil {
		return nil, fmt.Errorf("error reading resources directory: %w", err)
	}

	contentFile := make(static.ContentFile)
	for _, file := range dir {
		fn := file.Name()

		ext := filepath.Ext(fn)
		if ext != ".yaml" && ext != ".yml" {
			logger.Get().Warn().Msgf(`file (%s) be skipped. No matching with authorized extensions (yaml | yml) found`, fn)
			continue
		}

		b, err := extractBytes(filepath.Join(source, fn))
		if err != nil {
			return nil, fmt.Errorf("error extracting workflow: %w", err)
		}

		contentFile[fn] = b
	}

	return contentFile, nil
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
