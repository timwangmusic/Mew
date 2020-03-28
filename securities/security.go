package securities

import "html/template"

// maps a subset of data returned from Robinhood instrument API of a security
type Instrument struct {
	Id     string       `json:"id"`     // The instrument id of this security
	Symbol string       `json:"symbol"` // ticker
	Name   string       `json:"name"`   // security name
	State  string       `json:"state"`  // active or not
	Url    template.URL `json:"url"` // Robinhood instrument URL
	Quote  template.URL `json:"quote"`
}

type SecurityGroup struct {
	Securities []Instrument
}
