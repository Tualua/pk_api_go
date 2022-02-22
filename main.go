package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

// const tmpPath string = "/tmp"

type SmartCloneInfo struct {
	origin       string
	written      string
	lastsnapshot string
	actualclone  string
}

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func smartClone(apiZfs string, apiScst string, clonename string, clonesource string, deviceid string) (res SmartCloneInfo, err error) {
	var (
		lastSnapshot   string
		cloneinfo      map[string]string = make(map[string]string)
		zeroSnapExists bool
	)
	if lastSnapshot, err = ZfsGetLastSnapshot(apiZfs, clonesource); err != nil {
		fmt.Println(err.Error())
	} else {
		if lastSnapshot == "" {
			err = fmt.Errorf("there is no any snapshot in %s", clonesource)
		} else {
			res.lastsnapshot = lastSnapshot
			if cloneinfo, err = ZfsGetCloneInfo(apiZfs, clonename); err != nil {
				fmt.Println(err.Error())
			} else {
				// Check if dataset is clone
				if cloneinfo["origin"] == "" {
					err = fmt.Errorf("%s is not clone", clonename)
				} else {
					res.origin = cloneinfo["origin"]
					res.written = cloneinfo["written"]
					// Check if clone is modified or is not on last snapshot
					if cloneinfo["written"] != "0" || cloneinfo["origin"] != lastSnapshot {
						// Check if there are any established iSCSI session
						if err = ScstCheckIscsiSessions(apiScst, deviceid); err != nil {
							fmt.Println(err.Error())
						} else {
							// Deactivate device to make it avaliable for modifications
							if err = ScstDeactivateDevice(apiScst, deviceid); err != nil {
								fmt.Println(err.Error())
							} else {
								zeroSnapshot := clonename + "@0"
								if zeroSnapExists, err = ZfsCheckDatasetExists(apiZfs, zeroSnapshot); err != nil {
									fmt.Println(err.Error())
								} else {
									if cloneinfo["origin"] == lastSnapshot && zeroSnapExists {
										ZfsRollback(apiZfs, zeroSnapshot)
									} else {
										ZfsDestroy(apiZfs, clonename)
										ZfsCloneLast(apiZfs, clonename, clonesource)
										ZfsCreateSnapshot(apiZfs, clonename, "0")
									}
									if err = ScstActivateDevice(apiScst, deviceid); err != nil {
										fmt.Println(err.Error())
									}
								}
							}

						}
					} else {
						res.actualclone = "nothing to do"
					}
				}
			}
		}
	}
	return
}

func run(cfg *Config) {
	router := mux.NewRouter().StrictSlash(true)
	addrString := cfg.Server.Host + ":" + cfg.Server.Port
	router.Path("/").Queries("action", "snapshot",
		"snapsource", "{snapsource}",
		"snapname", "{snapname}").HandlerFunc(apiSnapshot(cfg.Apis.ZfsApi))
	router.Path("/").Queries("action", "bookmark").HandlerFunc(apiBookmark)
	router.Path("/").Queries("action", "clone").HandlerFunc(apiClone)
	router.Path("/").Queries("action", "destroy").HandlerFunc(apiDestroy)
	router.Path("/").Queries("action", "status").HandlerFunc(apiStatus(cfg.Apis.ZfsApi))
	router.Path("/").Queries("action", "ipcstats").HandlerFunc(apiIpcStats)
	router.Path("/").Queries("action", "targetmount").HandlerFunc(apiTargetMount)
	router.Path("/").Queries("action", "targetenable").HandlerFunc(apiTargetEnable)
	router.Path("/").Queries("action", "targetdisable").HandlerFunc(apiTargetDisable)
	router.Path("/").Queries("action", "release").HandlerFunc(apiRelease)
	router.Path("/").Queries("action", "reload").HandlerFunc(apiReload)
	router.Path("/").Queries("action", "send").HandlerFunc(apiSend)
	router.Path("/").Queries("action", "sendlist").HandlerFunc(apiSendList)
	router.Path("/").Queries("action", "senddetails").HandlerFunc(apiSendDetails)
	router.Path("/").Queries("action", "receivelist").HandlerFunc(apiReceiveList)
	router.Path("/").Queries("action", "targetconfig").HandlerFunc(apiTargetConfig)
	router.Path("/").Queries("action", "targetinfo").HandlerFunc(apiTargetInfo)
	router.Path("/").Queries("action", "rollback").HandlerFunc(apiRollback)
	router.Path("/").Queries("action", "version").HandlerFunc(apiVersion)
	router.Path("/").Queries("action", "targetcreate").HandlerFunc(apiTargetCreate)
	router.Path("/").Queries("action", "diffcreate").HandlerFunc(apiDiffCreate)
	router.Path("/").Queries("action", "smartclone",
		"clonesource", "{clonesource}",
		"clonename", "{clonename}",
		"deviceid", "{deviceid}",
	).HandlerFunc(apiSmartClone(cfg.Apis.ZfsApi, cfg.Apis.ScstApi))
	router.Path("/").Queries("action", "lastsnapshot").HandlerFunc(apiLastSnapshot)
	router.Path("/").Queries("action", "startreceiving").HandlerFunc(apiStartReceiving)
	router.Path("/").Queries("action", "replicate").HandlerFunc(apiReplicate)
	router.Path("/").Queries("action", "smartclone2",
		"systemmaster", "{systemmaster}",
		"gamesmaster", "{gamesmaster}",
		"gamesid", "{gamesid}").HandlerFunc(apiSmartClone2(cfg.Apis.ZfsApi, cfg.Apis.ScstApi))
	/*router.Path("/").Queries("action", "smartclone2",
	"systemmaster", "{systemmaster}",
	"gamesmaster", "{gamesmaster}",
	"systemclone", "{systemclone}",
	"gamesid", "{gamesid}").HandlerFunc(apiSmartClone2(cfg.Apis.ZfsApi, cfg.Apis.ScstApi))*/
	router.Path("/").Queries("action", "checkclone",
		"clonesource", "{clonesource}",
		"clonename", "{clonename}",
	).HandlerFunc(apiCheckClone(cfg.Apis.ZfsApi))
	router.Path("/").Queries("action", "test").HandlerFunc(apiTest)
	router.Use(loggingMiddleware)
	log.Fatal(http.ListenAndServe(addrString, router))
}

func apiSnapshot(apiZfs string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			res XmlResponse
			err error
		)
		res.SetAction("snapshot")

		if mux.Vars(r)["snapsource"] == "" {
			res.SetVal("snapsource", "null")
		} else {
			res.SetVal("snapsource", mux.Vars(r)["snapsource"])
		}
		if mux.Vars(r)["snapname"] == "" {
			res.SetVal("snapname", "null")
		} else {
			res.SetVal("snapname", mux.Vars(r)["snapname"])
		}

		if res.Fields["snapname"] == "null" || res.Fields["snapsource"] == "null" {
			res.Error("missing snapshot source or snapshot name.")
		} else {
			if err = ZfsCreateSnapshot(apiZfs, mux.Vars(r)["snapsource"], mux.Vars(r)["snapname"]); err != nil {
				res.Error("log file not empty.")
				Log := make([]string, 0)
				Log = append(Log, err.Error())
				res.Log = &XmlData{Entries: Log}
			} else {
				res.Success()
			}
		}

		res.Write(&w)
	}
}
func apiBookmark(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "bookmark")
}
func apiClone(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "clone")
}
func apiDestroy(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "destroy")
}
func apiStatus(apiZfs string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			res_in  []ZfsEntity
			res_out XmlResponse
			err     error
		)
		res_out.Action = "status"
		if res_in, err = ZfsListAll(apiZfs); err != nil {
			res_out.Error(err.Error())
		} else {
			res_out.Status = "success"
			res_out.Log = &XmlData{Entries: res_in}
		}
		res_out.Write(&w)
	}
}
func apiIpcStats(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "ipcstats")
}
func apiTargetMount(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetmount")
}
func apiTargetEnable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetenable")
}
func apiTargetDisable(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetdisable")
}
func apiRelease(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "release")
}
func apiReload(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "reload")
}
func apiSend(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "send")
}
func apiSendList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "sendlist")
}
func apiSendDetails(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "senddetails")
}
func apiReceiveList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "receivelist")
}
func apiTargetConfig(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetconfig")
}
func apiTargetInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetinfo")
}
func apiRollback(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "rollback")
}
func apiVersion(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetconfig")
}
func apiTargetCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "targetcreate")
}
func apiDiffCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "diffcreate")
}
func apiSmartClone(apiZfs string, apiScst string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			res_out   XmlResponse
			res_in    SmartCloneInfo
			tgtParams map[string]string
			err       error

			// lastSnapshot   string
			// cloneinfo      map[string]string = make(map[string]string)
			// zeroSnapExists bool
		)
		res_out.Action = "smartclone"
		res_out.SetVal("clonesource", mux.Vars(r)["clonesource"])
		res_out.SetVal("clonename", mux.Vars(r)["clonename"])
		res_out.SetVal("deviceid", mux.Vars(r)["deviceid"])
		if tgtParams, err = ScstGetIscsiTargetParams(apiScst, mux.Vars(r)["deviceid"]); err != nil {
			res_out.Error(err.Error())
		} else {
			if res_in, err = smartClone(apiZfs, apiScst, mux.Vars(r)["clonename"], mux.Vars(r)["clonesource"], mux.Vars(r)["deviceid"]); err != nil {
				res_out.Error(err.Error())
			} else {
				res_out.Success()
				res_out.SetVal("target", tgtParams["wwn"])
			}
		}
		if res_in.actualclone != "" {
			res_out.SetVal("actualclone", res_in.actualclone)
		}
		res_out.SetVal("lastsnapshot", res_in.lastsnapshot)
		res_out.SetVal("origin", res_in.origin)
		res_out.SetVal("written", res_in.written)

		/*
			//Get last snapshot of clone source
			if lastSnapshot, err = ZfsGetLastSnapshot(apiZfs, mux.Vars(r)["clonesource"]); err != nil {
				res.Error(err.Error())
			} else {
				if lastSnapshot == "" {
					res.Error(fmt.Sprintf("there is no any snapshot in %s", mux.Vars(r)["clonesource"]))
				} else {
					res.SetVal("lastsnapshot", lastSnapshot)
					if cloneinfo, err = ZfsGetCloneInfo(apiZfs, mux.Vars(r)["clonename"]); err != nil {
						res.Error(err.Error())
					} else {
						// Check if dataset is clone
						if cloneinfo["origin"] == "" {
							res.Error(fmt.Sprintf("%s is not clone.", mux.Vars(r)["clonename"]))
						} else {
							res.SetVal("origin", cloneinfo["origin"])
							res.SetVal("written", cloneinfo["written"])
							// Check if clone is modified or is not on last snapshot
							if cloneinfo["written"] != "0" || cloneinfo["origin"] != lastSnapshot {
								// Check if there are any established iSCSI session
								if err = ScstCheckIscsiSessions(apiScst, mux.Vars(r)["deviceid"]); err != nil {
									res.Error(err.Error())
								} else {
									// Deactivate device to make it avaliable for modifications
									if err = ScstDeactivateDevice(apiScst, mux.Vars(r)["deviceid"]); err != nil {
										res.Error(err.Error())
									} else {
										zeroSnapshot := mux.Vars(r)["clonename"] + "@0"
										if zeroSnapExists, err = ZfsCheckDatasetExists(apiZfs, zeroSnapshot); err != nil {
											res.Error(err.Error())
										} else {
											if cloneinfo["origin"] == lastSnapshot && zeroSnapExists {
												ZfsRollback(apiZfs, zeroSnapshot)
											} else {
												ZfsDestroy(apiZfs, mux.Vars(r)["clonename"])
												ZfsCloneLast(apiZfs, mux.Vars(r)["clonename"], mux.Vars(r)["clonesource"])
												ZfsCreateSnapshot(apiZfs, mux.Vars(r)["clonename"], "0")
											}
											if err = ScstActivateDevice(apiScst, mux.Vars(r)["deviceid"]); err != nil {
												res.Error(err.Error())
											} else {
												res.Success()
											}
										}
									}

								}
							} else {
								res.SetVal("actualclone", "nothing to do")
							}
						}
					}
				}
			}*/
		res_out.Write(&w)
	}
}
func apiZfsList(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "zfslist")
}
func apiLastSnapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "lastsnapshot")
}
func apiStartReceiving(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "startreceiving")
}
func apiReceivingLog(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "receivinglog")
}
func apiReportPrometheus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "reportprometheus")
}
func apiReplicate(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "replicate")
}
func apiSmartClone2(apiZfs string, apiScst string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			res         XmlResponseGeneric
			systemClone string
			gamesClone  string
			systemId    string
			gamesId     string
			res_in      SmartCloneInfo
			res_out     XmlResponseGeneric
			err         error
		)
		res.SetAction("smartclone2")
		params := r.URL.Query()

		if systemCloneCheck, ok := params["systemclone"]; !ok {
			if systemIdCheck, ok := params["systemid"]; !ok {
				res.Error("systemid not supplied")
			} else {
				systemClone = "data/kvm/desktop/" + systemIdCheck[0]
			}
		} else {
			systemClone = systemCloneCheck[0]
			if systemIdCheck, ok := params["systemid"]; !ok {
				systemCloneSplit := strings.Split(systemClone, "/")
				systemId = systemCloneSplit[len(systemCloneSplit)-1]
			} else {
				systemId = systemIdCheck[0]
			}
		}

		if gamesCloneCheck, ok := params["gamesclone"]; !ok {
			if gamesIdCheck, ok := params["gamesid"]; !ok {
				res.Error("gamesid not supplied")
			} else {
				systemClone = "data/kvm/desktop/" + gamesIdCheck[0]
			}
		} else {
			gamesClone = gamesCloneCheck[0]
			if gamesIdCheck, ok := params["gamesid"]; !ok {
				gamesCloneSplit := strings.Split(gamesClone, "/")
				gamesId = gamesCloneSplit[len(gamesCloneSplit)-1]
			} else {
				systemId = gamesIdCheck[0]
			}
		}

		if res_in, err = smartClone(apiZfs, apiScst, systemClone, mux.Vars(r)["systemmaster"], systemId); err != nil {
			res_out.Error(err.Error())
		}
		if res_in, err = smartClone(apiZfs, apiScst, systemClone, mux.Vars(r)["gamesmaster"], gamesId); err != nil {
			res_out.Error(err.Error())
		}
		log.Println(res_in)
	}
}
func apiCheckClone(apiZfs string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			lastSnapshot string
			res          XmlResponseGeneric
			err          error
			cloneinfo    map[string]string = make(map[string]string)
		)
		res.SetAction("checkclone")
		if lastSnapshot, err = ZfsGetLastSnapshot(apiZfs, mux.Vars(r)["clonesource"]); err != nil {
			res.Error(err.Error())
		} else {
			res.SetVal("lastsnapshot", lastSnapshot)
			if cloneinfo, err = ZfsGetCloneInfo(apiZfs, mux.Vars(r)["clonename"]); err != nil {
				res.Error(err.Error())
			} else {
				if cloneinfo["origin"] == "" {
					res.Error(fmt.Sprintf("%s is not clone.", mux.Vars(r)["clonename"]))
				} else {
					if lastSnapshot == cloneinfo["origin"] {
						res.SetVal("origin", cloneinfo["origin"])
						res.SetVal("written", cloneinfo["written"])
						res.Error("actual clone. nothing to do")
					} else {
						res.Success()
					}
				}
			}
			res.Write(&w)
		}
	}
}

func apiTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "test")
}

func main() {
	cfg, err := NewConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}
	run(cfg)

}
