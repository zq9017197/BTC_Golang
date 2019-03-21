package main

import (
	"github.com/boltdb/bolt"
	"fmt"
	"time"
)

func main() {
	bc := NewBlockChain()
	bc.AddBlock("Hello BTC")
	bc.AddBlock("Hello ETH")
	bc.AddBlock("Hello LTC")

	/*
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
	*/

	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBecket))

		cursor := bucket.Cursor() //遍历 key
		for hash, data := cursor.First(); hash != nil; hash, data = cursor.Next() {
			block := Deserialize(data) //反序列化
			fmt.Println(" ============== current block hash :", hash)
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

		return nil
	})

}
