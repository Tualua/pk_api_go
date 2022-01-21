package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Tualua/pk_api_go/zfs_go"
	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
)

const tmpPath string = "/tmp"

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
}

//Structs for action=status
type XmlLog struct {
	XMLName xml.Name `xml:"log"`
	Log     []zfs_go.ZfsEntity
}
type XmlStatus struct {
	XMLName xml.Name `xml:"response"`
	Action  string   `xml:"action"`
	Status  string   `xml:"status"`
	Log     XmlLog
}

type XmlResponse map[string]string

type xmlMapEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

//XML Encoder for map
func (m XmlResponse) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}
	start.Name.Local = "response"
	err := e.EncodeToken(start)
	if err != nil {
		return err
	}

	for k, v := range m {
		e.Encode(xmlMapEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return e.EncodeToken(start.End())
}

func NewConfig(configPath string) (*Config, error) {
	config := &Config{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

func run(cfg *Config) {
	router := mux.NewRouter().StrictSlash(true)
	addrString := cfg.Server.Host + ":" + cfg.Server.Port
	router.Path("/api").Queries("action", "snapshot").HandlerFunc(apiSnapshot)
	router.Path("/api").Queries("action", "bookmark").HandlerFunc(apiBookmark)
	router.Path("/api").Queries("action", "clone").HandlerFunc(apiClone)
	router.Path("/api").Queries("action", "destroy").HandlerFunc(apiDestroy)
	router.Path("/api").Queries("action", "status").HandlerFunc(apiStatus)
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
	router.Path("/api").Queries("action", "smartclone").HandlerFunc(apiSmartClone)
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
	).HandlerFunc(apiCheckClone)
	router.Path("/api").Queries("action", "test").HandlerFunc(apiTest)

	log.Fatal(http.ListenAndServe(addrString, router))
}

func apiSnapshot(w http.ResponseWriter, r *http.Request) {
	var (
		res      XmlStatus
		err      error
		snapshot string
	)

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
func apiStatus(w http.ResponseWriter, r *http.Request) {
	var (
		res      XmlStatus
		err      error
		datasets []zfs_go.ZfsEntity
	)
	res.Action = "status"
	if datasets, err = zfs_go.ZfsListAll(); err != nil {
		res.Status = "error"
	} else {
		res.Status = "success"
		res.Log.Log = datasets
	}

	fmt.Fprintf(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent(" ", "  ")
	enc.Encode(res)
	fmt.Fprintf(w, "\n")
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
	fmt.Fprintf(w, "smartclone")
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
func apiCheckClone(w http.ResponseWriter, r *http.Request) {
	var (
		res          XmlResponse
		err          error
		lastSnapshot string
		snapshotInfo map[string]string
	)
	res = make(map[string]string)
	res["action"] = "checkclone"
	if lastSnapshot, err = zfs_go.ZfsGetLastSnapshot(mux.Vars(r)["clonesource"]); err != nil {
		res["status"] = "error"
		res["errormessage"] = err.Error()
	} else {
		res["lastsnapshot"] = lastSnapshot
	}
	if snapshotInfo, err = zfs_go.ZfsGetSnapshotInfo(mux.Vars(r)["clonename"]); snapshotInfo["origin"] == "-" {
		res["status"] = "error"
		res["errormessage"] = fmt.Sprintf("%s is not clone.", res["clonename"])
	} else {
		if lastSnapshot == snapshotInfo["origin"] {
			res["origin"] = snapshotInfo["origin"]
			res["written"] = snapshotInfo["written"]
			res["status"] = "error"
			res["errormessage"] = "actual clone. nothing to do"
		} else {
			res["status"] = "success"
		}
	}

	fmt.Fprintf(w, xml.Header)
	enc := xml.NewEncoder(w)
	enc.Indent(" ", "  ")
	enc.Encode(res)
	fmt.Fprintf(w, "\n")
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
