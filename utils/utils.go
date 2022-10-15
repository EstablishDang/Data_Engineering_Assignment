package utils

import (
	"log"
	"os/exec"
	"strconv"
	"time"
)

const (
	PATH_NETWORK_INTERFACE        = "/etc/network/interfaces"
	PATH_NETWORK_INTERFACE_BACKUP = "/etc/network/interfaces_backup"
	PATH_PVE_DATA_CONFIG          = "/etc/pve/pve_nw_config/"
)

func StringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Printf("StringToInt error: %v", err)
	}
	return i
}

func IntToString(i int) string {
	return strconv.FormatInt(int64(i), 10)
}

func ExecCmd(cmd string) (string, error) {
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		log.Printf("[ERROR]_Exec command %s fail, err: %v", cmd, err)
		return "", err
	}
	return string(out), nil
}

// Convert Unix seconds to ISODate, UTC time
func UnixSecondsToTimestampRFC3339(secs int64) string {
	t := time.Unix(secs, 0)
	loc, _ := time.LoadLocation("UTC")
	t2 := t.In(loc)
	return t2.Format("2006-01-02T15:04:05.000Z")
}
