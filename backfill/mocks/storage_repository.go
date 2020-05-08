package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

type StorageRepository struct {
	GetUrnsCalled               bool
	GetUrnsUrnsToReturn         []backfill.Urn
	GetUrnsErr                  error
	InsertDiffErr               error
	InsertDiffPassedDiff        types.RawDiff
	VatIlkArtExistsBoolToReturn bool
	VatUrnArtExistsBoolToReturn bool
	VatUrnInkExistsBoolToReturn bool
}

func (mock *StorageRepository) GetUrns() ([]backfill.Urn, error) {
	mock.GetUrnsCalled = true
	return mock.GetUrnsUrnsToReturn, mock.GetUrnsErr
}

func (mock *StorageRepository) InsertDiff(diff types.RawDiff) error {
	mock.InsertDiffPassedDiff = diff
	return mock.InsertDiffErr
}

func (mock *StorageRepository) VatIlkArtExists(ilkID, headerID int) (bool, error) {
	return mock.VatIlkArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	return mock.VatUrnArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	return mock.VatUrnInkExistsBoolToReturn, nil
}
