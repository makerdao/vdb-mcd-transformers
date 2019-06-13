package vat_grab

import (
	"encoding/json"
	"errors"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"

	"github.com/vulcanize/vulcanizedb/libraries/shared/constants"

	"github.com/ethereum/go-ethereum/common"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
)

type VatGrabConverter struct{}

func (VatGrabConverter) ToModels(ethLogs []types.Log) ([]interface{}, error) {
	var models []interface{}
	for _, ethLog := range ethLogs {
		err := verifyLog(ethLog)
		if err != nil {
			return nil, err
		}
		ilk := ethLog.Topics[1].Hex()
		urn := common.BytesToAddress(ethLog.Topics[2].Bytes()).String()
		v := common.BytesToAddress(ethLog.Topics[3].Bytes()).String()
		wBytes, wErr := shared.GetLogNoteArgumentAtIndex(3, ethLog.Data)
		if wErr != nil {
			return nil, wErr
		}
		w := common.BytesToAddress(wBytes).String()
		dinkBytes, dinkErr := shared.GetLogNoteArgumentAtIndex(4, ethLog.Data)
		if dinkErr != nil {
			return nil, dinkErr
		}
		dink := shared.ConvertInt256HexToBigInt(hexutil.Encode(dinkBytes))
		dartBytes, dartErr := shared.GetLogNoteArgumentAtIndex(5, ethLog.Data)
		if dartErr != nil {
			return nil, dartErr
		}
		dart := shared.ConvertInt256HexToBigInt(hexutil.Encode(dartBytes))

		raw, err := json.Marshal(ethLog)
		if err != nil {
			return nil, err
		}
		model := VatGrabModel{
			Ilk:              ilk,
			Urn:              urn,
			V:                v,
			W:                w,
			Dink:             dink.String(),
			Dart:             dart.String(),
			LogIndex:         ethLog.Index,
			TransactionIndex: ethLog.TxIndex,
			Raw:              raw,
		}
		models = append(models, model)
	}
	return models, nil
}

func verifyLog(log types.Log) error {
	if len(log.Topics) < 4 {
		return errors.New("log missing topics")
	}
	if len(log.Data) < constants.DataItemLength {
		return errors.New("log missing data")
	}
	return nil
}
