package command

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	pb "github.com/zdnscloud/node-agent/proto"
)

func GetDisksInfo() (map[string]*pb.Diskinfo, error) {
	infos, err := getDeviceInfo()
	if err != nil {
		return nil, err
	}
	res := make(map[string]*pb.Diskinfo)
	for _, info := range infos {
		diskInfo := make(map[string]string)
		name := "/dev/" + info.name
		res[name] = &pb.Diskinfo{
			Diskinfo: diskInfo,
		}
		diskInfo["Size"] = info.size
		if info.fstype != "" {
			diskInfo["Filesystem"] = "true"
		}
		if info.mountpoint != "" {
			diskInfo["Mountpoint"] = "true"
		}

		if info.typ == "part" {
			output, err := execute("udevadm", []string{"info", "--query=property", name})
			if err != nil {
				fmt.Println(err)
				continue
			}
			scanner := bufio.NewScanner(strings.NewReader(output))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.HasPrefix(line, "DEVPATH=") {
					tmp := strings.Split(line, "/")
					Parted := tmp[len(tmp)-2]
					_name := "/dev/" + Parted
					v, ok := res[_name]
					if ok {
						v.Diskinfo["Parted"] = "true"
					}
				}
			}
		}
		if info.typ != "disk" {
			delete(res, name)
			continue
		}
		if info.maj != 8 && (info.maj < 240 || info.maj > 254) {
			delete(res, name)
			continue
		}
		gpt, err := checkGpt(name)
		if err != nil {
			return nil, err
		}
		if gpt {
			continue
		}
	}
	return res, nil
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

func getDeviceInfo() ([]Info, error) {
	/*
		# lsblk -a --bytes --noheadings -P -o NAME,TYPE,SIZE,FSTYPE,MOUNTPOINT,MAJ:MIN
		NAME="loop1" TYPE="loop" SIZE="95805440" FSTYPE="squashfs" MOUNTPOINT="" MAJ:MIN="7:1"
		NAME="vdd" TYPE="disk" SIZE="21474836480" FSTYPE="" MOUNTPOINT="" MAJ:MIN="252:48"
		NAME="loop6" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:6"
		NAME="vdb" TYPE="disk" SIZE="21474836480" FSTYPE="" MOUNTPOINT="" MAJ:MIN="252:16"
		NAME="loop4" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:4"
		NAME="sr0" TYPE="rom" SIZE="1073741312" FSTYPE="" MOUNTPOINT="" MAJ:MIN="11:0"
		NAME="loop2" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:2"
		NAME="loop0" TYPE="loop" SIZE="95748096" FSTYPE="squashfs" MOUNTPOINT="" MAJ:MIN="7:0"
		NAME="loop7" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:7"
		NAME="vdc" TYPE="disk" SIZE="21474836480" FSTYPE="" MOUNTPOINT="" MAJ:MIN="252:32"
		NAME="loop5" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:5"
		NAME="vda" TYPE="disk" SIZE="53687234560" FSTYPE="" MOUNTPOINT="" MAJ:MIN="252:0"
		NAME="vda2" TYPE="part" SIZE="53683945472" FSTYPE="ext4" MOUNTPOINT="/dev/termination-log" MAJ:MIN="252:2"
		NAME="vda1" TYPE="part" SIZE="1048576" FSTYPE="" MOUNTPOINT="" MAJ:MIN="252:1"
		NAME="loop3" TYPE="loop" SIZE="" FSTYPE="" MOUNTPOINT="" MAJ:MIN="7:3"
	*/
	opts := []string{
		"-a",
		"--bytes",
		"--noheadings",
		"-P",
		"-o", "NAME,TYPE,SIZE,KNAME,FSTYPE,MOUNTPOINT,MAJ:MIN",
	}
	output, err := execute("lsblk", opts)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	var infos []Info
	for scanner.Scan() {
		var info Info
		tmp := strings.Fields(scanner.Text())
		for i := 0; i < len(tmp); i++ {
			h := strings.Split(tmp[i], "=")
			switch h[0] {
			case "NAME":
				info.name = strings.Trim(h[1], "\"")
			case "TYPE":
				info.typ = strings.Trim(h[1], "\"")
			case "SIZE":
				info.size = strings.Trim(h[1], "\"")
			case "FSTYPE":
				info.fstype = strings.Trim(h[1], "\"")
			case "MOUNTPOINT":
				info.mountpoint = strings.Trim(h[1], "\"")
			case "MAJ:MIN":
				t := strings.Split(strings.Trim(h[1], "\""), ":")
				j, _ := strconv.Atoi(t[0])
				i, _ := strconv.Atoi(t[1])
				info.maj = j
				info.min = i
			}
		}
		infos = append(infos, info)
	}
	return infos, nil
}

type Info struct {
	name       string
	typ        string
	size       string
	fstype     string
	mountpoint string
	maj        int
	min        int
}
