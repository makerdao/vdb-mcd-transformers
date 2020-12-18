package initializer_test

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/initializer"
	"github.com/makerdao/vdb-mcd-transformers/generators/new_collateral/types"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("InitializerGenerator", func() {
	var (
		collateral                 = types.NewCollateral("TEST-COLLATERAL", "0.1.2")
		testProjectPath            = "./tmp"
		generator                  = initializer.NewInitializerGenerator(testProjectPath, collateral, true)
		testFlipCollateralPath     = generator.GetAbsoluteFlipStorageInitializersDirectoryPath()
		testFlipCollateralFullPath = filepath.Join(
			testFlipCollateralPath, "initializer.go",
		)
		testMedianCollateralPath     = generator.GetAbsoluteMedianStorageInitializersDirectoryPath()
		testMedianCollateralFullPath = filepath.Join(
			testMedianCollateralPath, "initializer.go",
		)
	)

	Context("get absolute initializer directory path", func() {
		It("formats the collateral as a flip initializer path", func() {
			flipInitializersPath := filepath.Join(testProjectPath, "transformers", "storage", "flip", "initializers")
			Expect(generator.GetAbsoluteFlipStorageInitializersDirectoryPath()).To(Equal(filepath.Join(flipInitializersPath, "test_collateral/v0_1_2")))
		})

		It("formats the collateral as a median initializer path", func() {
			medianInitializersPath := filepath.Join(testProjectPath, "transformers", "storage", "median", "initializers")
			Expect(generator.GetAbsoluteMedianStorageInitializersDirectoryPath()).To(Equal(filepath.Join(medianInitializersPath, "median_test_collateral/v0_1_2")))
		})
	})

	Context("get absolute initializer file path", func() {
		It("formats the collateral as a flip initializer file path", func() {
			flipInitializersPath := filepath.Join(testProjectPath, "transformers", "storage", "flip", "initializers")
			Expect(generator.GetAbsoluteFlipStorageInitializerFilePath()).To(Equal(filepath.Join(flipInitializersPath, "test_collateral/v0_1_2", "initializer.go")))
		})

		It("formats the collateral as a median initializer file path", func() {
			medianInitializersPath := filepath.Join(testProjectPath, "transformers", "storage", "median", "initializers")
			Expect(generator.GetAbsoluteMedianStorageInitializerFilePath()).To(Equal(filepath.Join(medianInitializersPath, "median_test_collateral/v0_1_2", "initializer.go")))
		})
	})

	Context("GenerateFlipInitializer", func() {
		It("creates a flip initializer for new collateral", func() {
			initializerErr := generator.GenerateFlipInitializer()
			Expect(initializerErr).NotTo(HaveOccurred())

			fileInfo, fileErr := os.Stat(testFlipCollateralFullPath)
			Expect(os.IsNotExist(fileErr)).To(BeFalse())
			Expect(fileInfo.IsDir()).To(BeFalse())

			fileContents, readFileErr := ioutil.ReadFile(testFlipCollateralFullPath)
			Expect(readFileErr).NotTo(HaveOccurred())
			Expect(string(fileContents)).To(Equal(expectedFlipInitializerFileContents))

		})
	})

	Context("GenerateMedianInitializer", func() {
		It("creates a median initializer for new collateral", func() {
			initializerErr := generator.GenerateMedianInitializer()
			Expect(initializerErr).NotTo(HaveOccurred())

			fileInfo, fileErr := os.Stat(testMedianCollateralFullPath)
			Expect(os.IsNotExist(fileErr)).To(BeFalse())
			Expect(fileInfo.IsDir()).To(BeFalse())

			fileContents, readFileErr := ioutil.ReadFile(testMedianCollateralFullPath)
			Expect(readFileErr).NotTo(HaveOccurred())
			Expect(string(fileContents)).To(Equal(expectedMedianInitializerFileContents))

			removeTestFile := os.RemoveAll(testMedianCollateralPath)
			Expect(removeTestFile).NotTo(HaveOccurred())
		})

		It("doesn't create a median initializer if it is not configured to do so", func() {
			generator = initializer.Generator{
				Collateral:                collateral,
				MedianInitializerRequired: false,
			}
			initializerErr := generator.GenerateMedianInitializer()
			Expect(initializerErr).NotTo(HaveOccurred())

			_, fileErr := os.Stat(testMedianCollateralFullPath)
			Expect(os.IsNotExist(fileErr)).To(BeTrue())
		})
	})

	AfterSuite(func() {
		removeTestProjectDirectoryErr := os.RemoveAll(testProjectPath)
		Expect(removeTestProjectDirectoryErr).NotTo(HaveOccurred())
	})
})

var expectedFlipInitializerFileContents = `// This is a plugin generated to export the configured transformer initializers

package v0_1_2

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/flip/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MCD_FLIP_TEST_COLLATERAL_0_1_2")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
`

var expectedMedianInitializerFileContents = `// This is a plugin generated to export the configured transformer initializers

package v0_1_2

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MEDIAN_TEST_COLLATERAL_0_1_2")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
`
