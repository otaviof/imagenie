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
}
