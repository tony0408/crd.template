package exchange

type ExchangeName string

const (
	BLANK     ExchangeName = "BLANK"
	CRYPTOPIA ExchangeName = "CRYPTOPIA"
	FCOIN     ExchangeName = "FCOIN"
	COINEAL   ExchangeName = "COINEAL"
	ITIGER    ExchangeName = "ITIGER"
	BITFOREX    ExchangeName = "BITFOREX"
	KRAKEN    ExchangeName = "KRAKEN"
	BITRUE    ExchangeName = "BITRUE"
)

func (e *ExchangeManager) initExchangeNames() {
	supportList = append(supportList, CRYPTOPIA)
}
