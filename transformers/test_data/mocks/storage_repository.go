package mocks

import (
	"github.com/vulcanize/vulcanizedb/libraries/shared/storage/utils"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockStorageRepository struct {
	CreateErr         error
	PassedBlockNumber int
	PassedBlockHash   string
	PassedMetadata    utils.StorageValueMetadata
	PassedValue       interface{}
}

func (repository *MockStorageRepository) Create(blockNumber int, blockHash string, metadata utils.StorageValueMetadata, value interface{}) error {
	repository.PassedBlockNumber = blockNumber
	repository.PassedBlockHash = blockHash
	repository.PassedMetadata = metadata
	repository.PassedValue = value
	return repository.CreateErr
}

func (*MockStorageRepository) SetDB(db *postgres.DB) {
	panic("implement me")
}
