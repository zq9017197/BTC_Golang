package main

import (
	"bytes"
	"encoding/binary"
	"log"
	"encoding/gob"
)

//将uint64转成[]byte
func uint64ToByte(num uint64) []byte {
	var buffer bytes.Buffer
	err := binary.Write(&buffer, binary.BigEndian, num)
	if err != nil {
		log.Panic("binary.Write err:", err)
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
