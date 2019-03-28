package main

func main() {
	bc := NewBlockChain("ing")
	//bc.AddBlock("Hello BTC")
	//bc.AddBlock("Hello ETH")
	//bc.AddBlock("Hello LTC")
	defer bc.db.Close()
	cli := Client{bc}
	cli.Run()

	/*
	for idx, block := range bc.blocks {
		fmt.Println("============== current block index :", idx)
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())
	}
	*/

	/*
	//bolt内部使用key的大小进行自动排序，而不是按照插入顺序排序
	bc.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(blockBecket))

		bucket.ForEach(func(k, v []byte) error {
			if bytes.Equal(k, []byte(lastHashKey)) {
				return nil
			}

			block := Deserialize(v) //反序列化v
			fmt.Printf("===========================\n\n")
			fmt.Printf("Version : %d\n", block.Version)
			fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
			fmt.Printf("Hash : %x\n", block.Hash)
			fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
			timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("TimeStamp : %s\n", timeFormat)
			fmt.Printf("Difficuty : %d\n", block.Difficulty)
			fmt.Printf("Nonce : %d\n", block.Nonce)
			fmt.Printf("Data : %s\n", block.Data)

			pow := NewProofOfWork(&block)
			fmt.Printf("IsValid : %v\n", pow.IsValid())
			return nil
		})

		return nil
	})
	*/

	/*
	it := NewBlockChainIterator(bc)
	for {
		block := it.GetBlockAndMoveLeft()
		fmt.Printf("===========================\n\n")
		fmt.Printf("Version : %d\n", block.Version)
		fmt.Printf("PrevBlockHash : %x\n", block.PreHash)
		fmt.Printf("Hash : %x\n", block.Hash)
		fmt.Printf("MerkleRoot : %x\n", block.MerKleRoot)
		timeFormat := time.Unix(int64(block.TimeStamp), 0).Format("2006-01-02 15:04:05")
		fmt.Printf("TimeStamp : %s\n", timeFormat)
		fmt.Printf("Difficuty : %d\n", block.Difficulty)
		fmt.Printf("Nonce : %d\n", block.Nonce)
		fmt.Printf("Data : %s\n", block.Data)

		pow := NewProofOfWork(&block)
		fmt.Printf("IsValid : %v\n", pow.IsValid())

		//终止条件
		if len(block.PreHash) == 0 {
			break
		}
	}
	*/

}
