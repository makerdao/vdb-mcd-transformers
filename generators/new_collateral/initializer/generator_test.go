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

var _ = Describe("GeneratorFlipInitializer", func() {
	var (
		collateral = types.NewCollateral("TEST-COLLATERAL", "0.1.2")
		generator  = initializer.Generator{
			Collateral:                collateral,
			MedianInitializerRequired: true,
		}
		testFlipCollateralPath     = collateral.GetAbsoluteFlipStorageInitializersDirectoryPath()
		testFlipCollateralFullPath = filepath.Join(
			testFlipCollateralPath, "initializer.go",
		)
		testMedianCollateralPath     = collateral.GetAbsoluteMedianStorageInitializersDirectoryPath()
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

		fileContents, readFileErr := ioutil.ReadFile(testFlipCollateralFullPath)
		Expect(readFileErr).NotTo(HaveOccurred())
		Expect(string(fileContents)).To(Equal(expectedFlipInitializerFileContents))

		removeTestFile := os.RemoveAll(testFlipCollateralPath)
		Expect(removeTestFile).NotTo(HaveOccurred())
	})

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

package median_test_collateral

import (
	initializers "github.com/makerdao/vdb-mcd-transformers/transformers/storage/median/initializers"
	constants "github.com/makerdao/vdb-transformer-utilities/pkg/shared/constants"
)

var contractAddress = constants.GetContractAddress("MEDIAN_TEST_COLLATERAL")
var StorageTransformerInitializer = initializers.GenerateStorageTransformerInitializer(contractAddress)
`
