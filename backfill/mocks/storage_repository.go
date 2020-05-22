package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

type StorageRepository struct {
	GetOrCreateUrnError           error
	GetOrCreateUrnIDsToReturn     map[string]int64
	GetUrnByIDUrnToReturn         repository.Urn
	GetUrnByIDError               error
	InsertDiffErr                 error
	InsertDiffPassedDiff          types.RawDiff
	VatIlkArtExistsBoolToReturn   bool
	VatIlkArtExistsPassedHeaderID int64
	VatIlkArtExistsPassedIlkID    int64
	VatUrnArtExistsBoolToReturn   bool
	VatUrnArtExistsPassedHeaderID int64
	VatUrnArtExistsPassedUrnID    int64
	VatUrnInkExistsBoolToReturn   bool
	VatUrnInkExistsPassedHeaderID int64
	VatUrnInkExistsPassedUrnID    int64
}

func (mock *StorageRepository) GetOrCreateUrn(urn, ilk string) (int64, error) {
	if id, ok := mock.GetOrCreateUrnIDsToReturn[urn]; ok {
		return id, nil
	}
	return 0, mock.GetOrCreateUrnError
}

func (mock *StorageRepository) GetUrnByID(id int64) (repository.Urn, error) {
	return mock.GetUrnByIDUrnToReturn, mock.GetUrnByIDError
}

func (mock *StorageRepository) InsertDiff(diff types.RawDiff) error {
	mock.InsertDiffPassedDiff = diff
	return mock.InsertDiffErr
}

func (mock *StorageRepository) VatIlkArtExists(ilkID, headerID int64) (bool, error) {
	mock.VatIlkArtExistsPassedHeaderID = headerID
	mock.VatIlkArtExistsPassedIlkID = ilkID
	return mock.VatIlkArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnArtExists(urnID, headerID int64) (bool, error) {
	mock.VatUrnArtExistsPassedHeaderID = headerID
	mock.VatUrnArtExistsPassedUrnID = urnID
	return mock.VatUrnArtExistsBoolToReturn, nil
}

func (mock *StorageRepository) VatUrnInkExists(urnID, headerID int64) (bool, error) {
	mock.VatUrnInkExistsPassedHeaderID = headerID
	mock.VatUrnInkExistsPassedUrnID = urnID
	return mock.VatUrnInkExistsBoolToReturn, nil
}
