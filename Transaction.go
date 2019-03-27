package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
)

/**
	交易
 */

//交易结构
type Transaction struct {
	TXID      []byte      //交易ID
	TXInputs  [] TXInput  //交易输⼊，可能是多个
	TXOutputs [] TXOutput //交易输出，可能是多个
}

//交易输入结构
type TXInput struct {
	PreTXID      []byte //引用utxo所在交易的ID
	VoutIndex int64  //所消费utxo在output中的索引
	ScriptSig string //解锁脚本（签名，公钥）
}

//交易输出结构
type TXOutput struct {
	Value float64 //接收金额
	//对方公钥的哈希，这个哈希可以通过地址反推出来，所以转账时知道地址即可！
	ScriptPubKey string //锁定脚本
}

//设置TXID
func (tx *Transaction) SetHash() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	if err != nil {
		log.Panic("encoder Encode err:", err)
	}

	//先序列化再hash，难得拼字符串！
	data := buffer.Bytes()
	hash := sha256.Sum256(data)
	tx.TXID = hash[:]
}
