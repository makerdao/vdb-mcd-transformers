package mocks

import "github.com/makerdao/vdb-mcd-transformers/backfill"

type UrnsRepository struct {
	GetUrnsCalled               bool
	GetUrnsUrnsToReturn         []backfill.Urn
	GetUrnsErr                  error
	VatUrnArtExistsBoolToReturn bool
	VatUrnArtExistsErr          error
	VatUrnInkExistsBoolToReturn bool
	VatUrnInkExistsErr          error
}

func (u *UrnsRepository) GetUrns() ([]backfill.Urn, error) {
	u.GetUrnsCalled = true
	return u.GetUrnsUrnsToReturn, u.GetUrnsErr
}

func (u *UrnsRepository) VatUrnArtExists(urnID, headerID int) (bool, error) {
	return u.VatUrnArtExistsBoolToReturn, u.VatUrnArtExistsErr
}

func (u *UrnsRepository) VatUrnInkExists(urnID, headerID int) (bool, error) {
	return u.VatUrnInkExistsBoolToReturn, u.VatUrnInkExistsErr
}
