package main

import (
	"github.com/boltdb/bolt"
	"log"
)

/**
	区块链迭代器
 */

//区块链迭代器的结构
type BlockChainIterator struct {
	db            *bolt.DB
	current_point []byte //指向当前的区块
}

//迭代器创建函数
func NewBlockChainIterator(bc *BlockChain) BlockChainIterator {
	return BlockChainIterator{
		db:            bc.db,
		current_point: bc.tail,
	}
}

//迭代器访问函数
func (it *BlockChainIterator) GetBlockAndMoveLeft() Block {
	var block Block

	it.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBecket))

		if bucket == nil {
			log.Panic("bucket should not be nil !")
		} else {
			//根据当前的current_pointer获取block
			data := bucket.Get(it.current_point)
			block = Deserialize(data)
			it.current_point = block.PreHash //将游标（指针）左移
		}

		return nil
	})

	return block
}
