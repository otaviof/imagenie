package imagenie

import (
	"github.com/containers/buildah"
	"github.com/containers/storage/pkg/unshare"
)

// ReInit executes buildah steps to initialize unshare.
func ReInit() {
	if buildah.InitReexec() {
		return
	}
	unshare.MaybeReexecUsingUserNamespace(false)
}
