package config

// Configurations exported
type Configurations struct {
	Broker BrokerConfigurations
}

// BrokerConfigurations exported
type BrokerConfigurations struct {
	Name               string
	EncodedCredentials string
}
