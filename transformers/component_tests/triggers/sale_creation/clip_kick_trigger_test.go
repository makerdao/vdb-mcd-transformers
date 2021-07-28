package sale_creation_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/triggers/sale_creation"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("clip created trigger", func() {
	var clipKickModel event.InsertionModel

	BeforeEach(func() {
		clipKickModel = test_data.ClipKickModel()
	})

	Describe("updating clip created", func() {
		sale_creation.SharedSaleCreationTriggerTests(constants.ClipTable, test_data.ClipEthAV150Address(), &clipKickModel)
	})
})
