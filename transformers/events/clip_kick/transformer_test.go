package clip_kick_test

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/events/clip_kick"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ClipKick transformer", func() {
	var (
		db          = test_config.NewTestDB(test_config.NewTestNode())
		transformer = clip_kick.Transformer{}
	)
	It("returns an error if converting the log to an entity fails", func() {
		_, err := transformer.ToModels("wrong abi", []core.EventLog{test_data.DealEventLog}, db)

		Expect(err).To(HaveOccurred())
	})
})
