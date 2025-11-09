package logger

import (
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	defaultPath = "./logs/app.log"
)

var (
	ErrInvalidFilePath    = errors.New("invalid file path")
	ErrInvalidServiceInfo = errors.New("invalid service info")
)

type config struct {
	ServiceInfo ServiceInfo

	// Config log to file
	EnableFileLog bool
	FilePath      string

	// Rotation options
	MaxSize    int  // Max size (MB) before rotation, default 100
	MaxAge     int  // Max age (days) to retain old files
	MaxBackups int  // Max number of backup files to keep
	Compress   bool // Compress rotated files (gzip)
}

type ConfigOption func(conf *config) error

func getDefaultConfig(serviceInfo ServiceInfo) *config {
	return &config{
		ServiceInfo:   serviceInfo,
		EnableFileLog: false,
		FilePath:      defaultPath,
		MaxSize:       100, // MB
		MaxAge:        28,  //days
		MaxBackups:    3,
		Compress:      false,
	}
}

// EnableFileLogging returns a ConfigOption that enables or disables
// writing logs to a file in addition to stdout.
//
// When set to true, the logger writes logs to both console (stdout)
// and the configured file path. When set to false, only console output is used.
//
// Example:
//
//	logger.New(EnableFileLogging(true))
func EnableFileLogging(enable bool) ConfigOption {
	return func(conf *config) error {
		conf.EnableFileLog = enable
		return nil
	}
}

// WithLogFile returns a ConfigOption that validates and sets the log file path.
//
// It ensures the following:
//   - Uses a default path if the provided one is empty.
//   - Converts the path to an absolute path.
//   - Rejects any path that points outside the allowed "./logs" directory.
//   - Returns an error if the path is invalid or not allowed.
//
// Example:
//
//	logger.New(WithLogFile("./logs/app.log"))
func WithLogFile(filePath string) ConfigOption {
	return func(conf *config) error {
		// Normalize and validate the input path
		fp := strings.TrimSpace(filePath)
		if fp == "" {
			// Warn and fallback to default path if user didn't provide one
			log.Warn().
				Str("service_name", conf.ServiceInfo.Name).
				Str("service_version", conf.ServiceInfo.Version).
				Str("service_env", conf.ServiceInfo.Env).
				Msg("log file path not provided, falling back to default path")
			fp = defaultPath
		}

		// Convert to absolute path to prevent relative path traversal (e.g., ../../)
		absPath, err := filepath.Abs(fp)
		if err != nil {
			return fmt.Errorf("%w: %v", ErrInvalidFilePath, err)
		}

		// Determine the allowed root folder for logs
		root, err := filepath.Abs("./logs")
		if err != nil {
			return fmt.Errorf("cannot resolve log root folder: %w", err)
		}

		// Ensure the provided file path is inside the allowed root directory
		if !strings.HasPrefix(absPath, root) {
			return fmt.Errorf("%w: log file path %s is outside allowed folder %s", ErrInvalidFilePath, absPath, root)
		}

		// Save validated path to configuration
		conf.FilePath = absPath
		return nil
	}
}

// WithLogRotation configures the log file rotation settings for the logger.
// Params:
//   - maxSize: maximum size in MB of a log file before it gets rotated (must > 0)
//   - maxAge: maximum number of days to retain old log files (cannot be negative, 0 = keep indefinitely)
//   - maxBackups: maximum number of old log files to retain (cannot be negative)
//   - compress: whether to gzip compress rotated log files
//
// Example:
//
//	logger.New(WithLogRotation(100, 7, 3, false))
func WithLogRotation(maxSize, maxAge, maxBackups int, compress bool) ConfigOption {
	return func(conf *config) error {
		if maxSize <= 0 {
			return fmt.Errorf("maxSize must be > 0")
		}
		if maxAge < 0 {
			return fmt.Errorf("maxAge cannot be negative")
		}
		if maxBackups < 0 {
			return fmt.Errorf("maxBackups cannot be negative")
		}

		conf.MaxSize = maxSize
		conf.MaxAge = maxAge
		conf.MaxBackups = maxBackups
		conf.Compress = compress

		return nil
	}
}
