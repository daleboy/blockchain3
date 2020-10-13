package blockchain3

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

//CLI 响应处理命令行参数
type CLI struct {
	BC *Blockchain
}

func (cli *CLI) printUsage() {
	fmt.Println("Usage:")
	fmt.Println("   addblock -data BLOCK_DATA - 将一个区块加入到区块链")
	fmt.Println("   printchain - 打印区块链中的所有区块")
}

//validateArgs 校验命令，如果无效，打印使用说明
func (cli *CLI) validateArgs() {
	if len(os.Args) < 2 { //所有命令至少有两个参数
		cli.printUsage()
		os.Exit(1)
	}
}

// printChain 打印区块，从最新到最旧，直到打印完成创始区块
func (cli *CLI) printChain() {
	bci := cli.BC.Iterator()

	for {
		block := bci.Next()

		fmt.Printf("Prev. Hash:%x\n", block.PrevBlockHash)
		fmt.Printf("Data:%s\n", block.Data)
		fmt.Printf("Hash:%x\n", block.Hash)
		pow := NewProofOfWork(block)
		fmt.Printf("PoW:%s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()

		if len(block.PrevBlockHash) == 0 { //创始区块的PrevBlockHash为byte[]{}
			break
		}
	}
}

//addBlock 添加区块
func (cli *CLI) addBlock(data string) {
	cli.BC.AddBlock(data)
	os.Exit(1)
}

// Run 读取命令行参数，执行相应的命令
//使用标准库里面的 flag 包来解析命令行参数：
func (cli *CLI) Run() {
	cli.validateArgs()

	//定义名称为"addblock"的空的flagset集合
	addBlockCmd := flag.NewFlagSet("addblock", flag.ExitOnError)
	//定义名称为"printchain"的空的flagset集合
	printChainCmd := flag.NewFlagSet("printchain", flag.ExitOnError)

	//String用指定的名称给addBlockCmd新增一个字符串flag
	//以指针的形式返回addBlockCmd
	addBlockData := addBlockCmd.String("data", "", "Block data")

	//os.Args包含以程序名称开始的命令行参数
	switch os.Args[1] { //os.Args[0]为程序名称，真正传递的参数index从1开始，一般而言Args[1]为命令名称
	case "addblock":
		//Parse调用之前，必须保证addBlockCmd所有的flag都已经定义在其中
		//执行完毕，将修改addBlockData
		err := addBlockCmd.Parse(os.Args[2:]) //仅解析参数，不含命令
		if err != nil {
			log.Panic(err)
		}
	case "printchain":
		//Parse调用之前，必须保证addBlockCmd所有的flag都已经定义在其中
		//根据命令设计，这里将返回nil，所以在前面没有定义接收解析后数据的flag
		//但printChainCmd的parsed=true
		err := printChainCmd.Parse(os.Args[2:]) //仅仅解析参数，不含命令
		if err != nil {
			log.Panic(err)
		}
	default:
		cli.printUsage()
		os.Exit(1)
	}

	if addBlockCmd.Parsed() {
		if *addBlockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}
		cli.addBlock(*addBlockData)
	}

	if printChainCmd.Parsed() {
		cli.printChain()
	}
}
