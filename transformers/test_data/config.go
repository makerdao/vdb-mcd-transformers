package test_data

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"io/ioutil"
)

func SetTestConfig() bool {
	logrus.SetOutput(ioutil.Discard)
	viper.SetConfigName("staging")
	viper.AddConfigPath("$GOPATH/src/github.com/vulcanize/mcd_transformers/environments/")
	return true
}
