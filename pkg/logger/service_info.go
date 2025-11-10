package logger

import (
	"fmt"
	"slices"
	"strings"
)

const (
	DevelopmentEnv           = "dev"
	ProductionEnv            = "prod"
	UserAcceptanceTestingEnv = "uat"
)

type ServiceInfo struct {
	Name    string
	Version string
	Env     string
}

// Validate checks whether the ServiceInfo fields are properly populated and valid.
//
// It performs the following validations:
//   - Trims leading/trailing spaces for all fields.
//   - Ensures Name, Version, and Env are not empty.
//   - Ensures Env is one of the allowed values: "production", "uat", or "develop".
func (s *ServiceInfo) Validate() error {
	// Normalize input fields
	s.Name = strings.TrimSpace(s.Name)
	s.Version = strings.TrimSpace(s.Version)
	s.Env = strings.TrimSpace(strings.ToLower(s.Env))

	// Check for missing required fields
	var missing []string
	if s.Name == "" {
		missing = append(missing, "Name")
	}
	if s.Version == "" {
		missing = append(missing, "Version")
	}
	if s.Env == "" {
		missing = append(missing, "Env")
	}
	if len(missing) > 0 {
		return fmt.Errorf("%w: missing fields [%s]", ErrInvalidServiceInfo, strings.Join(missing, ", "))
	}

	// Validate environment value
	validEnvs := []string{DevelopmentEnv, UserAcceptanceTestingEnv, ProductionEnv}
	if !slices.Contains(validEnvs, s.Env) {
		return fmt.Errorf("%w: invalid env %q (allowed: %s)", ErrInvalidServiceInfo, s.Env, strings.Join(validEnvs, ", "))
	}

	return nil
}
