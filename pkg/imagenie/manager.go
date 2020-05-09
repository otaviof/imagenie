package imagenie

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/buildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
	"github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
)

// Manager represents the container-manager, using buildah in the background.
type Manager struct {
	fromImage   string           // source container image
	targetImage string           // target image
	ctx         context.Context  // context instance
	store       storage.Store    // container storage instance
	b           *buildah.Builder // builder instance
	mountPoint  string           // source container root mount path
}

func (m *Manager) From() error {
	return m.b.Save()
}

// Mount the container and store the mount point path.
func (m *Manager) Mount() error {
	var err error
	m.mountPoint, err = m.b.Mount(m.b.MountLabel)
	if err != nil {
		return err
	}
	if m.mountPoint == "" {
		return fmt.Errorf("container-id '%s' is not mounted", m.b.ContainerID)
	}
	return nil
}

// Unmount container.
func (m *Manager) Unmount() error {
	return m.b.Unmount()
}

// Add local path to container.
func (m *Manager) Add(src, dst string) error {
	digester := digest.Canonical.Digester()
	opts := buildah.AddAndCopyOptions{
		Hasher: digester.Hash(),
	}
	log.Infof("Adding '%s' to '%s' (container '%s')", src, dst, m.b.Container)
	if err := m.b.Add(dst, false, opts, src); err != nil {
		return err
	}
	return m.b.Save()
}

func (m *Manager) SetLabel(k, v string) {
	m.b.SetLabel(k, v)
}

func (m *Manager) Labels() map[string]string {
	return m.b.Labels()
}

// Delete delete current container.
func (m *Manager) Delete() error {
	return m.b.Delete()
}

// getStore get a storage.Store instance.
func getStore() (storage.Store, error) {
	storeOpts, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	if err != nil {
		return nil, err
	}
	store, err := storage.GetStore(storeOpts)
	if err != nil {
		return nil, err
	}
	if store == nil {
		return nil, fmt.Errorf("unable to instantiate storage")
	}
	is.Transport.SetStore(store)
	return store, nil
}

// builderOptions returns the options buildah will use.
func (m *Manager) builderOptions() buildah.BuilderOptions {
	opts := &buildah.CommonBuildOptions{}
	return buildah.BuilderOptions{
		CommonBuildOpts:  opts,
		ConfigureNetwork: buildah.NetworkDefault,
		FromImage:        m.fromImage,
		Isolation:        buildah.IsolationChroot,
		ReportWriter:     os.Stderr,
		SystemContext:    &types.SystemContext{},
	}
}

// bootstrap instantiate builder with options.
func (m *Manager) bootstrap() error {
	store, err := getStore()
	if err != nil {
		return err
	}

	log.Debugf("Instantiating buildah.Builder with '%s' base image", m.fromImage)
	opts := m.builderOptions()
	m.b, err = buildah.NewBuilder(m.ctx, store, opts)
	return err
}

// NewManager instantiate and bootstrap manager.
func NewManager(fromImage, targetImage string) (*Manager, error) {
	if buildah.InitReexec() {
		return nil, nil
	}
	unshare.MaybeReexecUsingUserNamespace(false)

	m := &Manager{fromImage: fromImage, targetImage: targetImage, ctx: context.TODO()}
	return m, m.bootstrap()
}
