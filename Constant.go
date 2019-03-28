package main

/**
	常量
 */

const blockChainDb = "blockChain.db" //数据库文件名字

const blockBecket = "blockBecket" //bucket名字

const lastHashKey = "lastHashKey" //最后一个区块哈希的Key

//命令行参数
const Usage = `
	addBlock --data DATA "add a block"
	printChain "print block Chain"
`

//挖矿奖励
const reward = 12.5
