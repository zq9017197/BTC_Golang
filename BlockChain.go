package main

/**
	区块链
 */

//定义区块链结构
type BlockChain struct {
	blocks [] *Block
}

//创建区块链
func NewBlockChain() *BlockChain {
	//创建一个创世块，并作为第一个区块添加到区块链中
	genesisBlock := GenesisBlock()
	return &BlockChain{
		blocks: []*Block{genesisBlock},
	}
}

//定义创世块
func GenesisBlock() *Block {
	block := NewBlock("Genesis Block", []byte{})
	return block
}

//添加区块
func (bc *BlockChain) AddBlock(data string) {
	lastBlock := bc.blocks[len(bc.blocks)-1] //前区块
	block := NewBlock(data, lastBlock.Hash)  //创建新区块
	bc.blocks = append(bc.blocks, block)     //添加新区块
}
