package repositories

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFormatRepositorySize(t *testing.T) {
	const kb = 1024
	tests := []struct {
		name string
		size int64
		want string
	}{
		{name: "empty", size: 0, want: "-"},
		{name: "bytes", size: 512, want: "512 B"},
		{name: "kilobytes", size: 2 * kb, want: "2.0 KB"},
		{name: "megabytes", size: 3 * kb * kb, want: "3.0 MB"},
		{name: "gigabytes", size: 150 * kb * kb * kb, want: "150.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, formatRepositorySize(tt.size))
		})
	}
}
