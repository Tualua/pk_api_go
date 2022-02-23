package main

import (
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

type XmlApi interface {
	SetAction(action string)
	SetVal(name, val string)
	Success()
	Error()
	Write()
}

func (x *XmlResponseGeneric) SetAction(act string) {
	x.Action = act
}

func (x *XmlResponseGeneric) Success() {
	x.Status = "success"
}

func (x *XmlResponseGeneric) Error(message string) {
	x.Status = "error"
	if x.Fields == nil {
		x.Fields = make(XmlFieldsMap)
	}
	x.Fields["errormessage"] = message
}

func (x *XmlResponseGeneric) SetVal(name, val string) {
	if x.Fields == nil {
		x.Fields = make(XmlFieldsMap)
	}
	x.Fields[name] = val
}

func (x *XmlResponseGeneric) Write(w *http.ResponseWriter) {
	fmt.Fprintf(*w, xml.Header)
	enc := xml.NewEncoder(*w)
	enc.Indent(" ", "  ")
	enc.Encode(x)
	fmt.Fprintf(*w, "\n")
}

type XmlData struct {
	XMLName xml.Name    `xml:"log"`
	Entries interface{} `xml:"entry"`
}

type XmlResponseGeneric struct {
	XMLName xml.Name `xml:"response"`
	Action  string   `xml:"action"`
	Status  string   `xml:"status"`
	Fields  XmlFieldsMap
}

type XmlResponseSC2 struct {
	XmlResponseGeneric
	DesktopDeviceId      string `xml:"deviceid"`
	DesktopTarget        string `xml:"target"`
	DesktopFile          string `xml:"file"`
	DesktopLastSnapshot  string `xml:"lastsnapshot"`
	DesktopOrigin        string `xml:"origin"`
	DesktopWritten       string `xml:"written"`
	DesktopCloneSnapshot string `xml:"clonesnapshot"`
	GamesDeviceId        string `xml:"deviceid"`
	GamesTarget          string `xml:"target"`
	GamesFile            string `xml:"file"`
	GamesLastSnapshot    string `xml:"lastsnapshot"`
	GamesOrigin          string `xml:"origin"`
	GamesWritten         string `xml:"written"`
	GamesCloneSnapshot   string `xml:"clonesnapshot"`
}

type ZfsXmlResponseListAll struct {
	XmlResponseGeneric
	Data []ZfsEntity `xml:"zfsentity"`
}

type ZfsXmlResponseGeneric struct {
	XMLName      xml.Name `xml:"response"`
	Action       string   `xml:"action"`
	Status       string   `xml:"status"`
	ErrorMessage string   `xml:"errormessage"`
}
type XmlResponse struct {
	XmlResponseGeneric
	Log *XmlData
}

func (x *XmlResponse) Write(w *http.ResponseWriter) {
	fmt.Fprintf(*w, xml.Header)
	enc := xml.NewEncoder(*w)
	enc.Indent(" ", "  ")
	enc.Encode(x)
	fmt.Fprintf(*w, "\n")
}

type XmlFieldsMap map[string]string

type xmlFieldEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}

//XML Encoder for map
func (m XmlFieldsMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if len(m) == 0 {
		return nil
	}

	for k, v := range m {
		e.Encode(xmlFieldEntry{XMLName: xml.Name{Local: k}, Value: v})
	}

	return nil
}

//XML Decoder for map
func (m *XmlFieldsMap) UnmarshalXML(e *xml.Decoder, start xml.StartElement) error {
	log.Println("unmarshal")
	return nil
}
