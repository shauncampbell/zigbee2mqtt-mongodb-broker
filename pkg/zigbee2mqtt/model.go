package zigbee2mqtt

// Device is a structure which contains the information about an zigbee2mqtt device.
type Device struct {
	FriendlyName       string           `json:"friendly_name"`
	IEEEAddress        string           `json:"ieee_address"`
	InterviewCompleted bool             `json:"interview_completed"`
	Interviewing       bool             `json:"interviewing"`
	NetworkAddress     int              `json:"network_address"`
	Supported          bool             `json:"supported"`
	Type               string           `json:"type"`
	Definition         DeviceDefinition `json:"definition"`
	ModelID            string           `json:"model_id"`
}

// DeviceDefinition is a structure which contains the definition of a zigbee2mqtt device.
type DeviceDefinition struct {
	Description string            `json:"description"`
	Exposes     []DeviceAttribute `json:"exposes"`
	Model       string            `json:"model"`
	SupportsOTA bool              `json:"supports_ota"`
	Vendor      string            `json:"vendor"`
}

// DeviceAttribute is an attribute which the device exposes to the user.
type DeviceAttribute struct {
	Access      int    `json:"access"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Property    string `json:"property"`
	Type        string `json:"type"`
	Unit        string `json:"unit"`
	ValueMax    int    `json:"value_max"`
	ValueMin    int    `json:"value_min"`
}
