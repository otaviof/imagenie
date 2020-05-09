package imagenie

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.Level(99))
}

func TestManager_NewManager(t *testing.T) {
	fromImage := "alpine:latest"
	m, err := NewManager(fromImage, "imagenie-test:unit")
	require.NoError(t, err, "new manager should not error")
	require.NotNil(t, m, "manager instance should not be nil")

	t.Run("from", func(t *testing.T) {
		err := m.From()
		require.NoError(t, err, "from should not error")

		t.Logf("container: '%s'", m.b.Container)
		require.NotEmpty(t, m.b.Container, "container name should not be empty")
	})

	t.Run("mount", func(t *testing.T) {
		err := m.Mount()
		require.NoError(t, err, "mount should not error")
		t.Logf("mount-point: '%s'", m.mountPoint)
		require.NotEmpty(t, m.mountPoint, "mount point should be populated")

		stat, err := os.Stat(m.mountPoint)
		require.NoError(t, err, "mount point directory exists")
		require.True(t, stat.IsDir(), "mount point is a directory")
	})

	t.Run("add", func(t *testing.T) {
		cwd, err := os.Getwd()
		require.NoError(t, err, "should not error on getting cwd")

		filePath := path.Join(cwd, "../../README.md")
		t.Logf("Copying over '%s' file", filePath)
		err = m.Add(filePath, "/tmp")
		require.NoError(t, err, "should not error on copying file")
	})

	t.Run("unmount", func(t *testing.T) {
		err := m.Unmount()
		require.NoError(t, err, "unmount should not error")

		_, err = os.Stat(m.mountPoint)
		require.Error(t, err, "should error after unmount")
		require.True(t, os.IsNotExist(err), "error should be not-exists")
	})

	t.Run("delete", func(t *testing.T) {
		err := m.Delete()
		require.NoError(t, err, "delete should not error")
	})
}
