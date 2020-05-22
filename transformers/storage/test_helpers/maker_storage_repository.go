package test_helpers

import (
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"

	"github.com/makerdao/vdb-mcd-transformers/transformers/storage"
)

type MockMakerStorageRepository struct {
	Cdpis                            []string
	DaiKeys                          []string
	FlapBidIds                       []string
	FlipBidIds                       []string
	FlopBidIds                       []string
	GemKeys                          []storage.Urn
	Ilks                             []string
	MedianBudAddresses               []string
	MedianOrclAddresses              []string
	MedianSlotIds                    []string
	Owners                           []string
	PotPieUsers                      []string
	SinKeys                          []string
	Urns                             []storage.Urn
	VatWardsKeys                     []string
	WardsKeys                        []string
	GetCdpisCalled                   bool
	GetCdpisError                    error
	GetDaiKeysCalled                 bool
	GetDaiKeysError                  error
	GetFlapBidIdsCalled              bool
	GetFlapBidIdsError               error
	GetFlipBidIdsCalledWith          string
	GetFlipBidIdsError               error
	GetFlopBidIdsCalledWith          string
	GetFlopBidIdsError               error
	GetGemKeysCalled                 bool
	GetGemKeysError                  error
	GetIlksCalled                    bool
	GetIlksError                     error
	GetMedianBudAddressesCalledWith  string
	GetMedianBudAddressesError       error
	GetMedianOrclAddressesCalledWith string
	GetMedianOrclAddressesError      error
	GetMedianSlotIdCalled            bool
	GetMedianSlotIdError             error
	GetOwnersCalled                  bool
	GetOwnersError                   error
	GetPotPieUsersCalled             bool
	GetPotPieUsersError              error
	GetVatSinKeysCalled              bool
	GetVatSinKeysError               error
	GetVowSinKeysCalled              bool
	GetVowSinKeysError               error
	GetUrnsCalled                    bool
	GetUrnsError                     error
	GetVatWardsKeysCalled            bool
	GetVatWardsKeysError             error
	GetWardsKeysCalledWith           string
	GetWardsKeysError                error
}

func (repository *MockMakerStorageRepository) GetCdpis() ([]string, error) {
	repository.GetCdpisCalled = true
	return repository.Cdpis, repository.GetCdpisError
}

func (repository *MockMakerStorageRepository) GetDaiKeys() ([]string, error) {
	repository.GetDaiKeysCalled = true
	return repository.DaiKeys, repository.GetDaiKeysError
}

func (repository *MockMakerStorageRepository) GetFlapBidIDs(string) ([]string, error) {
	repository.GetFlapBidIDsCalled = true
	return repository.FlapBidIDs, repository.GetFlapBidIDsError
}

func (repository *MockMakerStorageRepository) GetFlipBidIDs(contractAddress string) ([]string, error) {
	repository.GetFlipBidIDsCalledWith = contractAddress
	return repository.FlipBidIDs, repository.GetFlipBidIDsError
}

func (repository *MockMakerStorageRepository) GetFlopBidIDs(contractAddress string) ([]string, error) {
	repository.GetFlopBidIDsCalledWith = contractAddress
	return repository.FlopBidIDs, repository.GetFlopBidIDsError
}

func (repository *MockMakerStorageRepository) GetGemKeys() ([]storage.Urn, error) {
	repository.GetGemKeysCalled = true
	return repository.GemKeys, repository.GetGemKeysError
}

func (repository *MockMakerStorageRepository) GetIlks() ([]string, error) {
	repository.GetIlksCalled = true
	return repository.Ilks, repository.GetIlksError
}

func (repository *MockMakerStorageRepository) GetMedianBudAddresses(address string) ([]string, error) {
	repository.GetMedianBudAddressesCalledWith = address
	return repository.MedianBudAddresses, repository.GetMedianBudAddressesError
}

func (repository *MockMakerStorageRepository) GetMedianOrclAddresses(address string) ([]string, error) {
	repository.GetMedianOrclAddressesCalledWith = address
	return repository.MedianOrclAddresses, repository.GetMedianOrclAddressesError
}

func (repository *MockMakerStorageRepository) GetMedianSlotIds(string) ([]string, error) {
	repository.GetMedianSlotIdCalled = true
	return repository.MedianSlotIds, repository.GetMedianSlotIdError
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

func (repository *MockMakerStorageRepository) GetVatWardsAddresses() ([]string, error) {
	repository.GetVatWardsKeysCalled = true
	return repository.VatWardsKeys, repository.GetVatWardsKeysError
}

func (repository *MockMakerStorageRepository) GetWardsAddresses(contractAddress string) ([]string, error) {
	repository.GetWardsKeysCalledWith = contractAddress
	return repository.WardsKeys, repository.GetWardsKeysError
}

func (repository *MockMakerStorageRepository) SetDB(db *postgres.DB) {}
