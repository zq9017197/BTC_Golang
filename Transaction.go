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

	//解锁脚本，我们用地址来模拟（签名，公钥）
	//ScriptSig string

	//真正的数字签名，由r，s拼成的[]byte
	Signature []byte

	//约定，这里的PubKey不存储原始的公钥，而是存储X和Y拼接的字符串，在校验端重新拆分（参考r,s传递）
	PubKey []byte //注意，是公钥，不是哈希，也不是地址
}

//交易输出结构
type TXOutput struct {
	Value float64 //接收金额

	//对方公钥的哈希，这个哈希可以通过地址反推出来，所以转账时知道地址即可！
	//ScriptPubKey string //锁定脚本,我们用地址模拟

	//收款方的公钥的哈希，注意，是哈希而不是公钥，也不是地址
	PubKeyHash []byte
}

//由于现在存储的字段是地址的公钥哈希，所以无法直接创建TXOutput，
//为了能够得到公钥哈希，我们需要处理一下，写一个Lock函数
func (output *TXOutput) Lock(address string) {
	//真正的锁定动作！！！！！
	output.PubKeyHash = GetPubKeyFromAddress(address)
}

//给TXOutput提供一个创建的方法，否则无法调用Lock
func NewTXOutput(value float64, address string) *TXOutput {
	output := TXOutput{
		Value: value,
	}

	output.Lock(address)
	return &output
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
	//coinbase总是新区块的第一条交易，这条交易中只有一个输出，即对矿工的奖励，没有输入。
	//address 是矿⼯地址，data是矿⼯自定义的附加信息
	if data == "" {
		data = fmt.Sprintf("reward %s %f\n", address, reward)
	}

	//比特币系统，对于这个input的id填0，对索引填0xffff，data由矿⼯填写，一般填所在矿池的名字
	input := TXInput{nil, -1, nil, []byte(data)}
	//output := TXOutput{reward, address}
	output := NewTXOutput(reward, address)

	tx := Transaction{nil, []TXInput{input}, []TXOutput{*output}}
	tx.SetHash()

	return &tx
}

//判断是否为挖矿交易
func (tx *Transaction) IsCoinbase() bool {
	if len(tx.TXInputs) == 1 && tx.TXInputs[0].PreTXID == nil && tx.TXInputs[0].VoutIndex == -1 {
		return true
	}

	return false
}

//创建普通交易
func NewTransaction(fromAddr string, toAddr string, amount float64, bc *BlockChain) *Transaction {
	//1. 创建交易之后要进行数字签名->所以需要私钥->打开钱包"NewWallets()"
	wallets := NewWallets()

	//2. 找到自己的钱包，根据地址返回自己的wallet
	wallet := wallets.WalletsMap[fromAddr]
	if wallet == nil {
		fmt.Printf("没有找到该地址的钱包，交易创建失败!\n")
		return nil
	}

	//3. 得到对应的公钥，私钥
	pubKey := wallet.PubKey
	//privateKey := wallet.Private

	//传递公钥的哈希，而不是传递地址
	pubKeyHash := HashPubKey(pubKey)

	//1.找到最合理的utxo集合 map[string][]int64
	utxos, calc := bc.FindNeedUTXOs(pubKeyHash, amount)
	if calc < amount {
		fmt.Println("余额不足，交易失败！")
		return nil
	}

	var inputs [] TXInput
	var outputs [] TXOutput

	//2.将这些utxo逐一转成inputs
	for txid, idxArr := range utxos {
		for _, idx := range idxArr {
			input := TXInput{[]byte(txid), int64(idx), nil, pubKey}
			inputs = append(inputs, input)
		}
	}

	//3.创建outputs
	//output := TXOutput{amount, toAddr}
	output := NewTXOutput(amount, toAddr)
	outputs = append(outputs, *output)

	//4.判断是否需要找零
	if calc > amount {
		output = NewTXOutput(calc-amount, fromAddr)
		outputs = append(outputs, *output)
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetHash()
	return &tx
}

/*
//解锁脚本，付款人会使用付款人的解锁脚本解开能够支配的UTXO
func (input *TXInput) CanUnlockUTXOWith(unlockData string) bool {
	//解锁脚本是检验input是否可以使用由某个地址锁定的utxo，所以对于解锁脚本来说，是外部提供锁定信息，我去检查一下能否解开它。
	//我们没有涉及到真实的非对称加密，所以使用字符串来代替加密和签名数据。即使用地址进行加密，同时使用地址当做签名，通过对比字符串来确定utxo能否解开。
	//ScriptSig是签名，v4版本中使用付款人的地址填充。unlockData是收款人的地址
	return input.ScriptSig == unlockData
}

//锁定脚本，使用收款人的地址对付款金额进行锁定
func (output *TXOutput) CanBeUnlockedWith(unlockData string) bool {
	//锁定脚本是用于指定比特币的新主人。在创建output时，应该是一直在等待一个签名的到来，检查这个签名能否解开自己锁定的比特币。
	//ScriptPubKey是锁定信息，v4版本中使用收款人的地址填充。unlockData是付款人的地址（签名）
	return output.ScriptPubKey == unlockData
}
*/
