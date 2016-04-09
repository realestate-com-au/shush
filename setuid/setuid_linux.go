// +build linux

package setuid

import (
	"github.com/opencontainers/runc/libcontainer/system"
)

// Setuid sets the uid of the calling thread to the specified uid.
func Setuid(uid int) (err error) {
	if e1 := system.Setuid(uid); e1 != nil {
		return e1
	}
	return
}

// Setgid sets the gid of the calling thread to the specified gid.
func Setgid(gid int) (err error) {
	if e1 := system.Setgid(gid); e1 != nil {
		return e1
	}
	return
}
