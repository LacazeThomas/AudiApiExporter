package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/go-resty/resty"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func getToStruct(client *resty.Client, url string, param map[string]string, target interface{}) {
	resp, err := client.R().
		SetQueryParams(param).
		Get(url)

	checkErr(err)
	err = json.Unmarshal(resp.Body(), &target)
	checkErr(err)
	return
}

func postToStruct(client *resty.Client, url string, param map[string]string, target interface{}) {
	resp, err := client.R().
		SetQueryParams(param).
		Post(url)

	checkErr(err)
	err = json.Unmarshal(resp.Body(), &target)
	checkErr(err)
	return
}

func getMetrics(u UserProfil) (r AudiInfo) {
	client := resty.New()

	postToStruct(client, "https://msg.audi.de/fs-car/core/auth/v1/Audi/DE/token", map[string]string{
		"grant_type": "password",
		"username":   u.Username,
		"password":   u.Password,
	}, &r.audiToken)

	getToStruct(client.
		SetHeader("Accept", "application/json").
		SetHeader("X-App-ID", "de.audi.mmiapp").
		SetHeader("X-App-Name", "MMIconnect").
		SetHeader("X-App-Version", "2.8.3").
		SetHeader("X-Brand", "audi").
		SetHeader("X-Country-Id", "DE").
		SetHeader("X-Language-Id", "de").
		SetHeader("X-Platform", "google").
		SetHeader("User-Agent", "okhttp/2.7.4").
		SetHeader("ADRUM_1", "isModule:true").
		SetHeader("ADRUM", "isAray:true").
		SetHeader("Authorization", "AudiAuth 1 "+r.audiToken.AccessToken),
		"https://msg.audi.de/fs-car/myaudi/carservice/v2/Audi/DE/vehicles", map[string]string{}, &r.audiCar)

	getToStruct(client.
		SetHeader("Accept", "application/json").
		SetHeader("X-App-ID", "de.audi.mmiapp").
		SetHeader("X-App-Name", "MMIconnect").
		SetHeader("X-App-Version", "2.8.3").
		SetHeader("X-Brand", "audi").
		SetHeader("X-Country-Id", "DE").
		SetHeader("X-Language-Id", "de").
		SetHeader("X-Platform", "google").
		SetHeader("User-Agent", "okhttp/2.7.4").
		SetHeader("ADRUM_1", "isModule:true").
		SetHeader("ADRUM", "isAray:true").
		SetHeader("Authorization", "AudiAuth 1 "+r.audiToken.AccessToken),
		"https://msg.audi.de/fs-car/bs/vsr/v1/Audi/DE/vehicles/"+r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN+"/status", map[string]string{}, &r.audiCarInfo)

	getToStruct(client.
		SetHeader("Accept", "application/json").
		SetHeader("X-App-ID", "de.audi.mmiapp").
		SetHeader("X-App-Name", "MMIconnect").
		SetHeader("X-App-Version", "2.8.3").
		SetHeader("X-Brand", "audi").
		SetHeader("X-Country-Id", "DE").
		SetHeader("X-Language-Id", "de").
		SetHeader("X-Platform", "google").
		SetHeader("User-Agent", "okhttp/2.7.4").
		SetHeader("ADRUM_1", "isModule:true").
		SetHeader("ADRUM", "isAray:true").
		SetHeader("Authorization", "AudiAuth 1 "+r.audiToken.AccessToken),
		"https://msg.audi.de/fs-car/bs/batterycharge/v1/Audi/DE/vehicles/"+r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN+"/charger", map[string]string{}, &r.audiCharger)
	return
}

func getMetricsInf(u UserProfil, m Metrics) {
	go func() {
		for {
			r := getMetrics(u)

			for i := 0; i < len(r.audiCarInfo.StoredVehicleDataResponse.VehicleData.Data); i++ {
				for y := 0; y < len(r.audiCarInfo.StoredVehicleDataResponse.VehicleData.Data[i].Field); y++ {

					val := r.audiCarInfo.StoredVehicleDataResponse.VehicleData.Data[i].Field[y]
					switch val.ID {
					case "0x0301030006":
						rangeHvacOff, err := strconv.ParseFloat(val.Value,64)
						checkErr(err)
						m.rangeHvacOff.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(rangeHvacOff)
					case "0x0101010002":
						totalMileage, err := strconv.ParseFloat(val.Value,64)
						checkErr(err)
						m.totalMileage.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(totalMileage)
					case "0x0301030002":
						batteryLevel, err := strconv.ParseFloat(val.Value,64)
						checkErr(err)
						m.batteryLevel.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(batteryLevel)
					}
				}
			}

			var plugStatus float64
			var err error
			if(r.audiCharger.Charger.Status.ChargingStatusData.EnergyFlow.Content == "off"){
				plugStatus = float64(-1)
			}else{
				plugStatus, err = strconv.ParseFloat(r.audiCharger.Charger.Status.ChargingStatusData.EnergyFlow.Content,64)
			}
			checkErr(err)
			m.plugStatus.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(plugStatus)

			m.lastUpdateTime.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(float64(r.audiCharger.Charger.Status.BatteryStatusData.StateOfCharge.Timestamp.Unix()))
			m.timeRequiredToFullSlow.WithLabelValues(r.audiCar.GetUserVINsResponse.CSIDVins[0].VIN).Set(float64(r.audiCharger.Charger.Status.BatteryStatusData.RemainingChargingTime.Content))
			time.Sleep(600 * time.Second)
		}
	}()
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	
	u := UserProfil{}
	err := env.Parse(&u)
	checkErr(err)


	metrics := initMetrics()
	prometheus.Register(metrics.chargeStatus)
	prometheus.Register(metrics.rangeHvacOff)
	prometheus.Register(metrics.plugStatus)
	prometheus.Register(metrics.lastUpdateTime)
	prometheus.Register(metrics.batteryLevel)
	prometheus.Register(metrics.totalMileage)
	prometheus.Register(metrics.timeRequiredToFullSlow)

	getMetricsInf(u, metrics)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":9158", nil)
}
