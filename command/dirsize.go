package command

import (
	"os/exec"
	"strconv"
	"strings"
)

func GetDirectorySize(path string) (map[string]int64, error) {
	infos := make(map[string]int64)
	out, _ := exec.Command("du", "-d1", path).Output()
	outputs := strings.Split(string(out), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, "/") {
			continue
		}
		line := strings.Fields(l)
		size, err := strconv.ParseInt(line[0], 10, 64)
		if err != nil {
			return infos, err
		}
		infos[line[1]] = size
	}
	delete(infos, path)
	return infos, nil
}
