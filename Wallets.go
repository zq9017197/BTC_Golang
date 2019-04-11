package main

/**
	保存所有的wallet以及它的地址
 */

type Wallets struct {
	WalletsMap map[string]*Wallet //map[地址]钱包
}

//创建方法，返回当前所有钱包的实例
func NewWallets() *Wallets {
	wallet := NewWallet()
	address := wallet.NewAddress()

	var wallets Wallets
	wallets.WalletsMap = make(map[string]*Wallet)
	wallets.WalletsMap[address] = wallet
	return &wallets
}
