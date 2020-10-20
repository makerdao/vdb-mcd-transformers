package new_collateral_generator

import (
	"os"
	"path/filepath"
)

func getProjectPath() string {
	projectPath := filepath.Join("$GOPATH", "src", "github.com", "makerdao", "vdb-mcd-transformers")
	return os.ExpandEnv(projectPath)
}

func GetEnvironmentsPath() string {
	return filepath.Join(getProjectPath(), "environments")
}

func GetExecutePluginsPath() string {
	return filepath.Join(getProjectPath(), "plugins", "execute")
}

func GetFullConfigFilePath(filePath, fileName string) string {
	fileNameWithExtension := fileName + ".toml"
	return filepath.Join(filePath, fileNameWithExtension)
}
