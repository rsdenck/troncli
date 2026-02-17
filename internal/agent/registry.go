package agent

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// AgentCapabilities defines the allowed actions and intents for the agent
type AgentCapabilities struct {
	AllowedIntents []string `yaml:"allowed_intents"`
	BlockedIntents []string `yaml:"blocked_intents"`
	RiskLevel      string   `yaml:"risk_level"` // low, medium, high
}

// CapabilityRegistry manages the agent capabilities
type CapabilityRegistry struct {
	Capabilities *AgentCapabilities
	Path         string
}

// NewCapabilityRegistry creates a new registry
func NewCapabilityRegistry(path string) *CapabilityRegistry {
	return &CapabilityRegistry{
		Path: path,
		Capabilities: &AgentCapabilities{
			AllowedIntents: []string{},
			BlockedIntents: []string{},
			RiskLevel:      "low",
		},
	}
}

// Load loads the capabilities from the YAML file
func (r *CapabilityRegistry) Load() error {
	data, err := os.ReadFile(r.Path)
	if err != nil {
		if os.IsNotExist(err) {
			// If file doesn't exist, use defaults
			return nil
		}
		return fmt.Errorf("failed to read capabilities file: %w", err)
	}

	if err := yaml.Unmarshal(data, r.Capabilities); err != nil {
		return fmt.Errorf("failed to parse capabilities file: %w", err)
	}

	return nil
}

// IsIntentAllowed checks if an intent is allowed
func (r *CapabilityRegistry) IsIntentAllowed(intent string) bool {
	// Check blocked first
	for _, blocked := range r.Capabilities.BlockedIntents {
		if blocked == intent {
			return false
		}
	}

	// Check allowed
	for _, allowed := range r.Capabilities.AllowedIntents {
		if allowed == intent {
			return true
		}
	}

	// Default to false if strict mode, or true if risk level is high?
	// For security, default to false unless explicitly allowed.
	return false
}

// Save saves the current capabilities to the file
func (r *CapabilityRegistry) Save() error {
	data, err := yaml.Marshal(r.Capabilities)
	if err != nil {
		return fmt.Errorf("failed to marshal capabilities: %w", err)
	}

	if err := os.WriteFile(r.Path, data, 0644); err != nil {
		return fmt.Errorf("failed to write capabilities file: %w", err)
	}

	return nil
}
