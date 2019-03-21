package main

import (
	"time"
	"bytes"
	"encoding/binary"
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
	Data [] byte //区块数据
}

//创建区块
func NewBlock(data string, preHash []byte) *Block {
	block := Block{
		Version:    00,
		PreHash:    preHash,
		MerKleRoot: []byte{}, //先填空，后面再计算
		TimeStamp:  uint64(time.Now().Unix()),
		Difficulty: 20, //前面4个零(00001)
		//Nonce:      100,
		//Hash:       []byte{}, //先填空，后面再计算
		Data: []byte(data),
	}

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

//将uint64转成[]byte
func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic("binary.Write err:", err)
	}
	return buffer.Bytes()
}

//序列化(转[]byte)
func (block *Block) toByte() []byte {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	err := encoder.Encode(block)
	if err != nil {
		log.Panic("encoder Encode err:", err)
	}

	return buffer.Bytes()
}

//反序列化
//gob是Golang包自带的一个数据结构序列化的编码/解码工具。编码使用Encoder，解码使用Decoder。
func Deserialize(data []byte) *Block {
	var block Block
	var buffer bytes.Buffer

	_, err := buffer.Write(data)
	if err != nil {
		log.Panic("buffer.Read err:", err)
	}

	decoder := gob.NewDecoder(&buffer)
	err = decoder.Decode(&block)
	if err != nil {
		log.Panic("decoder.Decode err:", err)
	}

	return &block
}
