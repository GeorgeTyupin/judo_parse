package filesutils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetRootFilePath(fileName string) (string, error) {
	exe, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("не удалось определить путь к бинарнику: %w", err)
	}
	p := filepath.Join(filepath.Dir(exe), fileName)
	if _, err := os.Stat(p); err == nil {
		return p, nil
	}

	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("не удалось определить рабочую директорию: %w", err)
	}
	dir := wd
	for {
		p = filepath.Join(dir, fileName)
		if _, err := os.Stat(p); err == nil {
			return p, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("файл %s не найден", fileName)
}
