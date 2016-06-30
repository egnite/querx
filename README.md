# querx v0.1
Interact with Querx Smart Ethernet Sensors using Go (golang)

## usage
```go
func main(){
  ip := "192.168.192.236"
  port  := 80
  querx := NewQuerx(QuerxIP, QuerxPort, false)

  temperatureSensor, err := q.SensorByID(THSensorTemperature)
  if err != nil {
		t.Error("Could not access TemperatureSensor")
	}
	currentTemperatue, err := q.CurrentValue(tSensor)
	if err != nil {
		t.Error("Could get Value from Temperature Sensor")
	}
}
```

## More information on Querx Smart Ethernet sensors
Find more information on Querx on the [product page](http://sensors.egnite.de)
and on [egnite's website](http://www.egnite.de)
