package main

import (
	"time"

	log "github.com/sirupsen/logrus"
)

func checkErr(err error) {
	if err != nil {
		log.Error(err)
	}
}

type audiToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type audiCar struct {
	GetUserVINsResponse struct {
		CSIDVins []struct {
			CSID       string    `json:"CSID"`
			VIN        string    `json:"VIN"`
			Registered time.Time `json:"registered"`
		} `json:"CSIDVins"`
		VinsOnBlacklist int `json:"vinsOnBlacklist"`
	} `json:"getUserVINsResponse"`
}

type userProfil struct {
	Username    string `env:"LOGIN"`
	Password string `env:"PASSWORD"`
}

type audiCarInfo struct {
	StoredVehicleDataResponse struct {
		Vin         string `json:"vin"`
		VehicleData struct {
			Data []struct {
				ID    string `json:"id"`
				Field []struct {
					ID               string    `json:"id"`
					TsCarSentUtc     time.Time `json:"tsCarSentUtc"`
					TsCarSent        string    `json:"tsCarSent"`
					TsCarCaptured    string    `json:"tsCarCaptured"`
					TsTssReceivedUtc time.Time `json:"tsTssReceivedUtc"`
					MilCarCaptured   int       `json:"milCarCaptured"`
					MilCarSent       int       `json:"milCarSent"`
					Value            string    `json:"value"`
				} `json:"field"`
			} `json:"data"`
		} `json:"vehicleData"`
	} `json:"StoredVehicleDataResponse"`
}


type audiCharger struct {
	Charger struct {
		Settings struct {
			MaxChargeCurrent struct {
				Content   int       `json:"content"`
				Timestamp time.Time `json:"timestamp"`
			} `json:"maxChargeCurrent"`
		} `json:"settings"`
		Status struct {
			ChargingStatusData struct {
				ChargingMode struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"chargingMode"`
				ChargingReason struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"chargingReason"`
				ExternalPowerSupplyState struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"externalPowerSupplyState"`
				EnergyFlow struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"energyFlow"`
				ChargingState struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"chargingState"`
			} `json:"chargingStatusData"`
			CruisingRangeStatusData struct {
				EngineTypeFirstEngine struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"engineTypeFirstEngine"`
				PrimaryEngineRange struct {
					Content   int       `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"primaryEngineRange"`
			} `json:"cruisingRangeStatusData"`
			LedStatusData struct {
				LedColor struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"ledColor"`
				LedState struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"ledState"`
			} `json:"ledStatusData"`
			BatteryStatusData struct {
				StateOfCharge struct {
					Content   int       `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"stateOfCharge"`
				RemainingChargingTime struct {
					Content   int       `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"remainingChargingTime"`
				RemainingChargingTimeTargetSOC struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"remainingChargingTimeTargetSOC"`
			} `json:"batteryStatusData"`
			PlugStatusData struct {
				PlugState struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"plugState"`
				LockState struct {
					Content   string    `json:"content"`
					Timestamp time.Time `json:"timestamp"`
				} `json:"lockState"`
			} `json:"plugStatusData"`
		} `json:"status"`
	} `json:"charger"`
}

//AudiInfo global info from Electic services
type AudiInfo struct {
	audiCharger audiCharger
	audiToken audiToken
	audiCar audiCar
	audiCarInfo audiCarInfo
}
