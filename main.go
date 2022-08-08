package main

import (
	"platform/logging"
	"platform/services"
)

func writeMessage(logger logging.Logger) {
	logger.Info("Brucheion")
}

func main() {
	services.RegisterDefaultServices()
	services.Call(writeMessage)
}
