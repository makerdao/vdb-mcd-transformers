package main

import (
	"os"

	"github.com/makerdao/vdb-mcd-transformers/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{
		PrettyPrint: true,
	})
	file, err := os.OpenFile("vdb-mcd-transformers.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logrus.SetOutput(file)
	} else {
		logrus.Info("Failed to log to file, using default stderr")
	}

	cmd.Execute()
}
