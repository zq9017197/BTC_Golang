package main

import "fmt"

func main() {
	bc := NewBlockChain()
	bc.AddBlock("老师转班长100枚比特币")
	bc.AddBlock("班长转我500枚比特币")

	for idx, block := range bc.blocks {
		fmt.Printf("======当前区块高度：%d======\n", idx)
		fmt.Printf("版本号：%d\n", block.Version)
		fmt.Printf("前区块哈希值：%x\n", block.PreHash)
		fmt.Printf("梅克尔根：%x\n", block.MerKleRoot)
		fmt.Printf("时间戳：%d\n", block.TimeStamp)
		fmt.Printf("难度值：%d\n", block.Difficulty)
		fmt.Printf("随机数：%d\n", block.Nonce)
		fmt.Printf("当前区块哈希值：%x\n", block.Hash)
		fmt.Printf("区块数据：%x\n", block.Data)
	}
}
