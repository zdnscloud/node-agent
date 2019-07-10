package server

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"

	pb "github.com/zdnscloud/node-agent/proto"
)

func getDirectorySize(path string) (map[string]int64, error) {
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

func getMountpointsSize(paths []string) (map[string]*pb.Sizes, error) {
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

type Info struct {
	Blockdevices []bd `json:"blockdevices"`
}

type bd struct {
	NAME       string     `json:"name"`
	TYPE       string     `json:"type"`
	SIZE       int64      `json:"size"`
	PKNAME     string     `json:"pkname"`
	FSTYPE     string     `json:"fstype"`
	MOUNTPOINT string     `json:"mountpoint"`
	CHILDREN   []children `json:"children"`
}

type children struct {
	NAME       string `json:"name"`
	TYPE       string `json:"type"`
	SIZE       int64  `json:"size"`
	PART       string `json:"part"`
	PKNAME     string `json:"pkname"`
	MOUNTPOINT string `json:"mountpoint"`
	FSTYPE     string `json:"fstype"`
}

func getDisksInfo(disk string) (map[string]*pb.Diskinfo, error) {
	infos := make(map[string]*pb.Diskinfo)
	args := []string{
		"--bytes",
		"-J",
		"--output",
		"NAME,TYPE,SIZE,PKNAME,FSTYPE,MOUNTPOINT",
	}
	if len(disk) > 0 {
		args = append(args, disk)
	}
	out, err := exec.Command("/bin/lsblk", args...).Output()
	if err != nil {
		return infos, err
	}
	var res Info
	json.Unmarshal(out, &res)
	for _, d := range res.Blockdevices {
		if d.TYPE != "disk" {
			continue
		}
		mountpoint := "false"
		fsexist := "false"
		parted := "false"
		if len(d.CHILDREN) > 0 {
			parted = "true"
			for _, children := range d.CHILDREN {
				if children.FSTYPE == "" {
					continue
				}
				fsexist = "true"
				if children.MOUNTPOINT == "" {
					continue
				}
				mountpoint = "true"
			}
		}
		info := make(map[string]string)
		info["Size"] = strconv.FormatInt(d.SIZE, 10)
		info["Parted"] = parted
		info["Filesystem"] = fsexist
		info["Mountpoint"] = mountpoint
		name := "/dev/" + d.NAME
		infos[name] = &pb.Diskinfo{
			Diskinfo: info,
		}
	}
	return infos, nil
}
