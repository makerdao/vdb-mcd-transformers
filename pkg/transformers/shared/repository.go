package shared

import (
	"database/sql"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
)

func MarkHeaderChecked(headerID int64, db *postgres.DB, checkedHeadersColumn string) error {
	_, err := db.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = $2`, headerID, true)
	return err
}

func MarkHeaderCheckedInTransaction(headerID int64, tx *sql.Tx, checkedHeadersColumn string) error {
	_, err := tx.Exec(`INSERT INTO public.checked_headers (header_id, `+checkedHeadersColumn+`)
		VALUES ($1, $2) 
		ON CONFLICT (header_id) DO
			UPDATE SET `+checkedHeadersColumn+` = $2`, headerID, true)
	return err
}

func MissingHeaders(startingBlockNumber, endingBlockNumber int64, db *postgres.DB, checkedHeadersColumn string) ([]core.Header, error) {
	var result []core.Header
	err := db.Select(
		&result,
		`SELECT headers.id, headers.block_number FROM headers
               	LEFT JOIN checked_headers on headers.id = header_id
                WHERE (header_id ISNULL OR `+checkedHeadersColumn+` IS FALSE)
               	AND headers.block_number >= $1
               	AND headers.block_number <= $2
               	AND headers.eth_node_fingerprint = $3`,
		startingBlockNumber,
		endingBlockNumber,
		db.Node.ID,
	)
	return result, err
}
