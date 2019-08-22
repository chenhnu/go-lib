package ptdb

import (
	"strconv"
	"strings"
)

func isHost(host string) bool {
	res := strings.Split(host, ".")
	if len(res) != 4 {
		return false
	}
	for _, s := range res {
		num, e := strconv.Atoi(s)
		if e != nil || num < 0 || num > 255 {
			return false
		}
	}
	return true
}
