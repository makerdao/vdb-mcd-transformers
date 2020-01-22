package bid_creation_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/triggers/bid_creation"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("flip created trigger", func() {
	var flipKickModel event.InsertionModel

	BeforeEach(func() {
		flipKickModel = test_data.FlipKickModel()
	})

	Describe("updating flop created", func() {
		bid_creation.SharedBidCreationTriggerTests(constants.FlipTable, test_data.EthFlipAddress(), &flipKickModel)
	})
})
