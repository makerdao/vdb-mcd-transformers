package integration_tests

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/libraries/shared/test_data"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("maker", func() {
	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("has a proper checked headers setup in the maker schema", func() {
		test_data.ExpectCheckedHeadersInThisSchema(db, "maker")
	})
})
