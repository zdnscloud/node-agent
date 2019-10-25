package command

import (
	"os/exec"
	"strconv"
	"strings"

	pb "github.com/zdnscloud/node-agent/proto"
)

func GetMountpointsSize(paths []string) (map[string]*pb.Sizes, error) {
	infos := make(map[string]*pb.Sizes)
	out, _ := exec.Command("df", paths...).Output()
	outputs := strings.Split(string(out), "\n")
	for i := 1; i < len(outputs); i++ {
		if !strings.Contains(outputs[i], "%") {
			continue
		}
		line := strings.Fields(outputs[i])
		num := len(line)
		tsize, err := strconv.ParseInt(line[num-5], 10, 64)
		if err != nil {
			return infos, err
		}
		usize, err := strconv.ParseInt(line[num-4], 10, 64)
		if err != nil {
			return infos, err
		}
		fsize, err := strconv.ParseInt(line[num-3], 10, 64)
		if err != nil {
			return infos, err
		}
		infos[line[num-1]] = &pb.Sizes{
			Tsize: tsize,
			Usize: usize,
			Fsize: fsize,
		}
	}
	return infos, nil
}
