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
func NewBlockChain(address string) *BlockChain {
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
			gBlock := GenesisBlock(address)
			//fmt.Printf("Genesis Block：%s\n",gBlock)

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
func GenesisBlock(address string) *Block {
	coinbase := NewCoinbaseTX(address, genesisInfo)
	block := NewBlock([]*Transaction{coinbase}, []byte{})
	return block
}

//添加区块
func (bc *BlockChain) AddBlock(txs []*Transaction) {
	/*
	lastBlock := bc.blocks[len(bc.blocks)-1] //前区块
	block := NewBlock(data, lastBlock.Hash)  //创建新区块
	bc.blocks = append(bc.blocks, block)     //添加新区块
	*/

	block := NewBlock(txs, bc.tail) //创建新区块
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

//返回指定地址能够支配的utxo的集合
func (bc *BlockChain) FindUTXOs(address string) []TXOutput {
	var utxos []TXOutput
	//保存消费过的output，key是这个output的txid，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)

	//1.遍历区块
	it := NewBlockChainIterator(bc)
	for {
		block := it.GetBlockAndMoveLeft()

		//2.遍历交易
		for _, tx := range block.Transactions {
			//3.遍历output，找到和自己相关的utxo（在添加output之前检查一下是否已经消耗过）
			for _, output := range tx.TXOutputs {
				if output.ScriptPubKey == address { //找到和自己相关的utxo
					utxos = append(utxos, output)
				}
			}

			//4.遍历input，找到自己花费过的utxo集合（把自己消费国的标识出来）
			for _, input := range tx.TXInputs {
				if input.ScriptSig == address {
					//交易输⼊，可能是多个。多个交易输入可能是同一个TXID，不同的索引
					idxArr := spentOutputs[string(input.PreTXID)]
					spentOutputs[string(input.PreTXID)] = append(idxArr, input.VoutIndex)
				}
			}
		}

		//终止条件
		if len(block.PreHash) == 0 {
			break
		}
	}

	//5.过滤已经消费过的utxo
	//TODO

	return utxos
}
