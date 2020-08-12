package test_helpers

import (
	. "github.com/onsi/gomega"
)

type BlockMetadata struct {
	HeaderID int64 `db:"header_id"`
}

type DiffMetadata struct {
	BlockMetadata
	DiffID int64 `db:"diff_id"`
	Value  string
}

type VariableRes struct {
	DiffMetadata
}

type VariableResWithAddress struct {
	VariableRes
	AddressID int64 `db:"address_id"`
}

type MappingRes struct {
	DiffMetadata
	Key string
}

type DoubleMappingRes struct {
	DiffMetadata
	KeyOne string `db:"key_one"`
	KeyTwo string `db:"key_two"`
}

type MappingResWithAddress struct {
	MappingRes
	AddressID string `db:"address_id"`
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

func AssertVariable(res VariableRes, diffID, headerID int64, value string) {
	Expect(res.DiffID).To(Equal(diffID))
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.Value).To(Equal(value))
}

func AssertVariableWithAddress(res VariableResWithAddress, diffID, headerID, addressID int64, value string) {
	Expect(res.DiffID).To(Equal(diffID))
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.AddressID).To(Equal(addressID))
	Expect(res.Value).To(Equal(value))
}

func AssertMapping(res MappingRes, diffID, headerID int64, key, value string) {
	Expect(res.DiffID).To(Equal(diffID))
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.Key).To(Equal(key))
	Expect(res.Value).To(Equal(value))
}

func AssertDoubleMapping(res DoubleMappingRes, diffID, headerID int64, keyOne, keyTwo, value string) {
	Expect(res.DiffID).To(Equal(diffID))
	Expect(res.HeaderID).To(Equal(headerID))
	Expect(res.KeyOne).To(Equal(keyOne))
	Expect(res.KeyTwo).To(Equal(keyTwo))
	Expect(res.Value).To(Equal(value))
}
