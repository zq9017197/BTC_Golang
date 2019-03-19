package main

import "fmt"

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
		fmt.Printf("TimeStamp : %d\n", block.TimeStamp)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)
	}
}
