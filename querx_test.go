package querx

import "testing"

var QuerxIPSSL = "192.168.192.222"
var QuerxPortSSL = 443
var QuerxIP = "192.168.192.236"
var QuerxPort = 80
var testQuerx = NewQuerx(QuerxIP, QuerxPort, false)

func TestQueryQuerx(t *testing.T) {
	queryQuerx(testQuerx, t)
}

func TestReadTemperatureSensors(t *testing.T) {
	readSensor(testQuerx, THSensorHumidity, t)
	readSensor(testQuerx, THSensorTemperature, t)
	readSensor(testQuerx, THSensorDewPoint, t)
}

func queryQuerx(q *Querx, t *testing.T) {
	err := q.QueryCurrent()
	if err != nil {
		t.Error("Could not querx "+q.Host, err)
	}
}

func readSensor(q *Querx, s int, t *testing.T) {
	tSensor, err := q.SensorByID(THSensorTemperature)
	if err != nil {
		t.Error("Could not get TemperatureSensor")
	}
	_, err = q.CurrentValue(tSensor)
	if err != nil {
		t.Error("Could get Value from Temperature Sensor")
	}
}
