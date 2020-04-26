package config

// Configurations exported
type Configurations struct {
	Broker BrokerConfigurations
}

// BrokerConfigurations exported
type BrokerConfigurations struct {
	Name               string
	User               string
	EncodedCredentials string
}
