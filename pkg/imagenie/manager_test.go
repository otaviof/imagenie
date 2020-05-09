package imagenie

import (
	"os"
	"path"
	"testing"

	is "github.com/containers/image/v5/storage"
	"github.com/stretchr/testify/require"
)

const (
	fromImage   = "alpine:latest"
	targetImage = "alpine-test:latest"
)

func init() {
	os.Setenv(LogLevelEnv, "99")
	SetLogLevel()
}

func Test(t *testing.T) {
	store, err := getStore()
	t.Logf("err: '%#v'", err)

	imageRef, err := is.Transport.ParseStoreReference(store, targetImage)
	t.Logf("image-ref: '%#v'", imageRef)
	t.Logf("err: '%#v'", err)
}

func TestManager_NewManager(t *testing.T) {
	m, err := NewManager(fromImage, targetImage)
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

	t.Run("commit", func(t *testing.T) {
		err := m.Commit()
		require.NoError(t, err, "should not error on committing")
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
