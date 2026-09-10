package kms

import (
	"fmt"
	"os"
	"strings"

	"github.com/thalassa-cloud/cli/internal/config/securefile"
	clientkms "github.com/thalassa-cloud/client-go/kms"
)

func resolveWireInput(flagValue, fromFile, flagName string) (string, error) {
	if flagValue == "" && fromFile == "" {
		return "", fmt.Errorf("provide --%s or --from-file", flagName)
	}
	if fromFile == "" {
		return flagValue, nil
	}
	data, err := os.ReadFile(fromFile)
	if err != nil {
		return "", fmt.Errorf("read --from-file: %w", err)
	}
	return clientkms.EncodeBytes(data), nil
}

func resolveTextFileInput(flagValue, fromFile, flagName string) (string, error) {
	if flagValue == "" && fromFile == "" {
		return "", fmt.Errorf("provide --%s or --from-file", flagName)
	}
	if fromFile == "" {
		return flagValue, nil
	}
	data, err := os.ReadFile(fromFile)
	if err != nil {
		return "", fmt.Errorf("read --from-file: %w", err)
	}
	return strings.TrimSpace(string(data)), nil
}

func writeTextResult(toFile, value string) error {
	value = strings.TrimSpace(value)
	if toFile != "" {
		if err := securefile.Write(toFile, []byte(value+"\n")); err != nil {
			return fmt.Errorf("write --to-file: %w", err)
		}
		return nil
	}
	fmt.Println(value)
	return nil
}

func writeBinaryResult(toFile string, data []byte) error {
	if toFile != "" {
		if err := securefile.Write(toFile, data); err != nil {
			return fmt.Errorf("write --to-file: %w", err)
		}
		return nil
	}
	fmt.Println(clientkms.EncodeBytes(data))
	return nil
}
