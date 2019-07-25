package main

import (
	"context"

	"github.com/boxgo/logger"
)

func main() {
	logger.Default.ConfigWillLoad(context.Background())
	logger.Default.ConfigDidLoad(context.Background())

	logger.Info("aaa")
	logger.Default.Info("aaa")
}
