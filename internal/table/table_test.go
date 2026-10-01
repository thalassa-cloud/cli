package table

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrintWithOptionsBorderless(t *testing.T) {
	var out strings.Builder
	PrintWithOptions(&out,
		[]string{"ID", "Name"},
		[][]string{{"k8s-d0ph3r39gg3c73f6h3ug", "support-cluster"}},
		DefaultOptions(),
	)

	rendered := out.String()
	require.NotEmpty(t, rendered)
	assert.NotContains(t, rendered, "┌")
	assert.NotContains(t, rendered, "│")

	lines := strings.Split(strings.TrimRight(rendered, "\n"), "\n")
	require.Len(t, lines, 2)

	headerCols := strings.Split(lines[0], "\t")
	bodyCols := strings.Split(lines[1], "\t")
	require.GreaterOrEqual(t, len(headerCols), 2)
	require.GreaterOrEqual(t, len(bodyCols), 2)
	assert.Equal(t, len(headerCols[0]), len(bodyCols[0]), "columns should share a width: %q", rendered)
	assert.Contains(t, headerCols[0], "ID")
	assert.Contains(t, bodyCols[0], "k8s-d0ph3r39gg3c73f6h3ug")
	assert.Contains(t, headerCols[1], "NAME")
	assert.Contains(t, bodyCols[1], "support-cluster")
}

func TestPrintWithOptionsBorders(t *testing.T) {
	options := DefaultOptions()
	options.UseBorders = true

	var out strings.Builder
	PrintWithOptions(&out, []string{"ID"}, [][]string{{"1"}}, options)

	rendered := out.String()
	require.NotEmpty(t, rendered)
	assert.Contains(t, rendered, "ID")
	assert.Contains(t, rendered, "1")
}
