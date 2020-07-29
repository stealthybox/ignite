package cleanup

import (
	"os"

	api "github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/util"
)

// DeactivateSnapshot deactivates the snapshot by removing it with dmsetup
func DeactivateSnapshot(vm *api.VM) error {
	dmArgs := []string{
		"remove",
		util.NewPrefixer().Prefix(vm.GetUID()),
	}

	// If the base device is visible in "dmsetup", we should remove it
	// The device itself is not forwarded to docker, so we can't query its path
	// TODO: Improve this detection
	baseDev := util.NewPrefixer().Prefix(vm.GetUID(), "base")
	if _, err := util.ExecuteCommand("dmsetup", "info", baseDev); err == nil {
		dmArgs = append(dmArgs, baseDev)
	}

	_, err := util.ExecuteCommand("dmsetup", dmArgs...)
	if err != nil {
		return err
	}

	// VM's from previous versions of ignite may not have a snapshot-dev symlink
	// Lstat it first before attempting to remove it
	if _, err = os.Lstat(vm.SnapshotDevLink()); err == nil {
		if err = os.Remove(vm.SnapshotDevLink()); err != nil {
			return err
		}
	}

	return nil
}
