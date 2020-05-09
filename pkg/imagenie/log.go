package imagenie

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

// LogLevelEnv log-level environment variable name.
const LogLevelEnv = "LOGLEVEL"

// init configure log level and output.
func init() {
	log.SetOutput(os.Stdout)
	SetLogLevel()
}

// SetLogLevel set logrus log-level based on environment variable value.
func SetLogLevel() {
	if level := os.Getenv(LogLevelEnv); level != "" {
		if logLevel, err := strconv.Atoi(level); err == nil {
			log.SetLevel(log.Level(logLevel))
		}
	} else {
		log.SetLevel(log.InfoLevel)
	}
}
