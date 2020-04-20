package test_helpers

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	. "github.com/onsi/gomega"
)

func FormatTimestamp(timestamp int64) string {
	return time.Unix(timestamp, 0).UTC().Format(time.RFC3339)
}

func CreateHeader(timestamp int64, blockNumber int, db *postgres.DB) core.Header {
	return CreateHeaderWithHash(strconv.Itoa(rand.Int()), timestamp, blockNumber, db)
}

func CreateHeaderWithHash(hash string, timestamp int64, blockNumber int, db *postgres.DB) core.Header {
	fakeHeader := fakes.GetFakeHeaderWithTimestamp(timestamp, int64(blockNumber))
	fakeHeader.Hash = hash
	headerRepo := repositories.NewHeaderRepository(db)
	headerId, headerErr := headerRepo.CreateOrUpdateHeader(fakeHeader)
	Expect(headerErr).NotTo(HaveOccurred())
	fakeHeader.Id = headerId
	return fakeHeader
}

func CreateFakeDiffRecord(db *postgres.DB) int64 {
	return CreateFakeDiffRecordWithHeader(db, fakes.FakeHeader)
}

func CreateFakeDiffRecordWithHeader(db *postgres.DB, header core.Header) int64 {
	fakeRawDiff := GetFakeStorageDiffForHeader(header, common.Hash{}, common.Hash{}, common.Hash{})
	storageDiffRepo := storage.NewDiffRepository(db)
	diffID, insertDiffErr := storageDiffRepo.CreateStorageDiff(fakeRawDiff)
	Expect(insertDiffErr).NotTo(HaveOccurred())

	return diffID
}

func CreateDiffRecord(db *postgres.DB, header core.Header, hashedAddress, key, value common.Hash) types.PersistedDiff {
	rawDiff := GetFakeStorageDiffForHeader(header, hashedAddress, key, value)

	repo := storage.NewDiffRepository(db)
	diffID, insertDiffErr := repo.CreateStorageDiff(rawDiff)
	Expect(insertDiffErr).NotTo(HaveOccurred())

	persistedDiff := types.PersistedDiff{
		RawDiff:  rawDiff,
		ID:       diffID,
		HeaderID: header.Id,
	}

	return persistedDiff
}

func GetFakeStorageDiffForHeader(header core.Header, hashedAddress, storageKey, storageValue common.Hash) types.RawDiff {
	return types.RawDiff{
		HashedAddress: hashedAddress,
		BlockHash:     common.HexToHash(header.Hash),
		BlockHeight:   int(header.BlockNumber),
		StorageKey:    storageKey,
		StorageValue:  storageValue,
	}
}
