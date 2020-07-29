package dmlegacy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	losetup "github.com/freddierice/go-losetup"
)

// loopDevice is a helper struct for handling loopback devices for devicemapper
type loopDevice struct {
	losetup.Device
}

func newLoopDev(file string, readOnly bool) (*loopDevice, error) {
	dev, err := losetup.Attach(file, 0, readOnly)
	if err != nil {
		return nil, fmt.Errorf("failed to setup loop device for %q: %v", file, err)
	}

	return &loopDevice{dev}, nil
}

func (ld *loopDevice) Size512K() (uint64, error) {
	data, err := ioutil.ReadFile(path.Join("/sys/class/block", path.Base(ld.Device.Path()), "size"))
	if err != nil {
		return 0, err
	}

	// Remove the trailing newline and parse to uint64
	return strconv.ParseUint(string(data[:len(data)-1]), 10, 64)
}

// dmsetup uses stdin to read multiline tables, this is a helper function for that
func runDMSetup(name string, table []byte) error {
	cmd := exec.Command(
		"dmsetup", "create",
		"--noudevsync", // we don't depend on udevd's /dev/mapper/* symlinks, we create our own
		name,
	)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if _, err := stdin.Write(table); err != nil {
		return err
	}

	if err := stdin.Close(); err != nil {
		return err
	}

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("command %q exited with %q: %w", cmd.Args, out, err)
	}

	return nil
}

// GetBlkDevPath returns the device path for a named device without the use of udevd's symlinks.
// This is useful for creating our own symlinks to track devices and pass to the sandbox container.
// The device path could be `/dev/<blkdevname>` (ex: `/dev/dm-0`)
// or `/dev/mapper/<name>` (ex: `/dev/mapper/ignite-47a6421c19b415ef`)
// depending on `dmsetup`'s udev-fallback related environment-variables and build-flags.
func GetBlkDevPath(name string) (string, error) {
	cmd := exec.Command(
		"dmsetup", "info",
		"--noudevsync", // we don't depend on udevd's /dev/mapper/* symlinks, we create our own
		"--columns", "--noheadings", "-o", "blkdevname",
		name,
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("command %q exited with %q: %w", cmd.Args, string(out), err)
	}
	blkdevname := strings.TrimSpace(string(out))

	// if dmsetup is compiled without udev-sync or the DM_DISABLE_UDEV env var is set,
	// `dmsetup info` will not return the correct blkdevname ("mapper/<name>") -- it still
	// returns "dm-<minor>" even though the path doesn't exist.
	// To work around this, we stat the returned blkdevname and try the fallback if it doesn't exist:
	blkDevPath := path.Join("/dev", blkdevname)
	if _, blkErr := os.Stat(blkDevPath); blkErr == nil {
		return blkDevPath, nil
	} else if !os.IsNotExist(blkErr) {
		return "", blkErr
	}

	fallbackDevPath := path.Join("/dev/mapper", name)
	if _, fallbackErr := os.Stat(fallbackDevPath); fallbackErr == nil {
		return fallbackDevPath, nil
	}

	return "", fmt.Errorf("Could not stat a valid block device path for %q or %q", blkDevPath, fallbackDevPath)
}
