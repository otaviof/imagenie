package imagenie

import (
	"context"
	"fmt"
	"os"

	"github.com/containers/buildah"
	"github.com/containers/buildah/imagebuildah"
	is "github.com/containers/image/v5/storage"
	"github.com/containers/image/v5/types"
	"github.com/containers/storage"
	"github.com/containers/storage/pkg/unshare"
	"github.com/opencontainers/go-digest"
	log "github.com/sirupsen/logrus"
)

// Manager buildaah based container manager.
type Manager struct {
	fromImage   string           // source container image
	targetImage string           // target image
	ctx         context.Context  // context instance
	store       storage.Store    // container storage instance
	b           *buildah.Builder // builder instance
	mountPoint  string           // source container root mount path
}

var systemContext = &types.SystemContext{}

// From save container image, therefore subsequent changes can take place.
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

// SetLabel set container label.
func (m *Manager) SetLabel(k, v string) {
	m.b.SetLabel(k, v)
}

// Labels return a map with container labels.
func (m *Manager) Labels() map[string]string {
	return m.b.Labels()
}

// SetEntrypoint set informed string slice as entrypoint.
func (m *Manager) SetEntrypoint(entrypoint []string) {
	m.b.SetEntrypoint(entrypoint)
}

// SetCMD set informed string slice as cmd.
func (m *Manager) SetCMD(cmd []string) {
	m.b.SetCmd(cmd)
}

// Run arbitrary commaand on container.
func (m *Manager) Run(command []string) error {
	opts := buildah.RunOptions{}
	return m.b.Run(command, opts)
}

// Commit execute commit by creating a image out of container in use.
func (m *Manager) Commit() error {
	targetRef, err := is.Transport.ParseStoreReference(m.store, m.targetImage)
	if err != nil {
		return err
	}

	opts := buildah.CommitOptions{
		Squash:        true,
		SystemContext: systemContext,
		Compression:   imagebuildah.Gzip,
	}
	id, ref, _, err := m.b.Commit(m.ctx, targetRef, opts)
	if err != nil {
		return err
	}

	log.Infof("Image-id: '%s'", id)
	log.Infof("Name: '%s'", ref.Name())
	log.Infof("Digest: '%s'", ref.Digest().String())
	return nil
}

// Delete delete current container.
func (m *Manager) Delete() error {
	return m.b.Delete()
}

// getStore get a storage.Store instance.
func getStore() (storage.Store, error) {
	opts, err := storage.DefaultStoreOptions(unshare.IsRootless(), unshare.GetRootlessUID())
	if err != nil {
		return nil, err
	}
	store, err := storage.GetStore(opts)
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
		Format:           buildah.OCIv1ImageManifest,
		FromImage:        m.fromImage,
		Isolation:        buildah.IsolationChroot,
		ReportWriter:     os.Stderr,
		SystemContext:    systemContext,
	}
}

// bootstrap instantiate builder with options.
func (m *Manager) bootstrap() error {
	var err error
	if m.store, err = getStore(); err != nil {
		return err
	}

	log.Debugf("Instantiating buildah.Builder with '%s' base image", m.fromImage)
	opts := m.builderOptions()
	m.b, err = buildah.NewBuilder(m.ctx, m.store, opts)
	return err
}

// NewManager instantiate and bootstrap manager.
func NewManager(fromImage, targetImage string) (*Manager, error) {
	ReInit()
	m := &Manager{fromImage: fromImage, targetImage: targetImage, ctx: context.TODO()}
	return m, m.bootstrap()
}
