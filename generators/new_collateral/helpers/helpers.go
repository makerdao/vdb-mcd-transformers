package helpers

import (
	"os"
	"path/filepath"
)

func GetProjectPath() string {
	projectPath := filepath.Join("$GOPATH", "src", "github.com", "makerdao", "vdb-mcd-transformers")
	return os.ExpandEnv(projectPath)
}

func GetEnvironmentsPath() string {
	return filepath.Join(GetProjectPath(), "environments")
}

func GetExecutePluginsPath() string {
	return filepath.Join(GetProjectPath(), "plugins", "execute")
}

func GetFullConfigFilePath(filePath, fileName string) string {
	fileNameWithExtension := fileName + ".toml"
	return filepath.Join(filePath, fileNameWithExtension)
}

func GetFlipStorageInitializersPath() string {
	return filepath.Join(GetProjectPath(), "transformers", "storage", "flip", "initializers")
}

func GetMedianStorageInitializersPath() string {
	return filepath.Join(GetProjectPath(), "transformers", "storage", "median", "initializers")
}
