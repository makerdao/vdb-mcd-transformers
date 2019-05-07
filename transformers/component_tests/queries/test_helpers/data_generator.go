package test_helpers

import (
	"fmt"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"strings"
)

const (
	headerSql = `INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id, eth_node_fingerprint)
				  VALUES ($1, $2, $3, $4, $5, $6)`
	nodeSql = `INSERT INTO public.eth_nodes (genesis_block, network_id, eth_node_id) VALUES ($1, $2, $3)`
	// Ilk data
	ilkSql     = `INSERT INTO maker.ilks (ilk, name) VALUES ($1, $2) RETURNING id`
	ilkRateSql = `INSERT INTO maker.vat_ilk_rate (block_number, block_hash, ilk_id, rate) VALUES ($1, $2, $3, $4)`
	ilkSpotSql = `INSERT INTO maker.vat_ilk_spot (block_number, block_hash, ilk_id, spot) VALUES ($1, $2, $3, $4)`
	ilkArtSql  = `INSERT INTO maker.vat_ilk_art (block_number, block_hash, ilk_id, art) VALUES ($1, $2, $3, $4)`
	ilkLineSql = `INSERT INTO maker.vat_ilk_line (block_number, block_hash, ilk_id, line) VALUES ($1, $2, $3, $4)`
	ilkDustSql = `INSERT INTO maker.vat_ilk_dust (block_number, block_hash, ilk_id, dust) VALUES ($1, $2, $3, $4)`
	ilkDutySql = `INSERT INTO maker.jug_ilk_duty (block_number, block_hash, ilk_id, duty) VALUES ($1, $2, $3, $4)`
	ilkRhoSql  = `INSERT INTO maker.jug_ilk_rho (block_number, block_hash, ilk_id, rho) VALUES ($1, $2, $3, $4)`
	ilkChopSql = `INSERT INTO maker.cat_ilk_chop (block_number, block_hash, ilk_id, chop) VALUES ($1, $2, $3, $4)`
	ilkFlipSql = `INSERT INTO maker.cat_ilk_flip (block_number, block_hash, ilk_id, flip) VALUES ($1, $2, $3, $4)`
	ilkLumpSql = `INSERT INTO maker.cat_ilk_lump (block_number, block_hash, ilk_id, lump) VALUES ($1, $2, $3, $4)`

	// Urn data
	urnSql    = `INSERT INTO maker.urns (ilk_id, guy) VALUES ($1, $2) RETURNING id`
	urnArtSql = `INSERT INTO maker.vat_urn_art (block_number, block_hash, urn_id, art) VALUES ($1, $2, $3, $4)`
	urnInkSql = `INSERT INTO maker.vat_urn_ink (block_number, block_hash, urn_id, ink) VALUES ($1, $2, $3, $4)`

	// Event data
	// TODO add event data
	// TODO add tx for events
)

type GeneratorState struct {
	db            *postgres.DB
	currentHeader core.Header // Current work header
	ilks          []int64     // Generated ilks
	urns          []int64     // Generated urns
}

func NewGenerator(db *postgres.DB) GeneratorState {
	return GeneratorState{
		db:            db,
		currentHeader: core.Header{},
		ilks:          []int64{},
		urns:          []int64{},
	}
}

// Runs probabilistic generator for random ilk/urn interaction.
func (state *GeneratorState) Run(steps int) {
	/* Unnecessary?
	nodeId := test_config.NewTestNode().ID
	_, nodeErr := state.db.Exec(nodeSql, "GENESIS", 1, nodeId)
	if nodeErr != nil {
		fmt.Println("Could not insert initial node: ", nodeErr)
	}*/

	state.currentHeader = fakes.GetFakeHeaderWithTimestamp(0, 0)
	headerErr := state.InsertCurrentHeader()
	if headerErr != nil {
		panic(fmt.Sprintf("Could not insert initial header: %v", headerErr))
	}

	ilkErr := state.CreateIlk()
	if ilkErr != nil {
		panic(fmt.Sprintf("Could not create initial ilk: %v", ilkErr))
	}
	urnErr := state.CreateUrn()
	if urnErr != nil {
		panic(fmt.Sprintf("Could not create initial urn: %v", urnErr))
	}

	var p float32
	var err error
	for i := 1; i <= steps; i++ {
		state.currentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.currentHeader.Hash = test_data.RandomString(10)
		headerErr := state.InsertCurrentHeader()
		if headerErr != nil {
			fmt.Println("Error inserting current header: ", headerErr)
			continue
		}

		p = rand.Float32()
		if p < 0.2 { // Interact with Ilks
			err = state.TouchIlks()
			if err != nil {
				fmt.Println("Error touching ilks: ", err)
			}
		} else { // Interact with Urns
			err = state.TouchUrns()
			if err != nil {
				fmt.Println("Error touching urns: ", err)
			}
		}
	}
}

// Creates a new ilk, or updates a random one
func (state *GeneratorState) TouchIlks() error {
	p := rand.Float32()
	if p < 0.05 {
		return state.CreateIlk()
	} else {
		return state.UpdateIlk()
	}
}

func (state *GeneratorState) CreateIlk() error {
	hexIlk := test_data.RandomString(10)
	name := strings.ToUpper(test_data.RandomString(5))
	ilkId, err := state.InsertIlk(hexIlk, name)
	if err != nil {
		return err
	}

	err = state.InsertInitialIlkData(ilkId)
	if err != nil {
		return err
	}

	state.ilks = append(state.ilks, ilkId)
	return nil
}

// Updates a random property of a randomly chosen ilk
func (state *GeneratorState) UpdateIlk() error {
	randomIlkId := state.ilks[rand.Intn(len(state.ilks))]
	blockNumber, blockHash := state.GetCurrentBlockAndHash()

	var err error
	p := rand.Float64()
	if p < 0.1 {
		_, err = state.db.Exec(ilkRateSql, blockNumber, blockHash, randomIlkId, rand.Int())
	} else {
		_, err = state.db.Exec(ilkSpotSql, blockNumber, blockHash, randomIlkId, rand.Int())
	}
	return err
}

func (state *GeneratorState) TouchUrns() error {
	p := rand.Float32()
	if p < 0.1 {
		return state.CreateUrn()
	} else {
		return state.UpdateUrn()
	}
}

// Creates a new urn associated with a random ilk
func (state *GeneratorState) CreateUrn() error {
	randomIlkId := state.ilks[rand.Intn(len(state.ilks))]
	guy := test_data.RandomString(10)
	urnId, err := state.InsertUrn(randomIlkId, guy)
	if err != nil {
		return err
	}

	blockNumber := state.currentHeader.BlockNumber
	blockHash := state.currentHeader.Hash

	tx, _ := state.db.Beginx()
	_, artErr := tx.Exec(urnArtSql, blockNumber, blockHash, urnId, rand.Int())
	_, inkErr := tx.Exec(urnInkSql, blockNumber, blockHash, urnId, rand.Int())

	// TODO insert urn event data?

	if artErr != nil && inkErr != nil {
		_ = tx.Rollback()
		return fmt.Errorf("Error creating urn.\n artErr: %v\ninkErr: %v", artErr, inkErr)
	}

	_ = tx.Commit()
	state.urns = append(state.urns, urnId)
	return nil
}

// Updates ink or art on a random urn
func (state *GeneratorState) UpdateUrn() error {
	randomUrnId := state.urns[rand.Intn(len(state.urns))]
	blockNumber := state.currentHeader.BlockNumber
	blockHash := state.currentHeader.Hash

	var err error
	p := rand.Float32()
	if p < 0.5 {
		// Update ink
		_, err = state.db.Exec(urnInkSql, blockNumber, blockHash, randomUrnId, rand.Int())
	} else {
		// Update art
		_, err = state.db.Exec(urnArtSql, blockNumber, blockHash, randomUrnId, rand.Int())
	}
	return err
}

// Inserts into `urns` table, returning the urn_id from the database
func (state *GeneratorState) InsertUrn(ilkId int64, guy string) (int64, error) {
	var id int64
	err := state.db.QueryRow(urnSql, ilkId, guy).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting urn: %v", err)
	}
	state.urns = append(state.urns, id)
	return id, nil
}

// Inserts into `ilks` table, returning the ilk_id from the database
func (state *GeneratorState) InsertIlk(hexIlk, name string) (int64, error) {
	var id int64
	err := state.db.QueryRow(ilkSql, hexIlk, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting ilk: %v", err)
	}
	state.ilks = append(state.ilks, id)
	return id, nil
}

func (state *GeneratorState) InsertInitialIlkData(ilkId int64) error {
	tx, _ := state.db.Beginx()
	blockNumber, blockHash := state.GetCurrentBlockAndHash()
	intInsertions := []string{
		ilkRateSql,
		ilkSpotSql,
		ilkArtSql,
		ilkLineSql,
		ilkDustSql,
		ilkDutySql,
		ilkRhoSql,
		ilkChopSql,
		ilkLumpSql,
	}

	for _, intInsertSql := range intInsertions {
		_, err := tx.Exec(intInsertSql, blockNumber, blockHash, ilkId, rand.Int())
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error inserting initial ilk data: %v", err)
		}
	}
	_, flipErr := tx.Exec(ilkFlipSql, blockNumber, blockHash, ilkId, test_data.RandomString(10))
	if flipErr != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error inserting initial ilk data: %v", flipErr)
	}

	_ = tx.Commit()
	return nil
}

func (state *GeneratorState) InsertCurrentHeader() error {
	header := state.currentHeader
	nodeId := test_config.NewTestNode().ID
	_, err := state.db.Exec(headerSql, header.Hash, header.BlockNumber, header.Raw, header.Timestamp, 1, nodeId)
	return err
}

func (state *GeneratorState) GetCurrentBlockAndHash() (int64, string) {
	return state.currentHeader.BlockNumber, state.currentHeader.Hash
}
