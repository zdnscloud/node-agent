package main

import (
        "context"
        "os"
        "os/exec"
        "path"
        "strings"
        "time"

        docker "github.com/docker/docker/client"
        "github.com/zdnscloud/cement/log"
)

const (
        DockerAPIVersion     = "1.24"
        kubeletContainerName = "kubelet"
        Sock                 = "unix:///var/run/docker.sock"
        dockerLogDir         = "/var/lib/docker/containers"
        kubeletLogDir        = "/var/lib/kubelet/pods"
        volumeDir            = "volumes/kubernetes.io~csi"
        flag                 = "errors similar to this. Turn up verbosity to see them."
        interval             = 5
)

func main() {
        log.InitLogger(log.Debug)

        kid, err := getKubeletID()
        if err != nil {
                log.Fatalf("get kubelet container id failed. Err: %s", err.Error())
        }
        log.Debugf("get kubelet container id %s, start watch", kid)
        logFile := path.Join(dockerLogDir, kid, kid+"-json.log")

        for {
                orphanedID, err := getOrphanedID(logFile)
                if err != nil {
                        log.Errorf("get orphaned container id failed. Err: %s", err.Error())
                }
                if orphanedID != "" {
                        dir := path.Join(kubeletLogDir, orphanedID, volumeDir)
                        if _, err := os.Stat(dir); err != nil {
                                continue
                        }
                        log.Debugf("orphaned pods %s found, clean it's volume now", orphanedID)
                        if err := removeContents(dir); err != nil {
                                log.Errorf("clean failed. Err: %s", err.Error())
                        }
                }
                time.Sleep(interval * time.Second)
        }
}

func getKubeletID() (string, error) {
        client, err := docker.NewClient(Sock, DockerAPIVersion, nil, nil)
        if err != nil {
                return "", err
        }
        r, err := client.ContainerInspect(context.Background(), kubeletContainerName)
        if err != nil {
                return "", err
        }
        return r.ID, nil
}

func getOrphanedID(log string) (string, error) {
        var id string
        info, err := exec.Command("tail", "-n", "20", log).Output()
        if err != nil {
                return id, err
        }
        lines := strings.Split(string(info), "\n")
        for _, line := range lines {
                if strings.Contains(line, flag) {
                        return strings.Replace(strings.Fields(line)[6], "\\\"", "", -1), nil
                }
        }
        return id, nil
}

func removeContents(dir string) error {
        d, err := os.Open(dir)
        if err != nil {
                return err
        }
        defer d.Close()
        names, err := d.Readdirnames(-1)
        if err != nil {
                return err
        }
        for _, name := range names {
                err = os.RemoveAll(path.Join(dir, name))
                if err != nil {
                        return err
                }
        }
        return nil
}
