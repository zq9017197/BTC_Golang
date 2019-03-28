package main

import (
	"time"
	"bytes"
	"log"
	"encoding/gob"
)

/**
	区块
 */

//定义区块结构
type Block struct {
	Version    uint64  //版本号
	PreHash    [] byte //前区块哈希值
	MerKleRoot []byte  //梅克尔根(就是一个哈希值)
	TimeStamp  uint64  //时间戳
	Difficulty uint64  //难度值(调整比特币挖矿的难度)
	Nonce      uint64  //随机数，这就是挖矿时所要寻找的数
	//正常比特币中没有当前区块的哈希值
	Hash [] byte //当前区块哈希值(为了方便实现，所以将区块的哈希值放到了区块中)
	//Data [] byte //区块数据
	Transactions []*Transaction //区块数据，真实交易数组
}

//创建区块
func NewBlock(txs []*Transaction, preHash []byte) *Block {
	block := Block{
		Version: 00,
		PreHash: preHash,
		//MerKleRoot: []byte{}, //先填空，后面再计算
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 20, //前面4个零(00001)
		//Nonce:      100,
		//Hash:       []byte{}, //先填空，后面再计算
		//Data: []byte(data),
		Transactions: txs,
	}

	block.MerKleRoot = block.MakeMerkelRoot() //设置梅克尔根

	//block.SetHash() //生成哈希值(v1)
	pow := NewProofOfWork(&block)
	hash, nonce := pow.Run() //POW 挖矿,计算符合hash的随机数
	block.Hash = hash
	block.Nonce = nonce

	return &block
}

/*
//生成哈希值
func (block *Block) SetHash() {
	//存储拼接好的数据，最后作为sha256函数的参数
	var blockInfo []byte
	*//*
	blockInfo = append(blockInfo, block.PreHash...)
	blockInfo = append(blockInfo, block.Data...)
	blockInfo = append(blockInfo, block.MerKleRoot...)
	blockInfo = append(blockInfo, uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Nonce)...)
	*//*
	tmp := [][]byte{
		block.PreHash,
		block.Data,
		block.MerKleRoot,
		uint64ToByte(block.Version),
		uint64ToByte(block.TimeStamp),
		uint64ToByte(block.Difficulty),
		uint64ToByte(block.Nonce),
	}
	blockInfo = bytes.Join(tmp, []byte(""))

	hash := sha256.Sum256(blockInfo) //生成hash
	block.Hash = hash[:]
}
*/

//序列化(转[]byte)
func (block *Block) Serialize() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic("encoder Encode err:", err)
	}

	return buffer.Bytes()
}

//模拟梅克尔根
func (block *Block) MakeMerkelRoot() []byte {
	//这里就不把所有交易数据两两哈希了，直接把所有 TXID连接起来
	//TODO
	return []byte{}
}
