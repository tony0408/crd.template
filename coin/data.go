package coin

// could be loaded from database in the future
// nowday, just expend this files to add more generic coins

func Init() {
	btc()
	bch()
	eth()
	ltc()
	usdt()

	ada()
}

func btc() {
	c := &Coin{}
	c.Code = BTC
	c.Name = "Bitcoin"
	c.Website = "https://www.bitcoin.org/"
	c.Explorer = "https://explore.bitcoinh.org/"
	coinmap[c.Code] = c
}

func bch() {
	c := &Coin{}
	c.Code = BCH
	c.Name = "Bitcoin Cash"
	c.Website = "https://www.bitcoincash.org/"
	c.Explorer = "https://explore.bitcoincash.org/"
	coinmap[c.Code] = c
}

func eth() {
	c := &Coin{}
	c.Code = ETH
	c.Name = "Ethereum"
	c.Website = "https://www.ethereum.org/"
	c.Explorer = "https://etherscan.io/"

	coinmap[c.Code] = c
}

func ltc() {
	c := &Coin{}
	c.Code = LTC
	c.Name = "Litecoin"
	c.Website = "https://litecoin.com/"
	c.Explorer = "https://chainz.cryptoid.info/ltc/"

	coinmap[c.Code] = c
}

func usdt() {
	c := &Coin{}
	c.Code = USDT
	c.Name = "Tether"
	c.Website = "https://tether.to/"
	c.Explorer = "https://omniexplorer.info/"

	coinmap[c.Code] = c
}

/****************************************************/

func ada() {
	c := &Coin{}
	c.Code = ADA
	c.Name = "Cardano"
	c.Website = "https://www.cardano.org/en/home/"
	c.Explorer = "https://www.cardano.org/en/home/"

	coinmap[c.Code] = c
}

func adx() {
	c := &Coin{}
	c.Code = ADX
	c.Name = "AdEx"
	c.Website = "https://www.adex.network/"
	c.Explorer = "https://etherscan.io/token/0x4470bb87d77b963a013db939be332f927f2b992e"

	coinmap[c.Code] = c
}

func appc() {
	c := &Coin{}
	c.Code = APPC
	c.Name = "AppCoins"
	c.Website = "https://appcoins.io/"
	c.Explorer = "https://etherscan.io/token/0x1a7a8bd9106f2b8d977e08582dc7d24c723ab0db"

	coinmap[c.Code] = c
}

func ark() {
	c := &Coin{}
	c.Code = ARK
	c.Name = "Ark"
	c.Website = "https://ark.io/"
	c.Explorer = "https://explorer.ark.io/"

	coinmap[c.Code] = c
}

func ast() {
	c := &Coin{}
	c.Code = AST
	c.Name = "AirSwap"
	c.Website = "https://www.airswap.io/"
	c.Explorer = "https://etherscan.io/token/0x27054b13b1b798b345b591a4d22e6562d47ea75a"

	coinmap[c.Code] = c
}

func bat() {
	c := &Coin{}
	c.Code = BAT
	c.Name = "Basic Attention Token"
	c.Website = "https://basicattentiontoken.org/"
	c.Explorer = "https://etherscan.io/token/Bat"

	coinmap[c.Code] = c
}

func bcd() {
	c := &Coin{}
	c.Code = BCD
	c.Name = "Bitcoin Diamond"
	c.Website = "http://btcd.io/"
	c.Explorer = "http://explorer.btcd.io/"

	coinmap[c.Code] = c
}

func bcpt() {
	c := &Coin{}
	c.Code = BCPT
	c.Name = "BlockMason Credit Protocol"
	c.Website = "https://blockmason.io/"
	c.Explorer = "https://etherscan.io/token/0x1c4481750daa5ff521a2a7490d9981ed46465dbd"

	coinmap[c.Code] = c
}

func blz() {
	c := &Coin{}
	c.Code = BLZ
	c.Name = "Bluzelle"
	c.Website = "https://bluzelle.com/"
	c.Explorer = "https://etherscan.io/token/0x5732046a883704404f284ce41ffadd5b007fd668"

	coinmap[c.Code] = c
}

func bnt() {
	c := &Coin{}
	c.Code = BNT
	c.Name = "Bancor"
	c.Website = "https://www.bancor.network/"
	c.Explorer = "https://etherscan.io/token/Bancor"

	coinmap[c.Code] = c
}

func btg() {
	c := &Coin{}
	c.Code = BTG
	c.Name = "Bitcoin Gold"
	c.Website = "https://bitcoingold.org/"
	c.Explorer = "https://explorer.bitcoingold.org/insight/"

	coinmap[c.Code] = c
}

func bts() {
	c := &Coin{}
	c.Code = BTS
	c.Name = "BitShares"
	c.Website = "https://bitshares.org/"
	c.Explorer = "http://cryptofresh.com/"

	coinmap[c.Code] = c
}

func chat() {
	c := &Coin{}
	c.Code = CHAT
	c.Name = "ChatCoin"
	c.Website = "http://www.openchat.co/"
	c.Explorer = "https://etherscan.io/token/0x442bc47357919446eabc18c7211e57a13d983469"

	coinmap[c.Code] = c
}

func cloak() {
	c := &Coin{}
	c.Code = CLOAK
	c.Name = "CloakCoin"
	c.Website = "https://www.cloakcoin.com/"
	c.Explorer = "https://chainz.cryptoid.info/cloak/"

	coinmap[c.Code] = c
}

func cmt() {
	c := &Coin{}
	c.Code = CMT
	c.Name = "CyberMiles"
	c.Website = "https://www.cybermiles.io/"
	c.Explorer = "https://etherscan.io/token/0xf85feea2fdd81d51177f6b8f35f0e6734ce45f5f"

	coinmap[c.Code] = c
}

func cvc() {
	c := &Coin{}
	c.Code = CVC
	c.Name = "Civic"
	c.Website = "https://www.civic.com/"
	c.Explorer = "https://etherscan.io/token/civic"

	coinmap[c.Code] = c
}
