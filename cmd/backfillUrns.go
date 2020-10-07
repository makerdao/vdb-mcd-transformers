package cmd

import (
	"fmt"
	"sync"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/makerdao/vdb-mcd-transformers/backfill"
	"github.com/makerdao/vdb-mcd-transformers/backfill/fork"
	"github.com/makerdao/vdb-mcd-transformers/backfill/frob"
	"github.com/makerdao/vdb-mcd-transformers/backfill/grab"
	"github.com/makerdao/vdb-mcd-transformers/backfill/repository"
	"github.com/makerdao/vdb-mcd-transformers/backfill/shared"
	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/makerdao/vulcanizedb/pkg/core"
	"github.com/makerdao/vulcanizedb/pkg/datastore/postgres/repositories"
	"github.com/makerdao/vulcanizedb/pkg/eth"
	"github.com/makerdao/vulcanizedb/pkg/eth/client"
	"github.com/makerdao/vulcanizedb/pkg/eth/converters"
	"github.com/makerdao/vulcanizedb/pkg/eth/node"
	"github.com/makerdao/vulcanizedb/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type backFillInitializer func(core.BlockChain, repository.EventsRepository, repository.StorageRepository, shared.DartDinkRetriever) backfill.BackFiller

var (
	eventsToBackFill []string
	initializers     = map[string]backFillInitializer{
		backfill.ForkEvent: fork.NewForkBackFiller,
		backfill.FrobEvent: frob.NewFrobBackFiller,
		backfill.GrabEvent: grab.NewGrabBackFiller,
	}
	startingBlock int
)

// backfillUrnsCmd represents the backfillUrns command
var backfillUrnsCmd = &cobra.Command{
	Use:   "backfillUrns",
	Short: "Backfill diffs for urns, looking up diffs based on associated events",
	Long: `Fetch diffs when events indicate the state of an Urn changed at a given block.
Optionally pass a starting block number to backfill since a given block.
Optionally pass events to watch (fork, frob, grab) to backfill based off of certain events.`,
	PreRun: setViperConfigs,
	Run: func(cmd *cobra.Command, args []string) {
		err := backfillUrns()
		if err != nil {
			logrus.Fatalf("error backfilling urns: %s", err.Error())
		}
		logrus.Println("Backfilling urns completed successfully")
		return
	},
}

func init() {
	rootCmd.AddCommand(backfillUrnsCmd)

	backfillUrnsCmd.Flags().IntVarP(&startingBlock, "starting-block", "s", 0, "starting block for backfilling diffs derived from urn events")
	backfillUrnsCmd.Flags().StringSliceVarP(&eventsToBackFill, "events-to-backfill", "e", []string{"fork", "frob", "grab"}, "events to back-fill")
	backfillUrnsCmd.Flags().String("database-name", "vulcanize_public", "database name")
	backfillUrnsCmd.Flags().Int("database-port", 5432, "database port")
	backfillUrnsCmd.Flags().String("database-hostname", "localhost", "database hostname")
	backfillUrnsCmd.Flags().String("database-user", "", "database user")
	backfillUrnsCmd.Flags().String("database-password", "", "database password")
	backfillUrnsCmd.Flags().String("client-ipcPath", "", "location of geth.ipc file")

	viper.BindPFlag("database.name", backfillUrnsCmd.Flags().Lookup("database-name"))
	viper.BindPFlag("database.port", backfillUrnsCmd.Flags().Lookup("database-port"))
	viper.BindPFlag("database.hostname", backfillUrnsCmd.Flags().Lookup("database-hostname"))
	viper.BindPFlag("database.user", backfillUrnsCmd.Flags().Lookup("database-user"))
	viper.BindPFlag("database.password", backfillUrnsCmd.Flags().Lookup("database-password"))
	viper.BindPFlag("client.ipcPath", backfillUrnsCmd.Flags().Lookup("client-ipcPath"))
}

func setViperConfigs(cmd *cobra.Command, args []string) {
	ipc = viper.GetString("client.ipcpath")
	databaseConfig = config.Database{
		Name:     viper.GetString("database.name"),
		Hostname: viper.GetString("database.hostname"),
		Port:     viper.GetInt("database.port"),
		User:     viper.GetString("database.user"),
		Password: viper.GetString("database.password"),
	}
	viper.Set("database.config", databaseConfig)
}

func backfillUrns() error {
	validationErr := backfill.ValidateArgs(eventsToBackFill)
	if validationErr != nil {
		return fmt.Errorf("invalid events-to-backfill: %w", validationErr)
	}

	blockChain := getBlockChain()
	db := utils.LoadPostgres(databaseConfig, blockChain.Node())
	eventRepository := repository.NewEventsRepository(&db)
	storageRepository := repository.NewStorageRepository(&db)

	var wg sync.WaitGroup
	done := make(chan bool)
	errs := make(chan error)

	for _, e := range eventsToBackFill {
		initializer := initializers[e]
		headerRepository := repositories.NewHeaderRepository(&db)
		dartDinkRetriever := shared.NewDartDinkRetriever(blockChain, eventRepository, headerRepository, storageRepository)
		backFiller := initializer(blockChain, eventRepository, storageRepository, dartDinkRetriever)
		wg.Add(1)
		go backFillEvents(backFiller, startingBlock, errs, &wg)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case err := <-errs:
		logrus.Errorf("error executing back-fill: %s", err.Error())
		return err
	}

	return nil
}

func backFillEvents(backFiller backfill.BackFiller, startingBlock int, errs chan error, wg *sync.WaitGroup) {
	defer wg.Done()
	err := backFiller.BackFill(startingBlock)
	if err != nil {
		errs <- err
	}
	return
}

func getBlockChain() *eth.BlockChain {
	rpcClient, ethClient := getClients()
	vdbEthClient := client.NewEthClient(ethClient)
	vdbNode := node.MakeNode(rpcClient)
	transactionConverter := converters.NewTransactionConverter(ethClient)
	return eth.NewBlockChain(vdbEthClient, rpcClient, vdbNode, transactionConverter)
}

func getClients() (client.RpcClient, *ethclient.Client) {
	rawRpcClient, err := rpc.Dial(ipc)

	if err != nil {
		logrus.Fatal(err)
	}
	rpcClient := client.NewRpcClient(rawRpcClient, ipc)
	ethClient := ethclient.NewClient(rawRpcClient)

	return rpcClient, ethClient
}
