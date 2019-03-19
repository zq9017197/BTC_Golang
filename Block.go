package main

import (
	"crypto/sha256"
	"time"
	"bytes"
	"encoding/binary"
	"log"
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
	Data [] byte //区块数据
}

//创建区块
func NewBlock(data string, preHash []byte) *Block {
	block := Block{
		Version:    00,
		PreHash:    preHash,
		MerKleRoot: []byte{}, //先填空，后面再计算
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 100,
		Nonce:      100,
		Hash:       []byte{}, //先填空，后面再计算
		Data:       []byte(data),
	}
	block.SetHash() //生成哈希值
	return &block
}

//生成哈希值
func (block *Block) SetHash() {
	//存储拼接好的数据，最后作为sha256函数的参数
	var blockInfo []byte
	blockInfo = append(blockInfo, block.PreHash...)
	blockInfo = append(blockInfo, block.Data...)
	blockInfo = append(blockInfo, block.MerKleRoot...)
	blockInfo = append(blockInfo, uint64ToByte(block.Version)...)
	blockInfo = append(blockInfo, uint64ToByte(block.TimeStamp)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Difficulty)...)
	blockInfo = append(blockInfo, uint64ToByte(block.Nonce)...)

	hash := sha256.Sum256(blockInfo) //生成hash
	block.Hash = hash[:]
}

//将uint64转成[]byte
func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}
	return buffer.Bytes()
}
