package bid_creation_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/triggers/bid_creation"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("flap created trigger", func() {
	var flapKickModel event.InsertionModel

	BeforeEach(func() {
		flapKickModel = test_data.FlapKickModel()
	})

	Describe("updating flap created", func() {
		bid_creation.SharedBidCreationTriggerTests(constants.FlapTable, test_data.FlapAddress(), &flapKickModel)
	})
})
