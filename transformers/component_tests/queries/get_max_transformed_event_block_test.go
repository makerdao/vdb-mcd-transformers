package queries

import (
	"math/rand"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Get max transformed event block query", func() {
	var (
		headerOne, headerTwo core.Header
		logOne, logTwo       core.EventLog
	)

	BeforeEach(func() {
		test_config.CleanTestDB(db)
		headerRepo := repositories.NewHeaderRepository(db)

		blockOne := rand.Int()
		timestampOne := int(rand.Int31())
		headerOne = createHeader(blockOne, timestampOne, headerRepo)
		headerTwo = createHeader(blockOne+1, timestampOne+1, headerRepo)
		logOne = test_data.CreateTestLog(headerOne.Id, db)
		logTwo = test_data.CreateTestLog(headerTwo.Id, db)
	})

	It("returns most recent block if log from that block is transformed", func() {
		_, err := db.Exec(`UPDATE public.event_logs SET transformed = true WHERE id = $1`, logOne.ID)
		Expect(err).NotTo(HaveOccurred())
		_, err = db.Exec(`UPDATE public.event_logs SET transformed = true WHERE id = $1`, logTwo.ID)
		Expect(err).NotTo(HaveOccurred())

		var result int64
		err = db.Get(&result, `SELECT * FROM api.get_max_transformed_event_block()`)
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(headerTwo.BlockNumber))
	})

	It("returns earlier block if log from most recent block is not transformed", func() {
		_, err := db.Exec(`UPDATE public.event_logs SET transformed = true WHERE id = $1`, logOne.ID)
		Expect(err).NotTo(HaveOccurred())

		var result int64
		err = db.Get(&result, `SELECT * FROM api.get_max_transformed_event_block()`)
		Expect(err).NotTo(HaveOccurred())

		Expect(result).To(Equal(headerOne.BlockNumber))
	})
})
