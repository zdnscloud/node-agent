package main

import (
	"bufio"
	"bytes"
	"fmt"
	pb "github.com/zdnscloud/node-agent/proto"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Info struct {
	name       string
	typ        string
	size       string
	fstype     string
	mountpoint string
	maj        int
	min        int
}

const (
	iscsiBinary     = "iscsiadm"
	multipathBinary = "multipath"
	lsscsiBinary    = "lsscsi"
	dmsetupBinary   = "dmsetup"
	lsblkBinary     = "lsblk"
	BUFFERSIZE      = 1000
)

var (
	DeviceWaitRetryCounts   = 5
	DeviceWaitRetryInterval = 1 * time.Second

	ScsiNodesDirs = []string{
		"/etc/iscsi/nodes/",
		"/var/lib/iscsi/nodes/",
	}
	cmdTimeout = time.Minute // one minute by default
)

func main() {
	res, err := gen()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(res)
}

func execute(binary string, args []string) (string, error) {
	return executeWithTimeout(cmdTimeout, binary, args)
}

func executeWithTimeout(timeout time.Duration, binary string, args []string) (string, error) {
	var err error
	cmd := exec.Command(binary, args...)
	done := make(chan struct{})

	var output, stderr bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &stderr

	go func() {
		err = cmd.Run()
		done <- struct{}{}
	}()

	select {
	case <-done:
	case <-time.After(timeout):
		if cmd.Process != nil {
			if err := cmd.Process.Kill(); err != nil {
				fmt.Sprintf("Problem killing process pid=%v: %s", cmd.Process.Pid, err)
			}

		}
		return "", fmt.Errorf("Timeout executing: %v %v, output %s, stderr, %s, error %v",
			binary, args, output.String(), stderr.String(), err)
	}

	if err != nil {
		return "", fmt.Errorf("Failed to execute: %v %v, output %s, stderr, %s, error %v",
			binary, args, output.String(), stderr.String(), err)
	}
	return output.String(), nil
}
func getDeviceInfo() ([]Info, error) {
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

func gen() (map[string]*pb.Diskinfo, error) {
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
					v, ok := res["/dev/"+Parted]
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
	}
	return res, nil
}
