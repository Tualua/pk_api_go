package main

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
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

func apiCall(api string, command string, param map[string]string) ([]byte, error) {
	var (
		err          error
		response     *http.Response
		responseData []byte
		res          []byte
	)

	u, _ := url.Parse(api)
	q := u.Query()
	q.Set("action", command)
	if param != nil {

		for k := range param {
			q.Set(k, param[k])
		}
	}
	u.RawQuery = q.Encode()
	// u.Path = path.Join(u.Path, command)
	apiUrl := u.String()
	if response, err = http.Get(apiUrl); err != nil {
		res = []byte(err.Error())
		log.Println(err.Error())
	} else {
		if responseData, err = ioutil.ReadAll(response.Body); err != nil {
			fmt.Println(err.Error())
		} else {
			res = responseData
		}
	}
	return res, err
}

func ZfsListAll(apiZfs string) ([]ZfsEntity, error) {
	var (
		apiResponse []byte
		err         error
		res         []ZfsEntity
		jsonData    jsonResponseListAll
	)
	if apiResponse, err = apiCall(apiZfs, "listall", nil); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		res = jsonData.ZfsEntities
	}
	return res, err
}

func ZfsGetLastSnapshot(apiZfs string, dataset string) (string, error) {
	var (
		apiResponse []byte
		err         error
		res         string
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseGeneric
	)
	param["dataset"] = dataset

	if apiResponse, err = apiCall(apiZfs, "lastsnapshot", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		if jsonData.Status != "error" {
			res = fmt.Sprintf("%v", jsonData.Data["lastsnapshot"])
		} else {
			err = errors.New(jsonData.ErrorMessage)
		}
	}
	return res, err
}

func ZfsGetCloneInfo(apiZfs string, dataset string) (map[string]string, error) {
	var (
		apiResponse []byte
		err         error
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseGeneric
		res         map[string]string = make(map[string]string)
	)
	param["dataset"] = dataset
	if apiResponse, err = apiCall(apiZfs, "cloneinfo", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		res = jsonData.GetData()
	}
	/*var (
		out []byte
		err error
		res map[string]string
	)

	cmd := exec.Command(ZFS_BINARY, "get", "-Hpo", "value", "origin,written", dataset)
	if out, err = cmd.CombinedOutput(); err != nil {
		log.Println(err)
	}
	out_split := strings.Split(string(out), "\n")
	res["origin"] = out_split[0]
	res["written"] = out_split[1]*/
	return res, err
}

func ZfsCreateSnapshot(apiZfs string, snapsource string, snapname string) error {
	var (
		err         error
		apiResponse []byte
		res         jsonResponseGeneric
	)

	param := make(map[string]string)
	param["snapsource"] = snapsource
	param["snapname"] = snapname
	if apiResponse, err = apiCall(apiZfs, "snapshot", param); err != nil {
		log.Println(res.ErrorMessage)
	} else {
		json.Unmarshal(apiResponse, &res)
		if res.Status == "error" {
			err = errors.New(res.ErrorMessage)
		}
	}

	return err
}

func ZfsRollback(snapname string) ([]string, error) {
	var (
		err error
		out []byte
		res []string
	)
	if snapname != "" {
		cmd := exec.Command(ZFS_BINARY, "rollback -r", snapname)
		if out, err = cmd.CombinedOutput(); err != nil {
			log.Println(err)
			res = strings.Split(string(out), "\n")
		}

	} else {
		err = errors.New("missing snapshot name.")
	}
	return res, err
}
