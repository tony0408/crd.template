package exchange

type ExchangeName string

const (
	BLANK     ExchangeName = "BLANK"
	CRYPTOPIA ExchangeName = "CRYPTOPIA"
	FCOIN	  ExchangeName = "FCOIN"
)

func (e *ExchangeManager) initExchangeNames() {
	supportList = append(supportList, CRYPTOPIA)
}
