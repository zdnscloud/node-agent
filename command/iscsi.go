package command

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	iscsiBinary     = "iscsiadm"
	multipathBinary = "multipath"
	lsscsiBinary    = "lsscsi"
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

func DiscoverTarget(ip, target string) (string, error) {
	opts := []string{
		"-m", "discovery",
		"-t", "sendtargets",
		"-p", ip,
	}
	output, err := execute(iscsiBinary, opts)
	if err != nil {
		return output, err
	}
	if strings.Contains(output, "Could not") {
		return output, fmt.Errorf("Cannot discover target: %s", output)
	}
	if !strings.Contains(output, target) {
		return output, fmt.Errorf("Cannot find target %s in discovered targets %s", target, output)
	}
	return output, nil
}

func DeleteDiscoveredTarget(ip, target string) (string, error) {
	opts := []string{
		"-m", "node",
		"-o", "delete",
		"-p", ip,
		"-T", target,
	}
	return execute(iscsiBinary, opts)
}

func IsTargetDiscovered(ip, target string) (string, bool) {
	opts := []string{
		"-m", "node",
		"-T", target,
		"-p", ip,
	}
	output, err := execute(iscsiBinary, opts)
	if err != nil {
		return output, false
	}
	return output, true
}

func IscsiChap(ip, target, username, password string) (string, error) {
	opts := []string{
		"-m", "node",
		"-T", target,
		"-p", ip,
		"-o", "update",
	}
	chapOpts := []string{
		"--name", "node.session.auth.authmethod",
		"--value=CHAP",
	}
	userOpts := []string{
		"--name", "node.session.auth.username",
		"--value=" + username,
	}
	passOpts := []string{
		"--name", "node.session.auth.password",
		"--value=" + password,
	}
	output, err := execute(iscsiBinary, append(opts, chapOpts...))
	if err != nil {
		return output, err
	}
	output, err = execute(iscsiBinary, append(opts, userOpts...))
	if err != nil {
		return output, err
	}
	output, err = execute(iscsiBinary, append(opts, passOpts...))
	if err != nil {
		return output, err
	}
	return output, nil
}

func LoginTarget(ip, target string) (string, error) {
	opts := []string{
		"-m", "node",
		"-T", target,
		"-p", ip,
		"--login",
	}
	return execute(iscsiBinary, opts)
}

func GetIscsiDevices(ip, target string) (map[string][]string, error) {
	var err error
	devs := make(map[string][]string)
	for i := 0; i < DeviceWaitRetryCounts; i++ {
		devs, err = findScsiDevice(ip, target)
		if err == nil {
			break
		}
		time.Sleep(DeviceWaitRetryInterval)
	}
	return devs, err
}

func findScsiDevice(ip, target string) (map[string][]string, error) {
	devs := make(map[string][]string)
	opts := []string{
		"-m", "session",
		"-P", "3",
	}
	output, err := execute(iscsiBinary, opts)
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	targetLine := "Target: " + target
	ipLine := " " + ip + ":"
	lunLine := "Lun: "
	diskPrefix := "Attached scsi disk"
	stateLine := "State:"

	inTarget := false
	inIP := false
	inLun := false
	var lun string
	for scanner.Scan() {
		if inTarget && (strings.HasPrefix(scanner.Text(), "Target: ")) {
			break
		}
		if !inTarget &&
			(strings.Contains(scanner.Text(), targetLine+" ") ||
				strings.HasSuffix(scanner.Text(), targetLine)) {
			inTarget = true
			continue
		}
		if inTarget && strings.Contains(scanner.Text(), ipLine) {
			inIP = true
			continue
		}
		if inIP && strings.Contains(scanner.Text(), lunLine) {
			lines := strings.Fields(scanner.Text())
			lun = lines[len(lines)-1]
			if _, ok := devs[lun]; !ok {
				devs[lun] = make([]string, 0)
			}
			inLun = true
			continue
		}
		// The line we need
		if inLun {
			line := scanner.Text()
			if !strings.Contains(line, diskPrefix) {
				return nil, fmt.Errorf("Invalid output format, cannot find disk in: %s\n %s", line, output)
			}
			line = strings.TrimSpace(strings.Split(line, stateLine)[0])
			line = strings.TrimPrefix(line, diskPrefix)
			dev := strings.TrimSpace(line)
			devs[lun] = append(devs[lun], dev)
			//devs[lun] = append(devs[lun], "/dev/"+dev)
			//break
			inLun = false
			lun = ""
		}
	}
	return devs, nil
}

func LogoutTarget(ip, target string) (string, error) {
	opts := []string{
		"-m", "node",
		"-T", target,
		"--logout",
	}
	if ip != "" {
		opts = append(opts, "-p", ip)
	}
	return execute(iscsiBinary, opts)
}

func CleanupScsiNodes(target string) (string, error) {
	for _, dir := range ScsiNodesDirs {
		if _, err := execute("ls", []string{dir}); err != nil {
			continue
		}
		targetDir := filepath.Join(dir, target)
		if _, err := execute("ls", []string{targetDir}); err != nil {
			continue
		}
		output, err := execute("find", []string{targetDir})
		if err != nil {
			return output, fmt.Errorf("Failed to search SCSI directory %v: %v", targetDir, err)
		}
		scanner := bufio.NewScanner(strings.NewReader(output))
		for scanner.Scan() {
			file := scanner.Text()
			output, err := execute("stat", []string{file})
			if err != nil {
				return output, fmt.Errorf("Failed to check SCSI node file %v: %v", file, err)
			}
			if strings.Contains(output, "regular empty file") {
				if output, err := execute("rm", []string{file}); err != nil {
					return output, fmt.Errorf("Failed to cleanup empty SCSI node file %v: %v", file, err)
				}
				execute("rmdir", []string{filepath.Dir(file)})
			}
		}
	}
	return "", nil
}

func GetIscsiMultipath(devs []string) (string, error) {
	var path string
	for _, dev := range devs {
		_path, err := getBlockMultipath(dev)
		if err != nil {
			return "", err
		}
		if path == "" {
			path = _path
		} else {
			if path != _path {
				return "", errors.New("devs have defferent multipath wwid")
			}
		}
	}
	if path == "" {
		return "", errors.New("can not get multipath for those devs")
	}
	return path, nil
}

func getBlockMultipath(dev string) (string, error) {
	opts := []string{
		"-v3",
	}
	output, err := execute(multipathBinary, opts)
	if err != nil {
		return "", err
	}
	scanner := bufio.NewScanner(strings.NewReader(output))
	pathList := " paths list "
	var inList bool
	for scanner.Scan() {
		if !inList &&
			(strings.Contains(scanner.Text(), pathList) ||
				strings.HasSuffix(scanner.Text(), pathList)) {
			inList = true
			continue
		}
		if inList && strings.Contains(scanner.Text(), dev) {
			return strings.Fields(scanner.Text())[0], nil
		}
	}
	return "", errors.New("can not get multipath for device")
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
				fmt.Printf("1111", err)
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
