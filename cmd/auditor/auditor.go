package main

import "github.com/sirupsen/logrus"

func main() {
	logger := logrus.New()
	logger.Info("started auditor in the background ...")
}
