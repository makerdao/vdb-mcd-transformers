package queries

import (
	"fmt"
	. "github.com/onsi/ginkgo"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/component_tests/queries/test_helpers"
)

var _ = Describe("data generator", func() {
	It("writes reasonable data", func() {
		db := test_config.NewTestDB(test_config.NewTestNode())
		test_config.CleanTestDB(db)
		state := test_helpers.NewGenerator(db)
		state.Run(10)
		fmt.Println("Stop!")
	})
})
