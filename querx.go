package querx

import (
	"crypto/tls"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"time"
)

const currentValueURL string = "/tpl/document.cgi?tpl/j/current.tpl&format=xml"
const loginURL string = "/login.cgi"

//THSensorTemperature is a constant holding the sensor number for temperature
const THSensorTemperature int = 0

//THSensorHumidity is a constant holding the sensor number for temperature
const THSensorHumidity int = 1

//THSensorDewPoint is a constant holding the sensor number for temperature
const THSensorDewPoint int = 2

//QuerxData provides a Querx
type QuerxData struct {
	Version   string  `xml:"version"`
	Hostname  string  `xml:"hostname"`
	IP        string  `xml:"ip"`
	Port      uint64  `xml:"port"`
	DateGmt   string  `xml:"date_gmt"`
	DateLocal string  `xml:"date_local"`
	Contact   string  `xml:"contact"`
	Location  string  `xml:"location"`
	Sensors   Sensors `xml:"sensors"`
	Data      Data    `xml:"data"`
}

//Sensors is a container for Sensors
type Sensors struct {
	Sensor []Sensor `xml:"sensor"`
}

//Querx provides a frame for a Querx
type Querx struct {
	Host           string
	Port           int
	TLS            bool
	Current        QuerxData
	Datalogger     QuerxData
	hasCurrentData bool
	hasLoggedData  bool
	Type           string
	client         *http.Client
}

//Sensor holds a sensor entry in
type Sensor struct {
	ID         string  `xml:"id,attr"`
	Name       string  `xml:"name,attr"`
	Unit       string  `xml:"unit,attr"`
	Status     int     `xml:"status,attr"`
	UpperLimit float64 `xml:"uplim,attr"`
	LowerLimit float64 `xml:"lolim,attr"`
}

type Data struct {
	Record []Record `xml:"record"`
}

//Record holds meta data of a single record
type Record struct {
	Timestamp string  `xml:"timestamp"`
	Date      string  `xml:"date"`
	Time      string  `xml:"datetime"`
	Entry     []Entry `xml:"entry"`
}

//Entry hold a data entry
type Entry struct {
	Sensorid string  `xml:"sensorid,attr"`
	Name     string  `xml:"name,attr"`
	Value    float64 `xml:"value,attr"`
	Trend    float64 `xml:"trend,attr"`
}

//NewQuerx creates a new Querx object from
func NewQuerx(host string, port int, tlsEnabled bool) *Querx {
	q := Querx{}
	q.Host = host
	q.Port = port
	q.TLS = tlsEnabled
	jar, _ := cookiejar.New(nil)
	if tlsEnabled {
		tlc := &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         host,
		}
		transport := &http.Transport{TLSClientConfig: tlc}
		timeout := 20 * time.Second
		q.client = &http.Client{
			Transport: transport,
			Timeout:   timeout,
			Jar:       jar,
		}
	} else {
		q.client = &http.Client{
			Jar: jar,
		}
	}
	return &q
}

func (q *Querx) Login(user string, password string) (err error) {
	protoPrefix := "http://"
	if q.TLS {
		protoPrefix = "https://"
	}
	posturl := protoPrefix + q.Host + ":" + strconv.Itoa(q.Port) + loginURL
	querxurl, err := url.Parse(protoPrefix + q.Host + "/")
	formdata := url.Values{}
	formdata.Add("login_user", user)
	formdata.Add("login_pass", password)
	resp, err := q.client.PostForm(posturl, formdata)
	q.client.Jar.SetCookies(querxurl, resp.Cookies())
	return err
}

//Returns the current value for a given sensor
func (q *Querx) CurrentValue(s Sensor) (value float64, err error) {
	value = 0.0
	if !q.hasCurrentData {
		err := q.QueryCurrent()
		if err != nil {
			return value, err
		}
	}
	for _, entry := range q.Current.Data.Record[0].Entry {
		if entry.Sensorid == s.ID {
			value = entry.Value
		}
	}
	return value, nil
}

//QueryCurrent retrieves current readings from Querx
func (q *Querx) QueryCurrent() (err error) {

	protoPrefix := "http://"
	if q.TLS {
		protoPrefix = "https://"
	}
	url := protoPrefix + q.Host + ":" + strconv.Itoa(q.Port) + currentValueURL

	resp, err := q.client.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return err
	}
	err = xml.Unmarshal(body, &q.Current)
	if err != nil {
		return err
	}
	q.hasCurrentData = true
	return nil
}

//SensorByID gets sensor by value
func (q *Querx) SensorByID(id int) (s Sensor, err error) {
	if id >= len(q.Current.Sensors.Sensor) {
		return Sensor{}, errors.New("Can not find Sensor with ID " + strconv.Itoa(id))
	}
	return q.Current.Sensors.Sensor[id], nil
}

func (s *Sensor) Alerts(messages []string) {
	messages = make([]string, 5)
	var i uint8
	for i = 0; i < 4; i++ {
		if ((1 << i) & s.Status) > 0 {
			messages = append(messages, "Fehlerchen")
		}
	}
}
