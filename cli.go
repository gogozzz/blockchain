package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	_ "strconv"
)

type CLI struct {
	//bc *Blockchain
}

func (cil *CLI) validateArgs() {
	if len(os.Args) < 2 {
		//cli.printUsage()
		fmt.Println("error : args < 2")
		os.Exit(1)
	}
}

func (cli *CLI) Run() {
	cli.validateArgs()

	addBlockCmd := flag.NewFlagSet("addBlock", flag.ExitOnError)
	printChainCmd := flag.NewFlagSet("printChain", flag.ExitOnError)
	createChainCmd := flag.NewFlagSet("createChain", flag.ExitOnError)
	getbalanceCmd := flag.NewFlagSet("getbalance", flag.ExitOnError)
	sendCmd := flag.NewFlagSet("send", flag.ExitOnError)

	blockData := addBlockCmd.String("data", "", "Block data")
	createChainAddress := createChainCmd.String("address", "", "Create chain address")
	getbalanceAddress := getbalanceCmd.String("address", "", "Get balance")

	sendFrom := sendCmd.String("from", "", "Source wallet address")
	sendTo := sendCmd.String("to", "", "Destination wallet address")
	sendAmount := sendCmd.Int("amount", 0, "Amount to send")

	var err error
	switch os.Args[1] {
	case "addBlock":
		err = addBlockCmd.Parse(os.Args[2:])
	case "printChain":
		err = printChainCmd.Parse(os.Args[2:])
	case "createChain":
		err = createChainCmd.Parse(os.Args[2:])
	case "getbalance":
		err = getbalanceCmd.Parse(os.Args[2:])
	case "send":
		err = sendCmd.Parse(os.Args[2:])
	default:
		os.Exit(1)
	}

	if err != nil {
		log.Panic(err)
	}

	if addBlockCmd.Parsed() {
		if *blockData == "" {
			addBlockCmd.Usage()
			os.Exit(1)
		}

		//cli.bc.AddBlock(*blockData)
	}

	if printChainCmd.Parsed() {

		cli.printChain()

	}

	if createChainCmd.Parsed() {
		cli.createChain(*createChainAddress)
	}

	if getbalanceCmd.Parsed() {
		cli.getBalance(*getbalanceAddress)
	}

	if sendCmd.Parsed() {
		if *sendFrom == "" || *sendTo == "" || *sendAmount <= 0 {
			sendCmd.Usage()
			os.Exit(1)
		}

		cli.send(*sendFrom, *sendTo, *sendAmount)
	}

}

func (cli *CLI) printChain() {

	//bci := cli.bc.Iterator()
	//
	//for {
	//	block := bci.Next()
	//
	//	fmt.Printf("Prev. hash: %x\n", block.PrevBlockHash)
	//	//fmt.Printf("Data: %s\n", block.Data)
	//	fmt.Printf("Hash: %x\n", block.Hash)
	//	//pow := NewProofOfWork(block)
	//	//fmt.Printf("PoW: %s\n", strconv.FormatBool(pow.Validate()))
	//	fmt.Println()
	//
	//	if len(block.PrevBlockHash) == 0 {
	//		break
	//	}
	//}

}

func (cli CLI) createChain(address string) {

	bc := NewBlockchain(address)

	defer bc.db.Close()

	fmt.Println("Create Chain Done!")
}

func (cli *CLI) getBalance(address string) {
	bc := NewBlockchain(address)
	defer bc.db.Close()

	balance := 0
	UTXOs := bc.FindUTXO(address)

	for _, out := range UTXOs {
		balance += out.Value
	}

	fmt.Printf("Balance of '%s': %d\n", address, balance)
}

func (cli *CLI) send(from, to string, amount int) {
	bc := NewBlockchain(from)
	defer bc.db.Close()

	tx := NewUTXOTransaction(from, to, amount, bc)
	bc.MineBlock([]*Transaction{tx})
	fmt.Println("Success!")
}
