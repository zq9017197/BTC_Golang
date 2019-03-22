package main

import (
	"math/big"
	"bytes"
	"crypto/sha256"
	"fmt"
)

/**
	POW 工作量证明
 */

//工作量证明的结构ProofOfWork
type ProofOfWork struct {
	block  *Block
	target *big.Int //挖矿目标值，先写成固定的值，后面再进行推导演算。
}

//创建POW的函数
func NewProofOfWork(block *Block) *ProofOfWork {
	pow := ProofOfWork{
		block: block,
	}

	//自定义的难度值，先写成固定值，最终挖矿的hash要比target值小
	/*
	targetString := "0000100000000000000000000000000000000000000000000000000000000000"
	bigIntTmp := big.Int{}
	bigIntTmp.SetString(targetString, 16) //将难度值赋值给big.Int，指定16进制的格式
	pow.target = &bigIntTmp
	*/

	/*
	目标值
		0000100000000000000000000000000000000000000000000000000000000000
	初始值
		0000000000000000000000000000000000000000000000000000000000000001
	左移256位(64 * 4)
		10000000000000000000000000000000000000000000000000000000000000000
	右移20位(5 * 4)
		0000100000000000000000000000000000000000000000000000000000000000
	*/
	targetLocal := big.NewInt(1)
	//targetLocal.Lsh(targetLocal, 256)
	//targetLocal.Rsh(targetLocal, uint(pow.block.Difficulty))
	targetLocal.Lsh(targetLocal, 256-uint(pow.block.Difficulty))
	pow.target = targetLocal

	return &pow
}

//不断计算hash的函数
func (pow *ProofOfWork) Run() (hash []byte, nonce uint64) {
	for {
		data := pow.prepareData(nonce)
		hashTmp := sha256.Sum256(data) //生成hash
		hash = hashTmp[:]              //计算的哈希值
		bigIntTmp := big.Int{}
		bigIntTmp.SetBytes(hash) //转换为big.Int

		if bigIntTmp.Cmp(pow.target) == -1 {
			fmt.Printf("mining success: %x, %d\n", hash, nonce)
			break //挖矿成功
		}
		nonce++
	}

	return hash, nonce
}

//辅助函数-类似于v1中的setHash函数的功能
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	block := pow.block
	tmp := [][]byte{
		block.PreHash,
		block.Data,
		block.MerKleRoot,
		uint64ToByte(block.Version),
		uint64ToByte(block.TimeStamp),
		uint64ToByte(block.Difficulty),
		uint64ToByte(nonce),
	}
	data := bytes.Join(tmp, []byte(""))
	return data
}

//校验函数-即对求出来的哈希和随机数进行验证，只需要对求出来的值进行反向的计算比较即可。
func (pow *ProofOfWork) IsValid() bool {
	hash := sha256.Sum256(pow.prepareData(pow.block.Nonce))
	fmt.Printf("is valid hash : %x, %d\n", hash[:], pow.block.Nonce)

	tmp := big.Int{}
	tmp.SetBytes(hash[:])
	if tmp.Cmp(pow.target) == -1 {
		return true
	}
	return false
}
