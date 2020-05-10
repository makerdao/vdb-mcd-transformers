package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

type StorageRepository struct {
	GetUrnByIDUrnToReturn         backfill.Urn
	GetUrnByIDError               error
	GetUrnsCalled                 bool
	GetUrnsUrnsToReturn           []backfill.Urn
	GetUrnsErr                    error
	InsertDiffErr                 error
	InsertDiffPassedDiff          types.RawDiff
	VatIlkArtExistsBoolToReturn   bool
	VatIlkArtExistsPassedHeaderID int
	VatIlkArtExistsPassedIlkID    int
	VatUrnArtExistsBoolToReturn   bool
	VatUrnArtExistsPassedHeaderID int
	VatUrnArtExistsPassedUrnID    int
	VatUrnInkExistsBoolToReturn   bool
	VatUrnInkExistsPassedHeaderID int
	VatUrnInkExistsPassedUrnID    int
}

func (mock *StorageRepository) GetUrnByID(id int) (backfill.Urn, error) {
	return mock.GetUrnByIDUrnToReturn, mock.GetUrnByIDError
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
	mock.VatIlkArtExistsPassedHeaderID = headerID
	mock.VatIlkArtExistsPassedIlkID = ilkID
	return mock.VatIlkArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	mock.VatUrnArtExistsPassedHeaderID = headerID
	mock.VatUrnArtExistsPassedUrnID = urnID
	return mock.VatUrnArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	mock.VatUrnInkExistsPassedHeaderID = headerID
	mock.VatUrnInkExistsPassedUrnID = urnID
	return mock.VatUrnInkExistsBoolToReturn, nil
}
