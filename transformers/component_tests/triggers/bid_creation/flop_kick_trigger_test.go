package bid_creation_test

import (
	"github.com/makerdao/vdb-mcd-transformers/transformers/component_tests/triggers/bid_creation"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	. "github.com/onsi/ginkgo"
)

var _ = Describe("flop created trigger", func() {
	var flopKickModel event.InsertionModel

	BeforeEach(func() {
		flopKickModel = test_data.FlopKickModel()
	})

	Describe("updating flop created", func() {
		bid_creation.SharedBidCreationTriggerTests(constants.FlopTable, test_data.FlopV101Address(), &flopKickModel)
	})
})
