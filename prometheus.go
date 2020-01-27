package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics listed are supported
type Metrics struct {
	chargeStatus           *prometheus.GaugeVec
	rangeHvacOff       *prometheus.GaugeVec
	plugStatus         *prometheus.GaugeVec
	lastUpdateTime     *prometheus.GaugeVec
	batteryLevel       *prometheus.GaugeVec
	totalMileage       *prometheus.GaugeVec
	timeRequiredToFullSlow *prometheus.GaugeVec
}

func initMetrics() (metrics Metrics) {

	metrics = Metrics{}
	metrics.chargeStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "chargeStatus",
		Help:      "Temperature in °C",
	}, []string{"vin"})

	metrics.rangeHvacOff = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "rangeHvacOff",
		Help:      "Atmospheric pressure in hPa",
	}, []string{"vin"})

	metrics.plugStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "plugStatus",
		Help:      "Wind speed in m/s",
	}, []string{"vin"})

	metrics.lastUpdateTime = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "lastUpdateTime",
		Help:      "Cloudiness in Percent",
	}, []string{"vin"})

	metrics.batteryLevel = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "batteryLevel",
		Help:      "Rain contents 3h",
	}, []string{"vin"})

	metrics.totalMileage = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "totalMileage",
		Help:      "The weather label.",
	}, []string{"vin"})

	metrics.timeRequiredToFullSlow = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "audi",
		Name:      "timeRequiredToFullSlow",
		Help:      "The weather label.",
	}, []string{"vin"})
	return
}
