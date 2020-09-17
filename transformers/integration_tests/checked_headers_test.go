package integration_tests

import (
	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Checked Headers model", func() {
	var db = test_config.NewTestDB(test_config.NewTestNode())

	BeforeEach(func() {
		test_config.CleanTestDB(db)
	})

	It("marks headers as checked in this schema", func() {
		// insert header
		fakeHeader := fakes.GetFakeHeader(1)
		headerRepo := repositories.NewHeaderRepository(db)
		headerID, headerErr := headerRepo.CreateOrUpdateHeader(fakeHeader)
		Expect(headerErr).NotTo(HaveOccurred())

		checkedHeaderRepo := repositories.NewCheckedHeadersRepository(db, "maker")
		uncheckedHeaders, uncheckedHeaderErr := checkedHeaderRepo.UncheckedHeaders(0, -1, 1)
		Expect(uncheckedHeaderErr).NotTo(HaveOccurred())
		Expect(uncheckedHeaders).To(ContainElement(fakeHeader))

		markHeaderErr := checkedHeaderRepo.MarkHeaderChecked(headerID)
		Expect(markHeaderErr).NotTo(HaveOccurred())

		noUncheckedHeaders, noUncheckedHeadersErr := checkedHeaderRepo.UncheckedHeaders(0, -1, 1)
		Expect(noUncheckedHeadersErr).NotTo(HaveOccurred())
		Expect(noUncheckedHeaders).To(BeEmpty())
	})
})

func getBlockNumbers(headers []core.Header) []int64 {
	var headerBlockNumbers []int64
	for _, header := range headers {
		headerBlockNumbers = append(headerBlockNumbers, header.BlockNumber)
	}
	return headerBlockNumbers
}
