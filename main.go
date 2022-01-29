package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

const tmpPath string = "/tmp"

func loggingMiddleware(next http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout, next)
}

func run(cfg *Config) {
	router := mux.NewRouter().StrictSlash(true)
	addrString := cfg.Server.Host + ":" + cfg.Server.Port
	router.Path("/api").Queries("action", "snapshot",
		"snapsource", "{snapsource}",
		"snapname", "{snapname}").HandlerFunc(apiSnapshot(cfg.Apis.ZfsApi))
	router.Path("/api").Queries("action", "bookmark").HandlerFunc(apiBookmark)
	router.Path("/api").Queries("action", "clone").HandlerFunc(apiClone)
	router.Path("/api").Queries("action", "destroy").HandlerFunc(apiDestroy)
	router.Path("/api").Queries("action", "status").HandlerFunc(apiStatus(cfg.Apis.ZfsApi))
	router.Path("/api").Queries("action", "ipcstats").HandlerFunc(apiIpcStats)
	router.Path("/api").Queries("action", "targetmount").HandlerFunc(apiTargetMount)
	router.Path("/api").Queries("action", "targetenable").HandlerFunc(apiTargetEnable)
	router.Path("/api").Queries("action", "targetdisable").HandlerFunc(apiTargetDisable)
	router.Path("/api").Queries("action", "release").HandlerFunc(apiRelease)
	router.Path("/api").Queries("action", "reload").HandlerFunc(apiReload)
	router.Path("/api").Queries("action", "send").HandlerFunc(apiSend)
	router.Path("/api").Queries("action", "sendlist").HandlerFunc(apiSendList)
	router.Path("/api").Queries("action", "senddetails").HandlerFunc(apiSendDetails)
	router.Path("/api").Queries("action", "receivelist").HandlerFunc(apiReceiveList)
	router.Path("/api").Queries("action", "targetconfig").HandlerFunc(apiTargetConfig)
	router.Path("/api").Queries("action", "targetinfo").HandlerFunc(apiTargetInfo)
	router.Path("/api").Queries("action", "rollback").HandlerFunc(apiRollback)
	router.Path("/api").Queries("action", "version").HandlerFunc(apiVersion)
	router.Path("/api").Queries("action", "targetcreate").HandlerFunc(apiTargetCreate)
	router.Path("/api").Queries("action", "diffcreate").HandlerFunc(apiDiffCreate)
	router.Path("/api").Queries("action", "smartclone",
		"clonesource", "{clonesource}",
		"clonename", "{clonename}",
		"deviceid", "{deviceid}",
	).HandlerFunc(apiSmartClone)
	router.Path("/api").Queries("action", "lastsnapshot").HandlerFunc(apiLastSnapshot)
	router.Path("/api").Queries("action", "startreceiving").HandlerFunc(apiStartReceiving)
	router.Path("/api").Queries("action", "replicate").HandlerFunc(apiReplicate)
	router.Path("/api").Queries("action", "smartclone2",
		"systemmaster", "{systemmaster}",
		"gamesmaster", "{gamesmaster}",
		"systemclone", "{systemclone}",
		"gamesid", "{gamesid}").HandlerFunc(apiSmartClone2)
	router.Path("/api").Queries("action", "checkclone",
		"clonesource", "{clonesource}",
		"clonename", "{clonename}",
	).HandlerFunc(apiCheckClone(cfg.Apis.ZfsApi))
	router.Path("/api").Queries("action", "test").HandlerFunc(apiTest)
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
		fmt.Fprintf(w, xml.Header)
		enc := xml.NewEncoder(w)
		enc.Indent(" ", "  ")
		enc.Encode(res_out)
		fmt.Fprintf(w, "\n")

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
func apiSmartClone(w http.ResponseWriter, r *http.Request) {
	/*var (
		res          XmlResponse
		err          error
		lastSnapshot string
		snapshotInfo map[string]string
	)
	res.Action = "smartclone"
	if lastSnapshot, err = ZfsGetLastSnapshot(mux.Vars(r)["clonesource"]); err != nil {
		res.Status = "error"
		res.Fields["errormessage"] = err.Error()
	} else {
		if lastSnapshot == "" {
			res.Status = "error"
			res.Fields["errormessage"] = fmt.Sprintf("there is no any snapshot in %s", mux.Vars(r)["clonesource"])
		} else {
			res.Fields["lastsnapshot"] = lastSnapshot
			if snapshotInfo, err = ZfsGetSnapshotInfo(mux.Vars(r)["clonename"]); snapshotInfo["origin"] == "-" {
				res.Status = "error"
				res.Fields["errormessage"] = fmt.Sprintf("%s is not clone.", mux.Vars(r)["clonename"])
			} else {
				res.Fields["origin"] = snapshotInfo["origin"]
				res.Fields["written"] = snapshotInfo["written"]
				// $lun2 - WTF?

				if snapshotInfo["origin"] != lastSnapshot || snapshotInfo["written"] != "0" {
					// lun2 checking?

				}
			}
		}
	}
	fmt.Fprintf(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent(" ", "  ")
	enc.Encode(res)
	fmt.Fprintf(w, "\n")*/

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
func apiSmartClone2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "smartclone2")
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
				log.Println(cloneinfo)
			}
			res.Write(&w)
		}
	}
	/*var (
		res          XmlResponseGeneric
		err          error
		lastSnapshot string
		snapshotInfo map[string]string
	)
	res.SetAction("checkclone")
	if lastSnapshot, err = ZfsGetLastSnapshot(mux.Vars(r)["clonesource"]); err != nil {
		res.Status = "error"
		res.Fields["errormessage"] = err.Error()
	} else {
		res.Fields["lastsnapshot"] = lastSnapshot
	}
	if snapshotInfo, err = ZfsGetSnapshotInfo(mux.Vars(r)["clonename"]); snapshotInfo["origin"] == "-" {
		res.Status = "error"
		res.Fields["errormessage"] = fmt.Sprintf("%s is not clone.", mux.Vars(r)["clonename"])
	} else {
		if lastSnapshot == snapshotInfo["origin"] {
			res.Fields["origin"] = snapshotInfo["origin"]
			res.Fields["written"] = snapshotInfo["written"]
			res.Status = "error"
			res.Fields["errormessage"] = "actual clone. nothing to do"
		} else {
			res.Status = "success"
		}
	}*/
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
