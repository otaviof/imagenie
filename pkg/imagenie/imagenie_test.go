package imagenie

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestImagenie_NewImagenie(t *testing.T) {
	cfg := &Config{
		FromImage:   "alpine:latest",
		BaseImage:   "alpine:latest",
		TargetImage: "imagenie:testing",
	}
	i, err := NewImagenie(cfg)

	require.NoError(t, err, "new imagenie should not error")
	require.NotNil(t, i, "imagenie instance should not be nil")

	t.Run("copy", func(t *testing.T) {
		copyPaths := CopyPaths{
			"/etc/alpine-release": "/tmp",
			"/etc/os-release":     "/tmp",
		}
		err := i.Copy(copyPaths)
		require.NoError(t, err, "copy paths should not return error")
	})

	t.Run("labels", func(t *testing.T) {
		i.sourceMgr.SetLabel("k", "v")

		i.Labels()

		labels := i.TargetMgr.Labels()
		require.Contains(t, labels, "k", "target container should have label")
	})

	t.Run("runall", func(t *testing.T) {
		// TODO: investigate and fix errors on running commands;
		t.Skip()
		err := i.RunAll([]string{"ls -l /tmp"})
		require.NoError(t, err, "should not error when running commands")
	})

	t.Run("cleanup", func(t *testing.T) {
		err := i.CleanUp()
		require.NoError(t, err, "should not error on cleanup")
	})
}
