package main

import (
	"github.com/boltdb/bolt"
	"log"
)

/**
	区块链
 */

//定义区块链结构
type BlockChain struct {
	//blocks [] *Block

	//使用boltdb代替数组
	db *bolt.DB //操作数据库的句柄
	//“lastHashKey” 这个key对应的值，这个值就是最后一个区块的哈希值，用于新区块的创建添加。
	tail []byte //尾巴，存储最后一个区块的哈希
}

//创建区块链
func NewBlockChain() *BlockChain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	/*
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
	*/

	//打开数据库
	db, err := bolt.Open(blockChainDb, 0600, nil)
	if err != nil {
		log.Panic("bolt.Open err:", err)
	}
	//defer db.Close() //千万不要关闭，后面AddBlock要用db
	var tail []byte //最后一个区块的哈希

	//获取bucket
	db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBecket))

		if bucket == nil {
			bucket, err = tx.CreateBucket([]byte(blockBecket))
			if err != nil {
				log.Panic("CreateBucket err:", err)
			}

			//创建一个创世块，并作为第一个区块添加到区块链中
			gBlock := GenesisBlock()
			bucket.Put(gBlock.Hash, gBlock.Serialize())  //写创世块
			bucket.Put([]byte(lastHashKey), gBlock.Hash) //写最后一个区块的哈希
			tail = gBlock.Hash
		} else {
			tail = bucket.Get([]byte(lastHashKey))
		}

		return nil
	})

	return &BlockChain{db, tail}
}

//定义创世块
func GenesisBlock() *Block {
	block := NewBlock("Genesis Block", []byte{})
	return block
}

//添加区块
func (bc *BlockChain) AddBlock(data string) {
	/*
	lastBlock := bc.blocks[len(bc.blocks)-1] //前区块
	block := NewBlock(data, lastBlock.Hash)  //创建新区块
	bc.blocks = append(bc.blocks, block)     //添加新区块
	*/

	block := NewBlock(data, bc.tail) //创建新区块
	bc.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBecket))

		if bucket == nil {
			log.Panic("bucket should not be nil !")
		} else {
			bucket.Put(block.Hash, block.Serialize())   //写创世块
			bucket.Put([]byte(lastHashKey), block.Hash) //写最后一个区块的哈希
			bc.tail = block.Hash
		}

		return nil
	})
}
