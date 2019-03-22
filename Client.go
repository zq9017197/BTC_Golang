package main

import (
	"os"
	"fmt"
	"time"
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
	case "addBlock":
		if len(list) > 3 && list[2] == "--data" {
			data := list[3]
			if data == "" {
				fmt.Println("data should not be empty!")
				os.Exit(1)
			}
			cli.addBlock(data)
		}
	case "printChain":
		cli.printChain()
	default:
		fmt.Println(Usage)
	}

}

//添加区块（挖矿）
func (cli *Client) addBlock(data string) {
	cli.bc.AddBlock(data)
}

//打印区块链
func (cli *Client) printChain() {
	it := NewBlockChainIterator(cli.bc)
	for {
		block := it.GetBlockAndMoveLeft()
		fmt.Printf("===========================\n\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(&block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

		//终止条件
		if len(block.PreHash) == 0 {
			break
		}
	}
}
