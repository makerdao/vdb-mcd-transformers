package test_data

import (
	"math/big"
	"math/rand"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/makerdao/vdb-mcd-transformers/transformers/shared/constants"
	"github.com/makerdao/vulcanizedb/libraries/shared/factories/event"
	"github.com/makerdao/vulcanizedb/pkg/core"
)

var (
	ilkIdentifier   = "4554482d41"
	UrnAddress      = "0x7fa9ef6609ca7921112231f8f195138ebba2977"
	id              = big.NewInt(32)
	ink             = big.NewInt(123)
	art             = big.NewInt(456)
	due             = big.NewInt(789)
	ClipAddress     = "0xBE8E3e3618f7474F8cB1d074A26afFef007E98FB"
	testBlockNumber = uint64(4)
	RawDogBarkLog   = types.Log{
		Address: common.HexToAddress(Dog130Address()),
		Topics: []common.Hash{
			common.HexToHash(constants.DogBarkSignature()),
			common.HexToHash("0x4554482d41000000000000000000000000000000000000000000000000000000"), // ilk
			common.HexToHash("0x00000000000000000000000007fa9ef6609ca7921112231f8f195138ebba2977"), // urn
			common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000020"), // id
		},
		Data: hexutil.MustDecode("0x" +
			"000000000000000000000000000000000000000000000000000000000000007b" + // ink
			"00000000000000000000000000000000000000000000000000000000000001C8" + // art
			"0000000000000000000000000000000000000000000000000000000000000315" + // due
			"000000000000000000000000be8e3e3618f7474f8cb1d074a26affef007e98fb"), // clip
		BlockNumber: testBlockNumber,
		TxHash:      common.Hash{},
		TxIndex:     0,
		BlockHash:   common.Hash{},
		Index:       0,
		Removed:     false,
	}

	DogBarkEventLog = core.EventLog{
		ID:          int64(rand.Int31()),
		HeaderID:    int64(rand.Int31()),
		Log:         RawDogBarkLog,
		Transformed: false,
	}
	dogBarkModel = event.InsertionModel{
		SchemaName: constants.MakerSchema,
		TableName:  constants.DogBarkTable,
		OrderedColumns: []event.ColumnName{
			event.HeaderFK,
			event.LogFK,
			event.AddressFK,
			constants.IlkColumn,
			constants.UrnColumn,
			constants.InkColumn,
			constants.ArtColumn,
			constants.DueColumn,
			constants.ClipColumn,
			constants.SaleIDColumn,
		},
		ColumnValues: event.ColumnValues{
			event.HeaderFK: DogBarkEventLog.HeaderID,
			event.LogFK:    DogBarkEventLog.ID,
			//event.AddressFK,
			//constants.IlkColumn,
			//constants.UrnColumn,
			constants.InkColumn: ink.String(),
			constants.ArtColumn: art.String(),
			constants.DueColumn: due.String(),
			//constants.ClipColumn,
			constants.SaleIDColumn: id.String(),
		},
	}
)

func DogBarkModel() event.InsertionModel { return CopyModel(dogBarkModel) }
