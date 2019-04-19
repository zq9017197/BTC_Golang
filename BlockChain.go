package main

import (
	"github.com/boltdb/bolt"
	"log"
	"bytes"
	"fmt"
	"errors"
	"crypto/ecdsa"
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

	//校验签名
	for _, tx := range txs {
		if !bc.VerifyTransaction(tx) {
			fmt.Printf("矿工发现无效交易!")
			return
		}
	}

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
func (bc *BlockChain) FindUTXOs(pubKeyHash []byte) []TXOutput {
	var utxos []TXOutput

	txs := bc.FindUTXOTransactions(pubKeyHash)
	for _, tx := range txs {
		for _, output := range tx.TXOutputs {
			if bytes.Equal(pubKeyHash, output.PubKeyHash) {
				utxos = append(utxos, output) //找到和自己相关的utxo
			}
		}
	}

	return utxos
}

//找到满足转账条件，未消费过的，合理的utxo的集合
func (bc *BlockChain) FindNeedUTXOs(fromPubKeyHash []byte, amount float64) (map[string][]int64, float64) {
	utxos := make(map[string][]int64) //找到的合理的utxo集合
	var calc float64                  //扎到的utxo里面包含的钱的总数

	txs := bc.FindUTXOTransactions(fromPubKeyHash)
	for _, tx := range txs {
		for idx, output := range tx.TXOutputs {
			if bytes.Equal(fromPubKeyHash, output.PubKeyHash) { //找到和自己相关的utxo
				if calc < amount {
					utxos[string(tx.TXID)] = append(utxos[string(tx.TXID)], int64(idx))
					calc += output.Value
					if calc >= amount {
						return utxos, calc
					}
				}
			}
		}
	}

	return utxos, calc
}

//返回指定地址的 Transaction集合
func (bc *BlockChain) FindUTXOTransactions(fromPubKeyHash []byte) []*Transaction {
	//存储所有包含utxo的交易
	var txs []*Transaction
	//保存消费过的output，key是这个output的txid，value是这个交易中索引的数组
	spentOutputs := make(map[string][]int64)

	//1.遍历区块
	it := NewBlockChainIterator(bc)
	for {
		block := it.GetBlockAndMoveLeft()

		//2.遍历交易
		for _, tx := range block.Transactions {
		OUTPUT:
		//3.遍历output，找到和自己相关的utxo（在添加output之前检查一下是否已经消耗过）
			for idx, output := range tx.TXOutputs {
				//5.过滤消费过的utxo
				arrs := spentOutputs[string(tx.TXID)]
				if arrs != nil {
					for _, index := range arrs {
						//当前准备添加的output已经消费过了，不要再添加了
						if int64(idx) == index {
							continue OUTPUT
						}
					}
				}

				if bytes.Equal(fromPubKeyHash, output.PubKeyHash) { //找到和自己相关的 Transaction
					txs = append(txs, tx)
					//同一个 Transaction添加一次就退出循环
					break //TODO
				}
			}

			//如果当前交易是挖矿交易的话，那么不做遍历，直接跳过
			if !tx.IsCoinbase() {
				//4.遍历input，找到自己花费过的utxo集合（把自己消费国的标识出来）
				for _, input := range tx.TXInputs {
					if bytes.Equal(fromPubKeyHash, HashPubKey(input.PubKey)) {
						//交易输入，可能是多个。多个交易输入可能是同一个TXID，不同的索引
						spentOutputs[string(input.PreTXID)] = append(spentOutputs[string(input.PreTXID)], input.VoutIndex)
					}
				}
			}
		}

		//终止条件
		if len(block.PreHash) == 0 {
			break
		}
	}

	return txs
}

//根据id查找交易本身，需要遍历整个区块链
func (bc *BlockChain) FindTransactionByTXid(id []byte) (Transaction, error) {
	it := NewBlockChainIterator(bc)
	//1. 遍历区块链
	for {
		block := it.GetBlockAndMoveLeft()
		//2. 遍历交易
		for _, tx := range block.Transactions {
			//3. 比较交易，找到了直接退出
			if bytes.Equal(tx.TXID, id) {
				return *tx, nil
			}
		}

		if len(block.PreHash) == 0 {
			fmt.Printf("区块链遍历结束!\n")
			break
		}
	}

	//4. 如果没找到，返回空Transaction，同时返回错误状态
	return Transaction{}, errors.New("无效的交易id，请检查!")
}

//交易签名
func (bc *BlockChain) SignTransaction(tx *Transaction, privateKey *ecdsa.PrivateKey) {
	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.PreTXID)
		if err != nil {
			log.Panic(err)
		}

		/*
		第一个input查找之后：prevTXs：
			map[2222]Transaction222

		第二个input查找之后：prevTXs：
			map[2222]Transaction222
			map[3333]Transaction333

		第三个input查找之后：prevTXs：
			map[2222]Transaction222
			map[3333]Transaction333(只不过是重新写了一次)
		*/
		prevTXs[string(input.PreTXID)] = tx

	}

	tx.Sign(privateKey, prevTXs)
}

//验证签名
func (bc *BlockChain) VerifyTransaction(tx *Transaction) bool {
	if tx.IsCoinbase() {
		return true
	}

	//签名，交易创建的最后进行签名
	prevTXs := make(map[string]Transaction)

	//找到所有引用的交易
	//1. 根据inputs来找，有多少input, 就遍历多少次
	//2. 找到目标交易，（根据TXid来找）
	//3. 添加到prevTXs里面
	for _, input := range tx.TXInputs {
		//根据id查找交易本身，需要遍历整个区块链
		tx, err := bc.FindTransactionByTXid(input.PreTXID)
		if err != nil {
			log.Panic(err)
		}

		prevTXs[string(input.PreTXID)] = tx
	}

	return tx.Verify(prevTXs)
}
