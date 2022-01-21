package zfs_go

import (
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"strings"
)

const ZFS_BINARY string = "/sbin/zfs"

type ZfsEntity struct {
	XMLName    xml.Name `xml:"zfsentity"`
	Name       string   `xml:"name"`
	Used       string   `xml:"used"`
	Avail      string   `xml:"avail"`
	Refer      string   `xml:"refer"`
	MountPoint string   `xml:"mountpoint"`
}

func ZfsListAll() ([]ZfsEntity, error) {
	var (
		out []byte
		err error
		res []ZfsEntity
	)
	cmd := exec.Command("zfs", "list", "-H")
	if out, err = cmd.CombinedOutput(); err != nil {
		log.Println(err)
	}
	datasets := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")

	for _, v := range datasets {
		v_split := strings.Split(v, "\t")
		res = append(
			res, ZfsEntity{
				Name:       v_split[0],
				Used:       v_split[1],
				Avail:      v_split[2],
				Refer:      v_split[3],
				MountPoint: v_split[4]})
	}
	return res, err
}

func ZfsGetLastSnapshot(dataset string) (string, error) {
	var (
		out []byte
		err error
		res string
	)
	cmd := exec.Command(ZFS_BINARY, "list", "-Ho", "name", "-t", "snapshot", "-r", dataset)
	if out, err = cmd.CombinedOutput(); err != nil {
		log.Println(err)
	}
	snapshots := strings.Split(strings.TrimSuffix(string(out), "\n"), "\n")
	if res = snapshots[len(snapshots)-1]; res == "" {
		err = errors.New(fmt.Sprintf("there is no any snapshot in %s", dataset))
	}
	return res, err
}

func ZfsGetSnapshotInfo(dataset string) (map[string]string, error) {
	var (
		out []byte
		err error
		res map[string]string
	)
	res = make(map[string]string)
	cmd := exec.Command(ZFS_BINARY, "get", "-Hpo", "value", "origin,written", dataset)
	if out, err = cmd.CombinedOutput(); err != nil {
		log.Println(err)
	}
	out_split := strings.Split(string(out), "\n")
	res["origin"] = out_split[0]
	res["written"] = out_split[1]
	return res, err
}

func ZfsCreateSnapshot(snapsource string, snapname string) (map[string]string, error) {
	var (
		// out []byte
		err error
		res map[string]string
	)
	res = make(map[string]string)
	if snapsource != "" && snapname != "" {
		cmd := exec.Command(ZFS_BINARY, "snapshot", snapsource+"@"+snapname)
		if _, err = cmd.CombinedOutput(); err != nil {
			log.Println(err)
		}

		// res =

	} else {
		err = errors.New("missing snapshot source or snapshot name.")
	}
	return res, err
}
