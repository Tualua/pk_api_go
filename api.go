package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-yaml/yaml"
	"github.com/gorilla/mux"
)

type Config struct {
	Server struct {
		Port string `yaml:"port"`
		Host string `yaml:"host"`
	} `yaml:"server"`
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
	router.HandleFunc("/snapshot", apiSnapshot)
	router.HandleFunc("/bookmark", apiBookmark)
	router.HandleFunc("/clone", apiClone)
	router.HandleFunc("/destroy", apiDestroy)
	router.HandleFunc("/status", apiStatus)
	router.HandleFunc("/ipcstats", apiIpcStats)
	router.HandleFunc("/targetmount", apiTargetMount)
	router.HandleFunc("/targetenable", apiTargetEnable)
	router.HandleFunc("/targetdisable", apiTargetDisable)
	router.HandleFunc("/release", apiRelease)
	router.HandleFunc("/reload", apiReload)
	router.HandleFunc("/send", apiSend)
	router.HandleFunc("/sendlist", apiSendList)
	router.HandleFunc("/senddetails", apiSendDetails)
	router.HandleFunc("/receivelist", apiReceiveList)
	router.HandleFunc("/targetconfig", apiTargetConfig)
	router.HandleFunc("/targetinfo", apiTargetInfo)
	router.HandleFunc("/rollback", apiRollback)
	router.HandleFunc("/version", apiVersion)
	router.HandleFunc("/targetcreate", apiTargetCreate)
	router.HandleFunc("/diffcreate", apiDiffCreate)
	router.HandleFunc("/smartclone", apiSmartClone)
	router.HandleFunc("/lastsnapshot", apiLastSnapshot)
	router.HandleFunc("/startreceiving", apiStartReceiving)
	router.HandleFunc("/replicate", apiReplicate)
	router.HandleFunc("/smartclone2", apiSmartClone2)
	router.HandleFunc("/checkclone", apiCheckClone)
	router.HandleFunc("/test", apiTest)

	log.Fatal(http.ListenAndServe(addrString, router))
}

func apiSnapshot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "snapshot")
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
	fmt.Fprintf(w, "status")
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
	fmt.Fprintf(w, "checkclone")
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
