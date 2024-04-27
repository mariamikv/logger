package main

import "logger/logger"

func main() {
	//  usage for file logging
	fileLogger := logger.GetInstance(logger.LogDestinationFile, "application.log", "")
	fileLogger.Info("Application started with file logging")
	fileLogger.Warning("Low memory condition detected")
	fileLogger.Error("An error occurred during processing")

	//  usage for stdout logging
	stdoutLogger := logger.GetInstance(logger.LogDestinationStdout, "", "")
	stdoutLogger.Info("Application can also log to stdout")

	networkLogger := logger.GetInstance(logger.LogDestinationNetwork, "", "http://localhost:8080/logs")
	networkLogger.Info("Network logging needs implementation")
}
