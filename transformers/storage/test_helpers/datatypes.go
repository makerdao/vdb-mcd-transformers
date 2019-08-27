package test_helpers

import (
	. "github.com/onsi/gomega"
)

type BlockMetadata struct {
	BlockNumber int    `db:"block_number"`
	BlockHash   string `db:"block_hash"`
}

type VariableRes struct {
	BlockMetadata
	Value string
}

type AuctionVariableRes struct {
	VariableRes
	ContractAddress string `db:"contract_address"`
}

type MappingRes struct {
	BlockMetadata
	Key   string
	Value string
}

type AuctionMappingRes struct {
	MappingRes
	ContractAddress string `db:"contract_address"`
}

type DoubleMappingRes struct {
	BlockMetadata
	KeyOne string `db:"key_one"`
	KeyTwo string `db:"key_two"`
	Value  string
}

type FlapRes struct {
	BlockMetadata
	ContractAddress string `db:"contract_address"`
	Id              string
	BidId           string `db:"bid_id"`
	Guy             string
	Tic             string
	End             string
	Lot             string
	Bid             string
	Gal             string
}

func AssertVariable(res VariableRes, blockNumber int, blockHash, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Value).To(Equal(value))
}

func AssertMapping(res MappingRes, blockNumber int, blockHash, key, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

func AssertDoubleMapping(res DoubleMappingRes, blockNumber int, blockHash, keyOne, keyTwo, value string) {
	Expect(res.BlockNumber).To(Equal(blockNumber))
	Expect(res.BlockHash).To(Equal(blockHash))
	Expect(res.KeyOne).To(Equal(keyOne))
	Expect(res.KeyTwo).To(Equal(keyTwo))
	Expect(res.Value).To(Equal(value))
}
