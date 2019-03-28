package main

import (
	"bytes"
	"encoding/gob"
	"log"
	"crypto/sha256"
	"fmt"
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
	PreTXID   []byte //引用utxo所在交易的ID
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

//创建挖矿交易
func NewCoinbaseTX(address string, data string) *Transaction {
	//address 是矿⼯地址，data是矿⼯自定义的附加信息
	if data == "" {
		data = fmt.Sprintf("reward %s %f\n", address, reward)
	}

	//比特币系统，对于这个input的id填0，对索引填0xffff，data由矿⼯填写，一般填所在矿池的名字
	input := TXInput{nil, -1, data}
	output := TXOutput{reward, address}

	tx := Transaction{nil, []TXInput{input}, []TXOutput{output}}
	tx.SetHash()

	return &tx
}
