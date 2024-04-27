package logger

import "sync"

type LogDestination int

const (
	LogDestinationStdout LogDestination = iota
	LogDestinationFile
	LogDestinationNetwork
)

type Logger struct {
	logMessages []string
	lock        sync.Mutex
	destination LogDestination
	filePath    string
	networkURL  string
}
