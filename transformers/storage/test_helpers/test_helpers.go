package test_helpers

import (
	"time"

	. "github.com/onsi/gomega"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
)

func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(time.RFC3339)
}

func CreateHeader(timestamp int64, blockNumber int, db *postgres.DB) {
	fakeHeader := fakes.GetFakeHeaderWithTimestamp(timestamp, int64(blockNumber))
	headerRepo := repositories.NewHeaderRepository(db)
	_, headerErr := headerRepo.CreateOrUpdateHeader(fakeHeader)
	Expect(headerErr).NotTo(HaveOccurred())
}
