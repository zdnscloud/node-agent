package command

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	pb "github.com/zdnscloud/node-agent/proto"
)

var supportType = []string{"disk"}

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
	MAJMIN     string     `json:"maj:min"`
}

type children struct {
	NAME       string `json:"name"`
	TYPE       string `json:"type"`
	SIZE       int64  `json:"size"`
	PART       string `json:"part"`
	PKNAME     string `json:"pkname"`
	MOUNTPOINT string `json:"mountpoint"`
	FSTYPE     string `json:"fstype"`
	MAJMIN     string `json:"maj:min"`
}

func GetDisksInfo(disk string) (map[string]*pb.Diskinfo, error) {
	infos := make(map[string]*pb.Diskinfo)
	devs, err := getAllDevs()
	if err != nil {
		return infos, err
	}
	for _, d := range devs {
		if len(d) == 0 {
			continue
		}
		dev := "/dev/" + d
		gpt, err := checkGpt(dev)
		if err != nil {
			return infos, err
		}
		if gpt {
			continue
		}
		conf := gen(dev)
		if conf == nil {
			continue
		}
		infos[dev] = &pb.Diskinfo{
			Diskinfo: conf,
		}
	}
	return infos, nil
}

func getAllDevs() ([]string, error) {
	args := []string{
		"--all",
		"--noheadings",
		"--list",
		"--output",
		"KNAME",
	}
	out, err := exec.Command("lsblk", args...).Output()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("get all devices failed, %s ", err.Error()))
	}
	return strings.Split(string(out), "\n"), nil
}

func gen(disk string) map[string]string {
	args := []string{
		"--bytes",
		"-J",
		"--output",
		"NAME,TYPE,SIZE,PKNAME,FSTYPE,MOUNTPOINT,MAJ:MIN",
	}
	args = append(args, disk)
	out, err := exec.Command("/bin/lsblk", args...).Output()
	if err != nil {
		fmt.Sprintf("get dev %s info failed, %s ", disk, err.Error())
		return nil
	}
	var res Info
	json.Unmarshal(out, &res)
	for _, d := range res.Blockdevices {
		if !checkMaj(d) || !checkType(d) {
			continue
		}
		return conf(d)
	}
	return nil
}

func conf(info bd) map[string]string {
	cfg := make(map[string]string)
	mountpoint := "false"
	fsexist := "false"
	parted := "false"
	if len(info.CHILDREN) > 0 {
		parted = "true"
		for _, children := range info.CHILDREN {
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
	if info.FSTYPE != "" {
		fsexist = "true"
	}
	if info.MOUNTPOINT != "" {
		mountpoint = "true"
	}
	cfg["Size"] = strconv.FormatInt(info.SIZE, 10)
	cfg["Parted"] = parted
	cfg["Filesystem"] = fsexist
	cfg["Mountpoint"] = mountpoint
	return cfg
}

func checkMaj(info bd) bool {
	maj := strings.Split(info.MAJMIN, ":")[0]
	m, _ := strconv.Atoi(maj)
	if m == 8 {
		return true
	}
	if m >= 240 && m <= 254 {
		return true
	}
	return false
}

func checkType(info bd) bool {
	for _, t := range supportType {
		if t == info.TYPE {
			return true
		}
	}
	return false
}

func checkGpt(disk string) (bool, error) {
	args := []string{
		"info",
		"--query=property",
	}
	args = append(args, disk)
	out, err := exec.Command("udevadm", args...).Output()
	if err != nil {
		return true, errors.New(fmt.Sprintf("get dev %s gpt info failed, %s ", disk, err.Error()))
	}
	outputs := strings.Split(string(out), "\n")
	for _, l := range outputs {
		if !strings.Contains(l, "ID_PART_TABLE_TYPE") {
			continue
		}
		tableType := strings.Split(l, "=")[0]
		if tableType == "gpt" {
			return true, nil
		}
	}
	return false, nil
}
