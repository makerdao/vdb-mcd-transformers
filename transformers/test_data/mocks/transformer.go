package mocks

import (
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/libraries/shared/transformer"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

type MockEventTransformer struct {
	ExecuteWasCalled bool
	ExecuteError     error
	PassedLogs       []types.Log
	PassedHeader     core.Header
	config           transformer.EventTransformerConfig
}

func (mh *MockEventTransformer) Execute(logs []types.Log, header core.Header) error {
	if mh.ExecuteError != nil {
		return mh.ExecuteError
	}
	mh.ExecuteWasCalled = true
	mh.PassedLogs = logs
	mh.PassedHeader = header
	return nil
}

func (mh *MockEventTransformer) GetConfig() transformer.EventTransformerConfig {
	return mh.config
}

func (mh *MockEventTransformer) SetEventTransformerConfig(config transformer.EventTransformerConfig) {
	mh.config = config
}

func (mh *MockEventTransformer) FakeEventTransformerInitializer(db *postgres.DB) transformer.EventTransformer {
	return mh
}

var FakeTransformerConfig = transformer.EventTransformerConfig{
	TransformerName:   "FakeTransformer",
	ContractAddresses: []string{"FakeAddress"},
	Topic:             "FakeTopic",
}
