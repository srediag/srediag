// TODO: Implement build error codes (0-5) for CI/CD integration (see docs/architecture/build.md §9)
// TODO: Document and enforce error code meanings for all mage targets and CLI commands (see docs/architecture/build.md §9)
// Package errorcodes defines standardized build error codes for CI/CD integration.
// See docs/architecture/build.md §7 for details.
package errorcodes

import "fmt"

// Build error codes for CI/CD integration.
const (
	ErrCodeSuccess         = 0 // Success
	ErrCodeGeneral         = 1 // General (unexpected) error
	ErrCodeVersionMismatch = 2 // Version mismatch (YAML ↔ go.mod)
	ErrCodeLintVetFailed   = 3 // Lint or vet failed
	ErrCodeUnitTestsFailed = 4 // Unit tests failed
	ErrCodeCosignFailed    = 5 // Cosign signing or verification failed
)

// ErrorCodeDescriptions maps error codes to their meaning.
var ErrorCodeDescriptions = map[int]string{
	ErrCodeSuccess:         "Success",
	ErrCodeGeneral:         "General (unexpected) error",
	ErrCodeVersionMismatch: "Version mismatch (YAML ↔ go.mod)",
	ErrCodeLintVetFailed:   "Lint or vet failed",
	ErrCodeUnitTestsFailed: "Unit tests failed",
	ErrCodeCosignFailed:    "Cosign signing or verification failed",
}

// PrintErrorCodes prints all error codes and their meanings.
func PrintErrorCodes() {
	fmt.Println("Build Error Codes:")
	for code, desc := range ErrorCodeDescriptions {
		fmt.Printf("  %d: %s\n", code, desc)
	}
}
