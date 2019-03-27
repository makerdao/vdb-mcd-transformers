package vat_fold

type VatFoldModel struct {
	Ilk              string
	Urn              string `db:"urn_id"`
	Rate             string
	LogIndex         uint   `db:"log_idx"`
	TransactionIndex uint   `db:"tx_idx"`
	Raw              []byte `db:"raw_log"`
}
