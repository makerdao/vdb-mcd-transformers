package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/jmoiron/sqlx"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/cat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/jug"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/test_helpers"
	"github.com/makerdao/vdb-mcd-transformers/transformers/storage/vat"
	"github.com/makerdao/vdb-mcd-transformers/transformers/test_data"
	"github.com/makerdao/vulcanizedb/libraries/shared/storage"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres"
	"github.com/makerdao/vulcanizedb/pkg/fakes"
	"golang.org/x/crypto/sha3"
)

const (
	headerSql = `INSERT INTO public.headers (hash, block_number, raw, block_timestamp, eth_node_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
	nodeSql = `INSERT INTO public.eth_nodes (genesis_block, network_id, eth_node_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING`
	txSql   = `INSERT INTO public.transactions (header_id, hash, tx_from, tx_index, tx_to)
		VALUES ($1, $2, $3, $4, $5)`
	insertIlkQuery = `INSERT INTO maker.ilks (ilk, identifier) VALUES ($1, $2) RETURNING id`
	insertUrnQuery = `INSERT INTO maker.urns (identifier, ilk_id) VALUES ($1, $2) RETURNING id`
	insertLogSql   = `
		WITH insertedAddressId AS (
			INSERT INTO public.addresses (address) VALUES ('0x1234567890123456789012345678901234567890') ON CONFLICT DO NOTHING RETURNING id
		),
		selectedAddressId AS (
			SELECT id FROM public.addresses WHERE address = '0x1234567890123456789012345678901234567890'
		)
		INSERT INTO public.event_logs (header_id, address) VALUES ($1, (
			SELECT id FROM insertedAddressId
			UNION
			SELECT id FROM selectedAddressId
		)) RETURNING id`
	insertVatFrobSql = `INSERT INTO maker.vat_frob (header_id, urn_id, v, w, dink, dart, log_id)
		VALUES($1::NUMERIC, $2::NUMERIC, $3, $4, $5::NUMERIC, $6::NUMERIC, $7::NUMERIC)
		ON CONFLICT (header_id, log_id)
		DO UPDATE SET urn_id = $2, v = $3, w = $4, dink = $5, dart = $6;`
	insertSpotPokeQuery = `INSERT INTO maker.spot_poke (header_id, ilk_id, value, spot, log_id)
		VALUES($1, $2, $3::NUMERIC, $4::NUMERIC, $5)
		ON CONFLICT (header_id, log_id) DO UPDATE SET ilk_id = $2, value = $3, spot = $5;`
)

var (
	node = core.Node{
		GenesisBlock: "GENESIS",
		NetworkID:    1,
		ID:           "b6f90c0fdd8ec9607aed8ee45c69322e47b7063f0bfb7a29c8ecafab24d0a22d24dd2329b5ee6ed4125a03cb14e57fd584e67f9e53e6c631055cbbd82f080845",
		ClientName:   "Geth/v1.7.2-stable-1db4ecdc/darwin-amd64/go1.9",
	}
	emptyRaw, _ = json.Marshal("nothing")
)

func main() {
	const defaultConnectionString = "postgres://vulcanize:vulcanize@localhost:5432/vulcanize_private?sslmode=disable"
	connectionStringPtr := flag.String("pg-connection-string", defaultConnectionString,
		"postgres connection string")
	stepsPtr := flag.Int("steps", 100, "number of interactions to generate")
	seedPtr := flag.Int64("seed", -1,
		"optional seed for repeatability. Running same seed several times will lead to database constraint violations.")
	flag.Parse()

	db, connectErr := sqlx.Connect("postgres", *connectionStringPtr)
	if connectErr != nil {
		fmt.Println("Could not connect to DB: ", connectErr)
		os.Exit(1)
	}

	pg := postgres.DB{
		DB:     db,
		Node:   node,
		NodeID: 0,
	}

	if *seedPtr != -1 {
		rand.Seed(*seedPtr)
		fmt.Println("\nUsing passed seed. If data from this seed is already in the DB, there will be database constraint errors.")
	} else {
		seed := time.Now().UnixNano()
		rand.Seed(seed)
		fmt.Printf("\nUsing current time as seed: %v. Pass this with '-seed' to reproduce results on a fresh DB.\n", seed)
	}

	fmt.Println("\nRunning this will write mock data to the DB you specified, possibly contaminating real data:")
	fmt.Println(*connectionStringPtr)
	fmt.Println("------------------------------")
	fmt.Print("Do you want to continue? (y/n)")

	var input string
	_, err := fmt.Scanln(&input)
	if input != "y" || err != nil {
		os.Exit(0)
	}

	startTime := time.Now()
	generatorState := NewGenerator(&pg)
	runErr := generatorState.Run(*stepsPtr)
	if runErr != nil {
		fmt.Println("Error occurred while running generator: ", runErr.Error())
		fmt.Println("Exiting without writing any data to DB.")
		os.Exit(1)
	}

	duration := time.Now().Sub(startTime)
	speed := float64(*stepsPtr) / duration.Seconds()
	fmt.Printf("Simulated %v interactions in %v. (%.f/s)\n",
		*stepsPtr, duration.Round(time.Duration(time.Second)).String(), speed)
}

type GeneratorState struct {
	db            *postgres.DB
	currentHeader core.Header // Current work header (Read-only everywhere except in Run)
	currentDiffID int64       // Current diff record
	ilks          []int64     // Generated ilks
	urns          []int64     // Generated urns
	pgTx          *sqlx.Tx
}

func NewGenerator(db *postgres.DB) GeneratorState {
	return GeneratorState{
		db:            db,
		currentHeader: core.Header{},
		ilks:          []int64{},
		urns:          []int64{},
		pgTx:          nil,
	}
}

// Runs probabilistic generator for random ilk/urn interaction.
func (state *GeneratorState) Run(steps int) error {
	pgTx, txErr := state.db.Beginx()
	if txErr != nil {
		return txErr
	}

	state.pgTx = pgTx
	initErr := state.doInitialSetup()
	if initErr != nil {
		return initErr
	}

	var p float32
	var err error

	for i := 1; i <= steps; i++ {
		state.currentHeader = fakes.GetFakeHeaderWithTimestamp(int64(i), int64(i))
		state.currentHeader.Hash = test_data.AlreadySeededRandomString(10)
		headerErr := state.insertCurrentHeader()
		if headerErr != nil {
			return fmt.Errorf("error inserting current header: %v", headerErr)
		}

		diffErr := state.insertDiffRecord()
		if diffErr != nil {
			return fmt.Errorf("error inserting current storage diff: %v", diffErr)
		}

		p = rand.Float32()
		if p < 0.2 { // Interact with Ilks
			err = state.touchIlks()
			if err != nil {
				return fmt.Errorf("error touching ilks: %v", err)
			}
		} else { // Interact with Urns
			err = state.touchUrns()
			if err != nil {
				return fmt.Errorf("error touching urns: %v", err)
			}
		}
	}
	return state.pgTx.Commit()
}

// Creates a starting ilk and urn, with the corresponding header.
func (state *GeneratorState) doInitialSetup() error {
	// This may or may not have been initialised, needed for a FK constraint
	_, nodeErr := state.pgTx.Exec(nodeSql, "GENESIS", 1, node.ID)
	if nodeErr != nil {
		return fmt.Errorf("could not insert initial node: %v", nodeErr)
	}

	state.currentHeader = fakes.GetFakeHeaderWithTimestamp(0, 0)
	state.currentHeader.Hash = test_data.AlreadySeededRandomString(10)
	headerErr := state.insertCurrentHeader()
	if headerErr != nil {
		return fmt.Errorf("could not insert initial header: %v", headerErr)
	}

	diffErr := state.insertDiffRecord()
	if diffErr != nil {
		return fmt.Errorf("error inserting initial storage diff: %v", diffErr)
	}

	ilkErr := state.createIlk()
	if ilkErr != nil {
		return fmt.Errorf("could not create initial ilk: %v", ilkErr)
	}
	urnErr := state.createUrn()
	if urnErr != nil {
		return fmt.Errorf("could not create initial urn: %v", urnErr)
	}
	return nil
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
	ilkName := strings.ToUpper(test_data.AlreadySeededRandomString(7))
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

	var eventErr, logErr, storageErr error
	p := rand.Float64()
	newValue := rand.Int()
	if p < 0.1 {
		_, storageErr = state.pgTx.Exec(vat.InsertIlkRateQuery, state.currentDiffID, state.currentHeader.Id, randomIlkId, newValue)
		// Rate is changed in fold, event which isn't included in spec
	} else {
		_, storageErr = state.pgTx.Exec(vat.InsertIlkSpotQuery, state.currentDiffID, state.currentHeader.Id, randomIlkId, newValue)
		var logID int64
		logErr = state.pgTx.QueryRow(insertLogSql, state.currentHeader.Id).Scan(&logID)
		_, eventErr = state.pgTx.Exec(insertSpotPokeQuery,
			state.currentHeader.Id, randomIlkId, newValue, newValue, logID)

		txErr := state.insertCurrentBlockTx()
		if txErr != nil {
			return txErr
		}
	}

	if storageErr != nil {
		return storageErr
	}
	if eventErr != nil {
		return eventErr
	}
	if logErr != nil {
		return logErr
	}

	return nil
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
	guy := getRandomAddress()
	urnId, insertUrnErr := state.insertUrn(randomIlkId, guy)
	if insertUrnErr != nil {
		return insertUrnErr
	}

	ink := rand.Int()
	art := rand.Int()
	_, artErr := state.pgTx.Exec(vat.InsertUrnArtQuery, state.currentDiffID, state.currentHeader.Id, urnId, art)
	_, inkErr := state.pgTx.Exec(vat.InsertUrnInkQuery, state.currentDiffID, state.currentHeader.Id, urnId, ink)
	var logID int64
	logErr := state.pgTx.QueryRow(insertLogSql, state.currentHeader.Id).Scan(&logID)
	_, frobErr := state.pgTx.Exec(insertVatFrobSql,
		state.currentHeader.Id, urnId, guy, guy, ink, art, logID)

	if artErr != nil || inkErr != nil || frobErr != nil || logErr != nil {
		return fmt.Errorf("error creating urn.\n artErr: %v\ninkErr: %v\nfrobErr: %v\n logErr: %v", artErr, inkErr, frobErr, logErr)
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return fmt.Errorf("error creating matching tx: %v", txErr)
	}

	state.urns = append(state.urns, urnId)
	return nil
}

// Updates ink or art on a random urn
func (state *GeneratorState) updateUrn() error {
	randomUrnId := state.urns[rand.Intn(len(state.urns))]
	randomGuy := getRandomAddress()
	newValue := rand.Int()

	// Computing correct diff complicated, also getting correct guy :(

	var frobErr, logErr, updateErr error
	p := rand.Float32()
	if p < 0.5 {
		// Update ink
		_, updateErr = state.pgTx.Exec(vat.InsertUrnInkQuery, state.currentDiffID, state.currentHeader.Id, randomUrnId, newValue)
		var logID int64
		logErr = state.pgTx.QueryRow(insertLogSql, state.currentHeader.Id).Scan(&logID)
		_, frobErr = state.pgTx.Exec(insertVatFrobSql,
			state.currentHeader.Id, randomUrnId, randomGuy, randomGuy, newValue, 0, logID)
	} else {
		// Update art
		_, updateErr = state.pgTx.Exec(vat.InsertUrnArtQuery, state.currentDiffID, state.currentHeader.Id, randomUrnId, newValue)
		var logID int64
		logErr = state.pgTx.QueryRow(insertLogSql, state.currentHeader.Id).Scan(&logID)
		_, frobErr = state.pgTx.Exec(insertVatFrobSql,
			state.currentHeader.Id, randomUrnId, randomGuy, randomGuy, 0, newValue, logID)
	}

	if updateErr != nil {
		return updateErr
	}
	if frobErr != nil {
		return frobErr
	}
	if logErr != nil {
		return logErr
	}

	txErr := state.insertCurrentBlockTx()
	if txErr != nil {
		return txErr
	}

	return nil
}

// Inserts into `urns` table, returning the urn_id from the database
func (state *GeneratorState) insertUrn(ilkId int64, guy string) (int64, error) {
	var id int64
	err := state.pgTx.QueryRow(insertUrnQuery, guy, ilkId).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting urn: %v", err)
	}
	state.urns = append(state.urns, id)
	return id, nil
}

// Inserts into `ilks` table, returning the ilk_id from the database
func (state *GeneratorState) insertIlk(hexIlk, name string) (int64, error) {
	var id int64
	err := state.pgTx.QueryRow(insertIlkQuery, hexIlk, name).Scan(&id)
	if err != nil {
		return -1, fmt.Errorf("error inserting ilk: %v", err)
	}
	state.ilks = append(state.ilks, id)
	return id, nil
}

// Skips initial events for everything, annoying to do individually
func (state *GeneratorState) insertInitialIlkData(ilkId int64) error {
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
		_, err := state.pgTx.Exec(intInsertSql, state.currentDiffID, state.currentHeader.Id, ilkId, rand.Int())
		if err != nil {
			return fmt.Errorf("error inserting initial ilk data: %v", err)
		}
	}
	_, flipErr := state.pgTx.Exec(cat.InsertCatIlkFlipQuery,
		state.currentDiffID, state.currentHeader.Id, ilkId, test_data.AlreadySeededRandomString(10))

	if flipErr != nil {
		return fmt.Errorf("error inserting initial ilk data: %v", flipErr)
	}

	return nil
}

func (state *GeneratorState) insertCurrentHeader() error {
	header := state.currentHeader
	var id int64
	// TODO: derive eth node id from db so this doesn't fail if id 1 does not exist
	err := state.pgTx.QueryRow(headerSql, header.Hash, header.BlockNumber, header.Raw, header.Timestamp, 1).Scan(&id)
	state.currentHeader.Id = id
	return err
}

// Inserts a tx for the current header, with index 0. This matches the events, that are all generated with index 0
func (state *GeneratorState) insertCurrentBlockTx() error {
	txHash := getRandomHash()
	txFrom := getRandomAddress()
	txIndex := 0
	txTo := getRandomAddress()
	_, txErr := state.pgTx.Exec(txSql, state.currentHeader.Id, txHash, txFrom, txIndex, txTo)
	return txErr
}

func (state *GeneratorState) insertDiffRecord() error {
	fakeRawDiff := test_helpers.GetFakeStorageDiffForHeader(state.currentHeader, common.Hash{}, common.Hash{}, common.Hash{})
	storageDiffRepo := storage.NewDiffRepository(state.db)
	diffID, insertDiffErr := storageDiffRepo.CreateStorageDiff(fakeRawDiff)
	state.currentDiffID = diffID
	return insertDiffErr
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

func getRandomAddress() string {
	hash := getRandomHash()
	address := hash[:42]
	return address
}

func getRandomHash() string {
	seed := test_data.AlreadySeededRandomString(10)
	hash := sha3.Sum256([]byte(seed))
	return fmt.Sprintf("0x%x", hash)
}
