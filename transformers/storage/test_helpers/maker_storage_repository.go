package test_helpers

import (
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"

	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
)

type MockMakerStorageRepository struct {
	Cdpis                   []string
	DaiKeys                 []string
	FlapBidIds              []string
	FlipBidIds              []string
	FlopBidIds              []string
	GemKeys                 []storage.Urn
	Ilks                    []string
	Owners                  []string
	PotPieUsers             []string
	SinKeys                 []string
	Urns                    []storage.Urn
	VatCanKeys              []storage.Can
	VatWardsKeys            []string
	WardsKeys               []string
	GetCdpisCalled          bool
	GetCdpisError           error
	GetDaiKeysCalled        bool
	GetDaiKeysError         error
	GetFlapBidIdsCalled     bool
	GetFlapBidIdsError      error
	GetFlipBidIdsCalledWith string
	GetFlipBidIdsError      error
	GetFlopBidIdsCalledWith string
	GetFlopBidIdsError      error
	GetGemKeysCalled        bool
	GetGemKeysError         error
	GetIlksCalled           bool
	GetIlksError            error
	GetOwnersCalled         bool
	GetOwnersError          error
	GetPotPieUsersCalled    bool
	GetPotPieUsersError     error
	GetUrnsCalled           bool
	GetUrnsError            error
	GetVatCanKeysCalled     bool
	GetVatCanKeysError      error
	GetVatSinKeysCalled     bool
	GetVatSinKeysError      error
	GetVatWardsKeysCalled   bool
	GetVatWardsKeysError    error
	GetVowSinKeysCalled     bool
	GetVowSinKeysError      error
	GetWardsKeysCalledWith  string
	GetWardsKeysError       error
}

func (repository *MockMakerStorageRepository) GetCdpis() ([]string, error) {
	repository.GetCdpisCalled = true
	return repository.Cdpis, repository.GetCdpisError
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetFlapBidIds(string) ([]string, error) {
	repository.GetFlapBidIdsCalled = true
	return repository.FlapBidIds, repository.GetFlapBidIdsError
}

func (repository *MockMakerStorageRepository) GetFlipBidIds(contractAddress string) ([]string, error) {
	repository.GetFlipBidIdsCalledWith = contractAddress
	return repository.FlipBidIds, repository.GetFlipBidIdsError
}

func (repository *MockMakerStorageRepository) GetFlopBidIds(contractAddress string) ([]string, error) {
	repository.GetFlopBidIdsCalledWith = contractAddress
	return repository.FlopBidIds, repository.GetFlopBidIdsError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]storage.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
}

func (repository *MockMakerStorageRepository) GetOwners() ([]string, error) {
	repository.GetOwnersCalled = true
	return repository.Owners, repository.GetOwnersError
}

func (repository *MockMakerStorageRepository) GetPotPieUsers() ([]string, error) {
	repository.GetPotPieUsersCalled = true
	return repository.PotPieUsers, repository.GetPotPieUsersError
}

func (repository *MockMakerStorageRepository) GetVatSinKeys() ([]string, error) {
	repository.GetVatSinKeysCalled = true
	return repository.SinKeys, repository.GetVatSinKeysError
}

func (repository *MockMakerStorageRepository) GetVowSinKeys() ([]string, error) {
	repository.GetVowSinKeysCalled = true
	return repository.SinKeys, repository.GetVowSinKeysError
}

func (repository *MockMakerStorageRepository) GetUrns() ([]storage.Urn, error) {
	repository.GetUrnsCalled = true
	return repository.Urns, repository.GetUrnsError
}

func (repository *MockMakerStorageRepository) GetVatCanKeys() ([]storage.Can, error) {
	repository.GetVatCanKeysCalled = true
	return repository.VatCanKeys, repository.GetVatCanKeysError
}

func (repository *MockMakerStorageRepository) GetVatWardsAddresses() ([]string, error) {
	repository.GetVatWardsKeysCalled = true
	return repository.VatWardsKeys, repository.GetVatWardsKeysError
}

func (repository *MockMakerStorageRepository) GetWardsAddresses(contractAddress string) ([]string, error) {
	repository.GetWardsKeysCalledWith = contractAddress
	return repository.WardsKeys, repository.GetWardsKeysError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}
