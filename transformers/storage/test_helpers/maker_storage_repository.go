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
	GetCdpisCalled          bool
	GetCdpisError           error
	GetDaiKeysCalled        bool
	GetDaiKeysError         error
	GetGemKeysCalled        bool
	GetGemKeysError         error
	GetFlapBidIdsCalled     bool
	GetFlapBidIdsError      error
	GetFlipBidIdsCalledWith string
	GetFlipBidIdsError      error
	GetFlopBidIdsCalledWith string
	GetFlopBidIdsError      error
	GetIlksCalled           bool
	GetIlksError            error
	GetOwnersCalled         bool
	GetOwnersError          error
	GetPotPieUsersCalled    bool
	GetPotPieUsersError     error
	GetVatSinKeysCalled     bool
	GetVatSinKeysError      error
	GetVowSinKeysCalled     bool
	GetVowSinKeysError      error
	GetUrnsCalled           bool
	GetUrnsError            error
}

func (repository *MockMakerStorageRepository) GetFlapBidIds(string) ([]string, error) {
	repository.GetFlapBidIdsCalled = true
	return repository.FlapBidIds, repository.GetFlapBidIdsError
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]storage.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
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

func (repository *MockMakerStorageRepository) GetFlipBidIds(contractAddress string) ([]string, error) {
	repository.GetFlipBidIdsCalledWith = contractAddress
	return repository.FlipBidIds, repository.GetFlipBidIdsError
}

func (repository *MockMakerStorageRepository) GetFlopBidIds(contractAddress string) ([]string, error) {
	repository.GetFlopBidIdsCalledWith = contractAddress
	return repository.FlopBidIds, repository.GetFlopBidIdsError
}

func (repository *MockMakerStorageRepository) GetCdpis() ([]string, error) {
	repository.GetCdpisCalled = true
	return repository.Cdpis, repository.GetCdpisError
}

func (repository *MockMakerStorageRepository) GetOwners() ([]string, error) {
	repository.GetOwnersCalled = true
	return repository.Owners, repository.GetOwnersError
}

func (repository *MockMakerStorageRepository) GetPotPieUsers() ([]string, error) {
	repository.GetPotPieUsersCalled = true
	return repository.PotPieUsers, repository.GetPotPieUsersError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}
