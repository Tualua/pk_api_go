package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
)

func scstGetIscsiSessions(apiScst string, tgtid string) (res []string, err error) {
	var (
		apiResponse []byte
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseList
	)
	param["tgtid"] = tgtid
	if apiResponse, err = apiCall(apiScst, "iscsisessions", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		if jsonData.Status != "error" {
			res = jsonData.Data
		} else {
			err = errors.New(jsonData.ErrorMessage)
		}
	}
	return
}

func ScstCheckIscsiSessions(apiScst string, tgtid string) (err error) {
	var (
		res []string
	)
	if res, err = scstGetIscsiSessions(apiScst, tgtid); err != nil {
		log.Println(err.Error())
	} else {
		if len(res) > 0 {
			err = fmt.Errorf("there is an active iscsi session: %s", res[0])
		}
	}
	return
}

func ScstDeactivateDevice(apiScst string, devid string) (err error) {
	var (
		apiResponse []byte
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseGeneric
	)
	param["devid"] = devid
	if apiResponse, err = apiCall(apiScst, "deactdev", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		if jsonData.Status == "error" {
			err = errors.New(jsonData.ErrorMessage)
		}
	}
	return
}

func ScstActivateDevice(apiScst string, devid string) (err error) {
	var (
		apiResponse []byte
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseGeneric
	)
	param["devid"] = devid
	if apiResponse, err = apiCall(apiScst, "actdev", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		if jsonData.Status == "error" {
			err = errors.New(jsonData.ErrorMessage)
		}
	}
	return
}

func ScstGetIscsiTargetParams(apiScst string, tgtid string) (res map[string]string, err error) {
	var (
		apiResponse []byte
		param       map[string]string = make(map[string]string)
		jsonData    jsonResponseGeneric
	)
	param["tgtid"] = tgtid
	if apiResponse, err = apiCall(apiScst, "iscsitargetparams", param); err != nil {
		log.Println(err.Error())
	} else {
		json.Unmarshal(apiResponse, &jsonData)
		if jsonData.Status == "error" {
			err = errors.New(jsonData.ErrorMessage)
		} else {
			res = jsonData.GetData()
		}
	}
	return
}
