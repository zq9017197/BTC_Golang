package main

import (
	"os"
	"fmt"
	"time"
	"strconv"
)

/**
	命令行客户端
*/

//客户端结构体
type Client struct {
	bc *BlockChain
}

//接收命令行参数
func (cli *Client) Run() {
	list := os.Args
	if len(list) < 2 {
		fmt.Println(Usage)
		os.Exit(1)
	}

	cmd := list[1]
	switch cmd {
	/*
	case "addBlock":
		if len(list) == 4 && list[2] == "--data" {
			data := list[3]
			if data == "" {
				fmt.Println("data should not be empty!")
				os.Exit(1)
			}
			cli.addBlock(nil) //TODO
		}
	*/
	case "printChain":
		cli.printChain()
	case "newWallet":
		cli.NewWallet()
	case "getBalance":
		if len(list) == 4 && list[2] == "--address" {
			address := list[3]
			if address == "" {
				fmt.Println("address should not be empty!")
				os.Exit(1)
			}
			cli.getBalance(address)
		}
	case "send":
		if len(list) != 7 {
			fmt.Println(usageSend)
			os.Exit(1)
		}

		fromAddr := list[2]
		toAddr := list[3]
		amount, _ := strconv.ParseFloat(list[4], 64)
		miner := list[5]
		data := list[6]
		cli.send(fromAddr, toAddr, amount, miner, data)
	default:
		fmt.Println(Usage)
	}

}

//添加区块（挖矿）
func (cli *Client) addBlock(txs []*Transaction) {
	cli.bc.AddBlock(txs)
}

//打印区块链
func (cli *Client) printChain() {
	it := NewBlockChainIterator(cli.bc)
	for {
		block := it.GetBlockAndMoveLeft()
		fmt.Printf("===========================\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Transactions[0].TXInputs[0].ScriptSig)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

		//终止条件
		if len(block.PreHash) == 0 {
			break
		}
	}
}

//获取余额
func (cli *Client) getBalance(address string) {
	utxos := cli.bc.FindUTXOs(address)

	var total float64
	for _, utxo := range utxos {
		total += utxo.Value
	}

	fmt.Printf("The balance of \"%s\" is : %f BTC\n", address, total)
}

//转账交易
func (cli *Client) send(fromAddr, toAddr string, amount float64, miner, data string) {
	//创建挖矿交易
	coinbase := NewCoinbaseTX(miner, data)

	//交接普通交易
	tx := NewTransaction(fromAddr, toAddr, amount, cli.bc)

	//添加区块
	if tx != nil {
		cli.bc.AddBlock([]*Transaction{coinbase, tx})
		fmt.Println("Send Successfully!")
	}
}

//创建一个新的钱包
func (cli *Client) NewWallet() {
	//wallet := NewWallet()
	//address := wallet.NewAddress()

	wallets := NewWallets()
	for address := range wallets.WalletsMap {
		fmt.Printf("地址：%s\n", address)
	}
}
