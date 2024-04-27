package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

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

var instance *Logger
var once sync.Once

func GetInstance(
	destination LogDestination,
	filePath string,
	networkUrl string,
) *Logger {
	once.Do(func() {
		instance = &Logger{
			logMessages: []string{},
			lock:        sync.Mutex{},
			destination: destination,
		}
		if destination == LogDestinationFile {
			instance.filePath = filePath
		} else if destination == LogDestinationNetwork {
			instance.networkURL = networkUrl
		}
	})
	return instance
}

func (l *Logger) writeToFile(messages []string) error {
	data := []byte(strings.Join(messages, "\n") + "\n")
	err := os.WriteFile(l.filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing logs to file: %w", err)
	}
	return nil
}

func (l *Logger) writeToNetwork(messages []string) error {
	data, err := json.Marshal(messages)
	if err != nil {
		return fmt.Errorf("error marshalling log messages to JSON: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, l.networkURL, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending logs to network: %w", err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			_ = fmt.Errorf("error ocuered")
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code from network logging: %d", resp.StatusCode)
	}
	return nil
}

func (l *Logger) writeLog() {
	switch l.destination {
	case LogDestinationStdout:
		for _, message := range l.logMessages {
			fmt.Println(message)
		}
	case LogDestinationFile:
		err := l.writeToFile(l.logMessages)
		if err != nil {
			fmt.Println("Error writing to file:", err)
		}
	case LogDestinationNetwork:
		err := l.writeToNetwork(l.logMessages)
		if err != nil {
			fmt.Println("Error sending logs to network:", err)
		}
	}
	l.logMessages = []string{}
}

func (l *Logger) Info(message string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logMessages = append(l.logMessages, "INFO: "+message)
	l.writeLog()
}

func (l *Logger) Warning(message string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logMessages = append(l.logMessages, "WARNING: "+message)
	l.writeLog()
}

func (l *Logger) Error(message string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.logMessages = append(l.logMessages, "ERROR: "+message)
	l.writeLog()
}
