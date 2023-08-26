package service

type Grid struct {
	MatchFound bool `json:"matchFound"`
	Name       bool `json:"name"`  //name 8
	Value      bool `json:"value"` // value 14
	IsDevice4G bool `json:"isDevice4G"`
	IsDevice5G bool `json:"isDevice5G"`
}

type ResourceSpecCharacteristicValue struct {
	Value string `json:"value"`
}

type ResourceSpecCharacteristic struct {
	Name                            string                            `json:"name"`
	ResourceSpecCharacteristicValue []ResourceSpecCharacteristicValue `json:"resourceSpecCharacteristicValue"`
}

type IoTDeviceSpecification struct {
	ID                         string                       `json:"_id"`
	Description                string                       `json:"description"`
	Href                       string                       `json:"href"`
	LastUpdate                 string                       `json:"lastUpdate"`
	Version                    string                       `json:"version"`
	Name                       string                       `json:"name"`
	ResourceSpecCharacteristic []ResourceSpecCharacteristic `json:"resourceSpecCharacteristic"`
}
