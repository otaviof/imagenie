package imagenie

import (
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Imagenie primary application representation.
type Imagenie struct {
	cfg       *Config  // global application configuration
	sourceMgr *Manager // source container manager
	TargetMgr *Manager // target container manager
}

// CopyPaths source directory as key, and destination as value.
type CopyPaths map[string]string

// Copy loop over paths in order to copy it from source-image into destination container.
func (i *Imagenie) Copy(paths CopyPaths) error {
	log.Infof("Entries to copy '%d'", len(paths))

	// mounting source container manager
	err := i.sourceMgr.Mount()
	if err != nil {
		return err
	}

	// looping over paths to be copied from one image to another
	for src, dst := range paths {
		if dst == "" {
			dst = src
		}

		// using mount point to compose source directory location
		src = path.Join(i.sourceMgr.mountPoint, src)
		log.Infof("Copying '%s' to '%s'...", src, dst)

		if err = i.TargetMgr.Add(src, dst); err != nil {
			return err
		}
	}
	return nil
}

// Labels inspect labels in source container, and set them in target container.
func (i *Imagenie) Labels() {
	for k, v := range i.sourceMgr.Labels() {
		log.Infof("Setting label: '%s=%v'", k, v)
		i.TargetMgr.SetLabel(k, v)
	}
}

// RunAll run all commands in informed slice, each entry is splitted again by spaces and executed on
// target container.
func (i *Imagenie) RunAll(commands []string) error {
	for _, command := range commands {
		commandSlice := strings.Split(command, " ")
		if err := i.TargetMgr.Run(commandSlice); err != nil {
			return err
		}
	}
	return nil
}

// CleanUp unmount source container, remove source and target working containers.
func (i *Imagenie) CleanUp() error {
	if err := i.sourceMgr.Unmount(); err != nil {
		return err
	}
	if err := i.sourceMgr.Delete(); err != nil {
		return err
	}
	return i.TargetMgr.Delete()
}

// bootstrap container managers and mount source-image.
func (i *Imagenie) bootstrap() error {
	var err error
	i.sourceMgr, err = NewManager(i.cfg.FromImage, "")
	if err != nil {
		return err
	}
	i.TargetMgr, err = NewManager(i.cfg.BaseImage, i.cfg.TargetImage)
	if err != nil {
		return err
	}
	return i.TargetMgr.From()
}

// NewImagenie instantiate application.
func NewImagenie(cfg *Config) (*Imagenie, error) {
	i := &Imagenie{cfg: cfg}
	return i, i.bootstrap()
}
