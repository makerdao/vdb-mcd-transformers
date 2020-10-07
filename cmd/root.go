package cmd

import (
	"strings"

	"github.com/makerdao/vulcanizedb/pkg/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	databaseConfig config.Database
	ipc            string
)

var rootCmd = &cobra.Command{
	Use:              "vdb-mcd-transformers",
	PersistentPreRun: initFuncs,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}

func init() {
	// When searching for env variables, replace dots in config keys with underscores
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	rootCmd.PersistentFlags().String("log-level", logrus.InfoLevel.String(), "Log level (trace, debug, info, warn, error, fatal, panic")
	viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))
}

func initFuncs(cmd *cobra.Command, args []string) {
	logLvlErr := logLevel()
	if logLvlErr != nil {
		logrus.Fatalf("Could not set log level: %s", logLvlErr.Error())
	}
}

func logLevel() error {
	lvl, err := logrus.ParseLevel(viper.GetString("log.level"))
	if err != nil {
		return err
	}
	logrus.SetLevel(lvl)
	if lvl > logrus.InfoLevel {
		logrus.SetReportCaller(true)
	}
	logrus.Info("Log level set to ", lvl.String())
	return nil
}
