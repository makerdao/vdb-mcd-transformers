package initializer_test

import (
	"os"
	"path/filepath"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/helpers"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GeneratorFlipInitializer", func() {
	var (
		collateral             = types.NewCollateral("TEST-COLLATERAL", "0.1.2")
		generator              = initializer.Generator{Collateral: collateral}
		testFlipCollateralPath = filepath.Join(
			helpers.GetFlipStorageInitializersPath(), "test_collateral",
		)
		testFlipCollateralFullPath = filepath.Join(
			testFlipCollateralPath, "v0_1_2", "initializer.go",
		)
		testMedianCollateralPath = filepath.Join(
			helpers.GetMedianStorageInitializersPath(), "median_test_collateral",
		)
		testMedianCollateralFullPath = filepath.Join(
			testMedianCollateralPath, "initializer.go",
		)
	)
	It("creates a flip initializer for new collateral", func() {
		initializerErr := generator.GenerateFlipInitializer()
		Expect(initializerErr).NotTo(HaveOccurred())

		fileInfo, fileErr := os.Stat(testFlipCollateralFullPath)
		Expect(os.IsNotExist(fileErr)).To(BeFalse())
		Expect(fileInfo.IsDir()).To(BeFalse())

		removeTestFile := os.RemoveAll(testFlipCollateralPath)
		Expect(removeTestFile).NotTo(HaveOccurred())
	})

	It("creates a median initializer for new collateral", func() {
		initializerErr := generator.GenerateMedianInitializer()
		Expect(initializerErr).NotTo(HaveOccurred())

		fileInfo, fileErr := os.Stat(testMedianCollateralFullPath)
		Expect(os.IsNotExist(fileErr)).To(BeFalse())
		Expect(fileInfo.IsDir()).To(BeFalse())

		removeTestFile := os.RemoveAll(testMedianCollateralPath)
		Expect(removeTestFile).NotTo(HaveOccurred())
	})
})
