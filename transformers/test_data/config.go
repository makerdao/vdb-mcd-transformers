package test_data

import (
	"io/ioutil"

	"github.com/makerdao/vdb-mcd-transformers/test_config"
	"github.com/sirupsen/logrus"
)

func SetTestConfig() bool {
	logrus.SetOutput(ioutil.Discard)
	test_config.SetTestConfig()
	return true
}
