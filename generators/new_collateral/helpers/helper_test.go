package helpers_test

import (
	"os"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("helpers", func() {
	var goPath = os.ExpandEnv("$GOPATH")
	It("returns the project path", func() {
		Expect(helpers.GetProjectPath()).To(
			Equal(goPath + "/src/github.com/makerdao/vdb-mcd-transformers"))
	})

	It("returns the environment path", func() {
		Expect(helpers.GetEnvironmentsPath()).To(
			Equal(goPath + "/src/github.com/makerdao/vdb-mcd-transformers/environments"))
	})

	It("returns the execute plugins path", func() {
		Expect(helpers.GetExecutePluginsPath()).To(
			Equal(goPath + "/src/github.com/makerdao/vdb-mcd-transformers/plugins/execute"))
	})

	It("returns the full config file path and file name", func() {
		Expect(helpers.GetFullConfigFilePath("path", "file")).To(
			Equal("path/file.toml"))
	})
})
