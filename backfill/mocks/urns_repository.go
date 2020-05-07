package mocks

import (
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage/types"
)

type UrnsRepository struct {
	GetUrnsCalled               bool
	GetUrnsUrnsToReturn         []backfill.Urn
	GetUrnsErr                  error
	InsertUrnDiffErr            error
	InsertUrnDiffPassedDiff     types.RawDiff
	VatUrnArtExistsBoolToReturn bool
	VatUrnArtExistsErr          error
	VatUrnInkExistsBoolToReturn bool
	VatUrnInkExistsErr          error
}

func (u *UrnsRepository) GetUrns() ([]backfill.Urn, error) {
	u.GetUrnsCalled = true
	return u.GetUrnsUrnsToReturn, u.GetUrnsErr
}

func (u *UrnsRepository) InsertUrnDiff(diff types.RawDiff) error {
	u.InsertUrnDiffPassedDiff = diff
	return u.InsertUrnDiffErr
}

func (u *UrnsRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	return u.VatUrnArtExistsBoolToReturn, u.VatUrnArtExistsErr
}

func (u *UrnsRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	return u.VatUrnInkExistsBoolToReturn, u.VatUrnInkExistsErr
}
