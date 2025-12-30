// Package domain defines core domain types and invariants.
package domain

// SandboxID identifies a runtime-isolated execution instance.
type SandboxID string

// SandboxSpec describes an isolated execution environment.
type SandboxSpec struct {
	// Execution environemnt (OCI image for now)
	Env EnvRef

	// Command to execute inside the sandbox
	Cmd     []string

	// Working directory inside the sandbox
	WorkDir string

	// Environment variables
	EnvVar map[string]string

	// User to run as (optional, runtime-dependent)
	User string

	// Filesystem mounts shared from host into the sandbox
	Mounts []MountSpec
}

// SandboxResult represents result returned from the Sandbox
type SandboxResult struct {
	ExitCode int
}
