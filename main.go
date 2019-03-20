package main

import (
	"fmt"
	"time"
)

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Hello BTC")
	bc.AddBlock("Hello ETH")

	for idx, block := range bc.blocks {
		fmt.Println(" ============== current block index :", idx)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())
	}

}
