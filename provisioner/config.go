package provisioner

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Config struct {
	ApplicationId         string  `json:"applicationId"`
	ApiKey                string  `json:"apiKey"`
	UserId                string  `json:"userId"`
	DeviceType            string  `json:"deviceType"`
	RegisteredAt          float64 `json:"registered_at,omitempty"`
	AppUpdatePollInterval string  `json:"appUpdatePollInterval"`
	ListenPort            string  `json:"listenPort"`
	VpnPort               string  `json:"vpnPort"`
	ApiEndpoint           string  `json:"apiEndpoint"`
	VpnEndpoint           string  `json:"vpnEndpoint"`
	RegistryEndpoint      string  `json:"registryEndpoint"`
	DeltaEndpoint         string  `json:"deltaEndpoint"`
	PubnubSubscribeKey    string  `json:"pubnubSubscribeKey"`
	PubnubPublishKey      string  `json:"pubnubPublishKey"`
	MixpanelToken         string  `json:"mixpanelToken"`

	// See json.go/parseConfig() for more details on what this is for.
	InitialRaw map[string]interface{} `json:"-"`
}

func (a *Api) readConfig() (*Config, error) {
	if bytes, err := ioutil.ReadFile(a.ConfigPath); os.IsNotExist(err) {
		// We'll create a new config.json.
		log.Printf("Empty %s, will create new on write.\n", a.ConfigPath)
		return parseConfig("{}")
	} else if err != nil {
		return nil, err
	} else {
		return parseConfig(string(bytes))
	}
}

func (a *Api) writeConfig(conf *Config) error {
	if str, err := stringifyConfig(conf); err != nil {
		return err
	} else {
		return atomicWrite(a.ConfigPath, str)
	}
}

func (c *Config) ProvisionedState() ProvisionedState {
	if c.ApplicationId == "" || c.ApiKey == "" {
		return Unprovisioned
	}

	if c.RegisteredAt == 0 {
		return Provisioning
	}

	return Provisioned
}

func (c *Config) DetectDeviceType() error {
	if deviceType, err := ScanDeviceType(); err != nil {
		return err
	} else {
		c.DeviceType = deviceType.String()
	}

	return nil
}
