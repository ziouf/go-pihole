package utils

import (
	"os/user"
	"strconv"
)

func FindUidAndGidFromName(name string) (uint32, uint32) {
	user, _ := user.Lookup("dnsmasq")
	uid, _ := strconv.Atoi(user.Uid)
	gid, _ := strconv.Atoi(user.Gid)
	return uint32(uid), uint32(gid)
}
