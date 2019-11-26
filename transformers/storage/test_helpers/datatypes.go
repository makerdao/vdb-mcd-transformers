package test_helpers

import (
	. "github.com/onsi/gomega"
)

type BlockMetadata struct {
	HeaderID int64 `db:"header_id"`
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

type DoubleMappingRes struct {
	BlockMetadata
	KeyOne string `db:"key_one"`
	KeyTwo string `db:"key_two"`
	Value  string
}

type FlapRes struct {
	BlockMetadata
	BlockNumber     int64  `db:"block_number"`
	ContractAddress string `db:"contract_address"`
	Id              string
	BidId           string `db:"bid_id"`
	Guy             string
	Tic             string
	End             string
	Lot             string
	Bid             string
}

type FlopRes struct {
	BlockMetadata
	BlockNumber     int64  `db:"block_number"`
	ContractAddress string `db:"contract_address"`
	Id              string
	BidId           string `db:"bid_id"`
	Guy             string
	Tic             string
	End             string
	Lot             string
	Bid             string
}

func AssertVariable(res VariableRes, headerID int64, value string) {
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.Value).To(Equal(value))
}

func AssertMapping(res MappingRes, headerID int64, key, value string) {
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

func AssertDoubleMapping(res DoubleMappingRes, headerID int64, keyOne, keyTwo, value string) {
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.KeyOne).To(Equal(keyOne))
	Expect(res.KeyTwo).To(Equal(keyTwo))
	Expect(res.Value).To(Equal(value))
}
