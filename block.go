package blockchain3

import (
	"bytes"
	"encoding/gob"
	"log"
	"time"
)

//Block 区块结构新版，增加了计数器nonce，主要目的是为了校验区块是否合法
//即挖出的区块是否满足工作量证明要求的条件
type Block struct {
	Timestamp     int64
	Data          []byte
	PrevBlockHash []byte
	Nonce         int
	Hash          []byte
}

//NewBlock 创建普通区块
func NewBlock(data string, prevBlockHash []byte) *Block {
	block := &Block{time.Now().Unix(), []byte(data), prevBlockHash, 0, []byte{}}

	//挖矿实质上是算出符合要求的哈希
	pow := NewProofOfWork(block) //注意传递block指针作为参数
	nonce, hash := pow.Run()

	//设置block的计数器和哈希
	block.Nonce = nonce
	block.Hash = hash[:]

	return block
}

//NewGenesisBlock 创建创始区块。注意，创建创始区块也需要挖矿。
func NewGenesisBlock() *Block {
	return NewBlock("创始区块", []byte{})
}

//Serialize Block序列化
func (b *Block) Serialize() []byte {
	var result bytes.Buffer //定义一个buffer存储序列化后的数据

	//初始化一个encoder，gob是标准库的一部分
	//encoder根据参数的类型来创建，这里将编码为字节数组
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b) //编码
	if err != nil {
		log.Panic(err) //如果出错，将记录log后，Panic调用，立即终止当前函数的执行
	}

	return result.Bytes()
}

// DeserializeBlock 反序列化，注意返回的是Block的指针（引用）
func DeserializeBlock(d []byte) *Block {
	var block Block //一般都不会通过指针来创建一个struct。记住struct是一个值类型

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)
	if err != nil {
		log.Panic(err)
	}

	return &block //返回block的引用
}
