package sunnytripower

// https://sma3013215929/dyn/getDashValues.json
type livePowerResponse struct {
	PvPowerW int
}

type rawApiResponse struct {
	Result struct {
		Zero17AXxxxx6B9 struct {
			Six1000046C200 struct {
				Num1 []struct {
					Val *int `json:"val,omitempty"`
				} `json:"1"`
			} `json:"6100_0046C200"`
		} `json:"017A-xxxxx6B9"`
	} `json:"result"`
}
