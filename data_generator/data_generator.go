package main

import (
	"flag"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/vulcanize/mcd_transformers/test_config"
	"github.com/vulcanize/mcd_transformers/transformers/shared"
	"github.com/vulcanize/mcd_transformers/transformers/storage/cat"
	"github.com/vulcanize/mcd_transformers/transformers/storage/jug"
	"github.com/vulcanize/mcd_transformers/transformers/storage/vat"
	"github.com/vulcanize/mcd_transformers/transformers/test_data"
	"github.com/vulcanize/vulcanizedb/pkg/core"
	"github.com/vulcanize/vulcanizedb/pkg/datastore/postgres"
	"github.com/vulcanize/vulcanizedb/pkg/fakes"
	"math/rand"
	"os"
	"strings"
)

const (
	headerSql = `INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id, eth_node_fingerprint)
				  VALUES ($1, $2, $3, $4, $5, $6)`
	nodeSql = `INSERT INTO public.eth_nodes (genesis_block, network_id, eth_node_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`

	// Event data
	// TODO add event data
	// TODO add tx for events
)

func main() {
	const defaultConnectionString = "postgres://vulcanize:vulcanize@localhost:5432/vulcanize_private?sslmode=disable"
	connectionStringPtr := flag.String("pg-connection-string", defaultConnectionString,
		"postgres connection string")
	stepsPtr := flag.Int("steps", 100, "number of interactions to generate")
	flag.Parse()

	db, connectErr := sqlx.Connect("postgres", *connectionStringPtr)
	if connectErr != nil {
		fmt.Println("Could not connect to DB: ", connectErr)
		os.Exit(1)
	}

	pg := postgres.DB{
		DB:     db,
		Node:   test_config.NewTestNode(),
		NodeID: 0,
	}

	fmt.Println("\nRunning this will write mock data to the DB you specified, possibly contaminating real data.")
	fmt.Println("------------------------------")
	fmt.Print("Do you want to continue? (y/n)")

	var input string
	_, err := fmt.Scanln(&input)
	if input != "y" || err != nil {
		os.Exit(0)
	}

	generatorState := NewGenerator(&pg)
	generatorState.Run(*stepsPtr)
}

type GeneratorState struct {
	db            *postgres.DB
	currentHeader core.Header // Current work header (Read-only everywhere except in Run)
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
	state.doInitialSetup()

	var p float32
	var err error
	for i := 1; i <= steps; i++ {
		state.currentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.currentHeader.Hash = test_data.RandomString(10)
		headerErr := state.insertCurrentHeader()
		if headerErr != nil {
			fmt.Println("Error inserting current header: ", headerErr)
			continue
		}

		p = rand.Float32()
		if p < 0.2 { // Interact with Ilks
			err = state.touchIlks()
			if err != nil {
				fmt.Println("Error touching ilks: ", err)
			}
		} else { // Interact with Urns
			err = state.touchUrns()
			if err != nil {
				fmt.Println("Error touching urns: ", err)
			}
		}
	}
}

// Creates a starting ilk and urn, with the corresponding header.
func (state *GeneratorState) doInitialSetup() {
	// This may or may not have been initialised, needed for a FK constraint
	nodeId := test_config.NewTestNode().ID
	_, nodeErr := state.db.Exec(nodeSql, "GENESIS", 1, nodeId)
	if nodeErr != nil {
		panic(fmt.Sprintf("Could not insert initial node: %v", nodeErr))
	}

	state.currentHeader = fakes.GetFakeHeaderWithTimestamp(0, 0)
	headerErr := state.insertCurrentHeader()
	if headerErr != nil {
		panic(fmt.Sprintf("Could not insert initial header: %v", headerErr))
	}

	ilkErr := state.createIlk()
	if ilkErr != nil {
		panic(fmt.Sprintf("Could not create initial ilk: %v", ilkErr))
	}
	urnErr := state.createUrn()
	if urnErr != nil {
		panic(fmt.Sprintf("Could not create initial urn: %v", urnErr))
	}
}

// Creates a new ilk, or updates a random one
func (state *GeneratorState) touchIlks() error {
	p := rand.Float32()
	if p < 0.05 {
		return state.createIlk()
	} else {
		return state.updateIlk()
	}
}

func (state *GeneratorState) createIlk() error {
	ilkName := strings.ToUpper(test_data.RandomString(5))
	hexIlk := GetHexIlk(ilkName)

	ilkId, insertIlkErr := state.insertIlk(hexIlk, ilkName)
	if insertIlkErr != nil {
		return insertIlkErr
	}

	initIlkErr := state.insertInitialIlkData(ilkId)
	if initIlkErr != nil {
		return initIlkErr
	}

	state.ilks = append(state.ilks, ilkId)
	return nil
}

// Updates a random property of a randomly chosen ilk
func (state *GeneratorState) updateIlk() error {
	randomIlkId := state.ilks[rand.Intn(len(state.ilks))]
	blockNumber, blockHash := state.getCurrentBlockAndHash()

	var err error
	p := rand.Float64()
	if p < 0.1 {
		_, err = state.db.Exec(vat.InsertIlkRateQuery, blockNumber, blockHash, randomIlkId, rand.Int())
	} else {
		_, err = state.db.Exec(vat.InsertIlkSpotQuery, blockNumber, blockHash, randomIlkId, rand.Int())
	}
	return err
}

func (state *GeneratorState) touchUrns() error {
	p := rand.Float32()
	if p < 0.1 {
		return state.createUrn()
	} else {
		return state.updateUrn()
	}
}

// Creates a new urn associated with a random ilk
func (state *GeneratorState) createUrn() error {
	randomIlkId := state.ilks[rand.Intn(len(state.ilks))]
	guy := test_data.RandomString(10)
	urnId, insertUrnErr := state.insertUrn(randomIlkId, guy)
	if insertUrnErr != nil {
		return insertUrnErr
	}

	blockNumber := state.currentHeader.BlockNumber
	blockHash := state.currentHeader.Hash

	tx, _ := state.db.Beginx()
	_, artErr := tx.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, urnId, rand.Int())
	_, inkErr := tx.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, urnId, rand.Int())

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
func (state *GeneratorState) updateUrn() error {
	randomUrnId := state.urns[rand.Intn(len(state.urns))]
	blockNumber := state.currentHeader.BlockNumber
	blockHash := state.currentHeader.Hash

	var err error
	p := rand.Float32()
	if p < 0.5 {
		// Update ink
		_, err = state.db.Exec(vat.InsertUrnInkQuery, blockNumber, blockHash, randomUrnId, rand.Int())
	} else {
		// Update art
		_, err = state.db.Exec(vat.InsertUrnArtQuery, blockNumber, blockHash, randomUrnId, rand.Int())
	}
	return err
}

// Inserts into `urns` table, returning the urn_id from the database
func (state *GeneratorState) insertUrn(ilkId int64, guy string) (int64, error) {
	var id int64
	err := state.db.QueryRow(shared.InsertUrnQuery, guy, ilkId).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting urn: %v", err)
	}
	state.urns = append(state.urns, id)
	return id, nil
}

// Inserts into `ilks` table, returning the ilk_id from the database
func (state *GeneratorState) insertIlk(hexIlk, name string) (int64, error) {
	var id int64
	err := state.db.QueryRow(shared.InsertIlkQuery, hexIlk, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting ilk: %v", err)
	}
	state.ilks = append(state.ilks, id)
	return id, nil
}

func (state *GeneratorState) insertInitialIlkData(ilkId int64) error {
	tx, _ := state.db.Beginx()
	blockNumber, blockHash := state.getCurrentBlockAndHash()
	intInsertions := []string{
		vat.InsertIlkRateQuery,
		vat.InsertIlkSpotQuery,
		vat.InsertIlkArtQuery,
		vat.InsertIlkLineQuery,
		vat.InsertIlkDustQuery,
		jug.InsertJugIlkDutyQuery,
		jug.InsertJugIlkRhoQuery,
		cat.InsertCatIlkChopQuery,
		cat.InsertCatIlkLumpQuery,
	}

	for _, intInsertSql := range intInsertions {
		_, err := tx.Exec(intInsertSql, blockNumber, blockHash, ilkId, rand.Int())
		if err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("error inserting initial ilk data: %v", err)
		}
	}
	_, flipErr := tx.Exec(cat.InsertCatIlkFlipQuery, blockNumber, blockHash, ilkId, test_data.RandomString(10))
	if flipErr != nil {
		_ = tx.Rollback()
		return fmt.Errorf("error inserting initial ilk data: %v", flipErr)
	}

	_ = tx.Commit()
	return nil
}

func (state *GeneratorState) insertCurrentHeader() error {
	header := state.currentHeader
	nodeId := test_config.NewTestNode().ID
	_, err := state.db.Exec(headerSql, header.Hash, header.BlockNumber, header.Raw, header.Timestamp, 1, nodeId)
	return err
}

func (state *GeneratorState) getCurrentBlockAndHash() (int64, string) {
	return state.currentHeader.BlockNumber, state.currentHeader.Hash
}

// UTF-oblivious, names generated with alphanums anyway
func GetHexIlk(ilkName string) string {
	hexIlk := fmt.Sprintf("%x", ilkName)
	unpaddedLength := len(hexIlk)
	for i := unpaddedLength; i < 64; i++ {
		hexIlk = hexIlk + "0"
	}
	return hexIlk
}
