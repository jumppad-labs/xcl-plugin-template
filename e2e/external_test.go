package e2e

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jumppad-labs/xcl"

	"github.com/jumppad-labs/xcl-plugin-template/entities"
)

func TestExternalAppliesPlansNoChangesAndDestroys(t *testing.T) {
	outDir := t.TempDir()
	c := newExternalConfig(t, t.TempDir(), outDir)

	err := c.Apply(basicExample)
	require.NoError(t, err)

	notePath := filepath.Join(outDir, "greeting.txt")

	content, err := os.ReadFile(notePath)
	require.NoError(t, err)
	require.Equal(t, "Hello from xcl", string(content))

	changes, err := c.Diff([]string{basicExample})
	require.NoError(t, err)
	require.Equal(t, 0, changes.Changed())

	err = c.Destroy()
	require.NoError(t, err)

	_, err = os.Stat(notePath)
	require.ErrorIs(t, err, os.ErrNotExist)
}

func TestExternalAppliesTheFunctionResultAsTheNoteName(t *testing.T) {
	outDir := t.TempDir()
	c := newExternalConfig(t, t.TempDir(), outDir)

	err := c.Apply(basicExample)
	require.NoError(t, err)

	note, err := xcl.Find[entities.Note](c, noteAddress)
	require.NoError(t, err)
	require.Equal(t, "greeting.txt", note.Name)
	require.True(t, strings.HasSuffix(note.Path, "greeting.txt"), "path %q does not end with greeting.txt", note.Path)

	_, err = os.Stat(filepath.Join(outDir, "greeting.txt"))
	require.NoError(t, err)

	err = c.Destroy()
	require.NoError(t, err)
}

func TestExternalReportsTheFunctionsError(t *testing.T) {
	c := newExternalConfig(t, t.TempDir(), t.TempDir())

	err := c.Apply(writeConfiguration(t, untitledNote))

	require.Error(t, err)
	require.Contains(t, err.Error(), "at least one letter or digit")
}
