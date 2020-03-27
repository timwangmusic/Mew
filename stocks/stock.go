package stocks

type Stock struct {
	Ticker      string
	CompanyName string
}

type Group struct {
	Stocks []Stock
}
