package main

/**
	常量
 */

const blockChainDb = "blockChain.db" //数据库文件名字

const blockBecket = "blockBecket" //bucket名字

const lastHashKey = "lastHashKey" //最后一个区块哈希的Key

/**
	命令行参数
	addBlock --data "Hello Btc" (废弃)
	addBlock --data DATA "add a block" (废弃)
	打印区块链：printChain
	获取“ing”余额：getBalance --address ing
	转账：send ing baibai 10 ing 挖矿收益

	测试案例：
	getBalance --address 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv
	getBalance --address 112rFUMPby5r15yZkYxWeZWMS5MLbd2p8o
	getBalance --address 12DA7Kv7jNdLz9zJF5UwNYgV3XQNeBgdrS
	getBalance --address 1919HLw9gDqdpB1hZtCKu6HBKVQESgvafq

	send 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 112rFUMPby5r15yZkYxWeZWMS5MLbd2p8o 7 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 挖矿收益
	send 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 12DA7Kv7jNdLz9zJF5UwNYgV3XQNeBgdrS 8 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 挖矿收益
	send 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 1919HLw9gDqdpB1hZtCKu6HBKVQESgvafq 9 1BFNQ79MCgEQXxVNj3xqjm1fo2ipuDrtfv 挖矿收益
 */
const Usage = `
	printChain "反向打印区块链"
	getBalance --address ADDRESS "获取指定地址的余额"
	send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"
	newWallet "创建一个新的钱包(私钥公钥对)"
	listAddresses "列举所有的钱包地址"
`

const usageSend = `send FROM TO AMOUNT MINER DATA "由FROM转AMOUNT给TO，由MINER挖矿，同时写入DATA"`
const usagegetBalance = `getBalance --address ADDRESS "获取指定地址的余额"`

//挖矿奖励
const reward = 12.5

//创世块中保存的信息
const genesisInfo = "Genesis Block"

//钱包数据文件名
const walletsFile = "wallets.dat"
