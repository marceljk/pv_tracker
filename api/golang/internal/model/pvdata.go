package model

type PvData struct {
	BatteryPercent    uint8 `json:"batteryPercent"`
	PowerConsumptionW int   `json:"powerConsumption"`
	GridPowerW        int   `json:"gridPower"`
	BatteryPowerW     int   `json:"batteryPower"`
	PvPowerW          int   `json:"pvPower"`
}
