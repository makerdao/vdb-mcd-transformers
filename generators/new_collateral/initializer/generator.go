package initializer

import (
	"os"
	"path/filepath"

	"github.com/dave/jennifer/jen"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
)

type IGenerate interface {
	GenerateFlipInitializer() error
	GenerateMedianInitializer() error
}

type Generator struct {
	Collateral                types.Collateral
	MedianInitializerRequired bool
}

var initializerFileName = "initializer.go"

func (g *Generator) GenerateFlipInitializer() error {
	initializer := jen.NewFile(g.Collateral.FormattedVersion())
	initializer.HeaderComment("This is a plugin generated to export the configured transformer initializers")

	collateralContractName := g.Collateral.FormattedForFlipContractName()
	initializer.Var().Id("contractAddress").Op("=").Qual(
		"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants",
		"GetContractAddress").Params(jen.Lit(collateralContractName))
	initializer.Var().Id("StorageTransformerInitializer").Op("=").Qual(
		"github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers",
		"GenerateStorageTransformerInitializer").Params(jen.Id("contractAddress"))

	//create the path to the initializer file
	path := g.createFlipPath()
	mkDirErr := os.MkdirAll(path, os.ModePerm)
	if mkDirErr != nil {
		return mkDirErr
	}

	writeFileErr := initializer.Save(g.createFullFlipInitializerPath())
	if writeFileErr != nil {
		return writeFileErr
	}

	return nil
}

func (g *Generator) GenerateMedianInitializer() error {
	if g.MedianInitializerRequired {
		initializer := jen.NewFile(g.Collateral.FormattedVersion())
		initializer.HeaderComment("This is a plugin generated to export the configured transformer initializers")

		collateralContractName := g.Collateral.GetFlipContractName()
		initializer.Var().Id("contractAddress").Op("=").Qual(
			"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants",
			"GetContractAddress").Params(jen.Lit(collateralContractName))
		initializer.Var().Id("StorageTransformerInitializer").Op("=").Qual(
			"github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers",
			"GenerateStorageTransformerInitializer").Params(jen.Id("contractAddress"))

		path := g.Collateral.GetAbsoluteMedianStorageInitializersDirectoryPath()
		mkDirErr := os.MkdirAll(path, os.ModePerm)
		if mkDirErr != nil {
			return mkDirErr
		}

		writeFileErr := initializer.Save(g.Collateral.GetAbsoluteMedianStorageInitializerFilePath())
		if writeFileErr != nil {
			return writeFileErr
		}
	}
	return nil
}

func (g *Generator) createFlipPath() string {
	return filepath.Join(helpers.GetFlipStorageInitializersPath(),
		g.Collateral.FormattedForFlipInitializerFileName(),
	)
}

func (g *Generator) createFullFlipInitializerPath() string {
	return filepath.Join(g.createFlipPath(), initializerFileName)
}

func (g *Generator) createMedianPath() string {
	return filepath.Join(helpers.GetMedianStorageInitializersPath(),
		g.Collateral.FormattedForMedianInitializerFileName(),
	)
}

func (g *Generator) createFullMedianInitializerPath() string {
	return filepath.Join(g.createMedianPath(), initializerFileName)
}
