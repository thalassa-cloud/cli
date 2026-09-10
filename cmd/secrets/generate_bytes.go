package secrets

import "fmt"

const (
	generateBytesMin     = 16
	generateBytesMax     = 4096
	generateBytesDefault = 32
)

func generateBytesFlagUsage() string {
	return fmt.Sprintf("Generate a random secret of this many bytes (%d-%d; bare --generate-bytes uses %d)",
		generateBytesMin, generateBytesMax, generateBytesDefault)
}

func validateGenerateBytes(n int) error {
	if n < generateBytesMin || n > generateBytesMax {
		return fmt.Errorf("--generate-bytes must be between %d and %d (got %d)", generateBytesMin, generateBytesMax, n)
	}
	return nil
}
