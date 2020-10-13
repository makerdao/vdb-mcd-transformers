package initializer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type IGenerate interface {
	GenerateFlipInitializer() error
	GenerateMedianInitializer() error
}

type Generator struct {
	ProjectPath               string
	Collateral                types.Collateral
	MedianInitializerRequired bool
}

func NewInitializerGenerator(projectPath string, collateral types.Collateral, medianInitializerRequired bool) Generator {
	return Generator{
		ProjectPath:               projectPath,
		Collateral:                collateral,
		MedianInitializerRequired: medianInitializerRequired,
	}
}

func (g *Generator) GenerateFlipInitializer() error {
	initializer := g.createInitializer(g.Collateral.FormattedVersion(), g.Collateral.GetFlipContractName(), "flip")
	//create the path to the initializer file
	path := g.GetAbsoluteFlipStorageInitializersDirectoryPath()
	mkDirErr := os.MkdirAll(path, os.ModePerm)
	if mkDirErr != nil {
		return mkDirErr
	}

	writeFileErr := initializer.Save(g.GetAbsoluteFlipStorageInitializerFilePath())
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func (g Generator) GenerateMedianInitializer() error {
	if g.MedianInitializerRequired {
		initializer := g.createInitializer(g.Collateral.GetMedianInitializerDirectory(), g.Collateral.GetMedianContractName(), "median")

		path := g.GetAbsoluteMedianStorageInitializersDirectoryPath()
		mkDirErr := os.MkdirAll(path, os.ModePerm)
		if mkDirErr != nil {
			return mkDirErr
		}

		writeFileErr := initializer.Save(g.GetAbsoluteMedianStorageInitializerFilePath())
		if writeFileErr != nil {
			return writeFileErr
		}
	}
	return nil
}

func (g Generator) createInitializer(packageName, contractName, initializerType string) *jen.File {
	initializer := jen.NewFile(packageName)
	initializer.HeaderComment("This is a plugin generated to export the configured transformer initializers")

	initializer.Var().Id("contractAddress").Op("=").Qual(
		"github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants",
		"GetContractAddress").Params(jen.Lit(contractName))
	initializer.Var().Id("StorageTransformerInitializer").Op("=").Qual(
		fmt.Sprintf("github.com/makerdao/vdb-mcd-transformers/transformers/storage/%s/initializers", initializerType),
		"GenerateStorageTransformerInitializer").Params(jen.Id("contractAddress"))

	return initializer
}

func (g Generator) GetAbsoluteFlipStorageInitializersDirectoryPath() string {
	// example: $GOPATH/transformer/storage/flip/initializers/eth_b/v1_1_3
	return filepath.Join(
		g.ProjectPath, "transformers", "storage", "flip", "initializers", g.Collateral.GetFlipInitializerDirectory())
}

func (g Generator) GetAbsoluteMedianStorageInitializersDirectoryPath() string {
	// example: $GOPATH/transformer/storage/median/initializers/median_eth_b
	return filepath.Join(
		g.ProjectPath, "transformers", "storage", "median", "initializers", g.Collateral.GetMedianInitializerDirectory())
}

func (g Generator) GetAbsoluteFlipStorageInitializerFilePath() string {
	// example: $GOPATH/transformer/storage/flip/initializers/eth_b/v1_1_3/initializer.go
	return filepath.Join(
		g.ProjectPath, "transformers", "storage", "flip", "initializers", g.Collateral.GetFlipInitializerDirectory(), "initializer.go")
}

func (g Generator) GetAbsoluteMedianStorageInitializerFilePath() string {
	// example: $GOPATH/transformer/storage/median/initializers/median_eth_b/initializer.go
	return filepath.Join(
		g.ProjectPath, "transformers", "storage", "median", "initializers", g.Collateral.GetMedianInitializerDirectory(), "initializer.go")
}
