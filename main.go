package main

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/joho/godotenv"
)

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Index     int
	Timestamp string
	inventory  int
	Hash      string
	PrevHash  string
	Validator string
}

type percent struct {
	Per       int 	`json:"Percent"`
}

// Blockchain is a series of validated Blocks
var Blockchain []Block
var tempBlocks []Block

// candidateBlocks handles incoming blocks for validation
var candidateBlocks = make(chan Block)

// announcements broadcasts winning validator to all nodes
var announcements = make(chan string)

var mutex = &sync.Mutex{}

// validators keeps track of open validators and balances
var validators = make(map[string]int)

// validates keeps track of participating validators
var validates = make(map[string]int)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	
	

	// create genesis block
	t := time.Now()
	genesisBlock := Block{}
	genesisBlock = Block{0, t.String(), 0, calculateBlockHash(genesisBlock), "", ""}
	spew.Dump(genesisBlock)
	Blockchain = append(Blockchain, genesisBlock)

	tcpPort := os.Getenv("TCPPORT")
	httpPort := os.Getenv("HTTPPORT")

	go anotherConn(httpPort)

	// start TCP and serve TCP server
	server, err := net.Listen("tcp", ":"+tcpPort)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("TCP Server Listening on port :", tcpPort)
	defer server.Close()

	go func() {
		for candidate := range candidateBlocks {
			mutex.Lock()
			tempBlocks = append(tempBlocks, candidate)
			mutex.Unlock()
		}
	}()


	for {
		conn, err := server.Accept()
		if err != nil {
			log.Fatal(err)
		}
		
		go handleConn(conn)
	}

	
}

func getValidators(c *gin.Context) {
    c.IndentedJSON(http.StatusOK, validators)
}

func getActiveValidatorsPer(c *gin.Context) {

	var per = percent{
		Per: 0,
	}

	

	if len(validators) == 0 {
		per.Per = 0
	} else {
		per.Per = (len(validates)*100/len(validators))
    	// c.IndentedJSON(http.StatusOK, (len(validates)*100/len(validators)))
	}

	fmt.Println("ActiveValidatorsPercent is " , per)
	c.IndentedJSON(http.StatusOK, per)
}

func getValidatorPer(c *gin.Context) {

	var per = percent{
		Per: 0,
	}

	if len(validators) == 0 {
		per.Per = 0
	} else {
		per.Per = (100/len(validators))
	}

	fmt.Println("validatorsPercent is " , per)
	c.IndentedJSON(http.StatusOK, per)
}

func anotherConn(httpPort string) {

	router := gin.Default()
    router.GET("/getActiveValidatorsPer", getActiveValidatorsPer)
	router.GET("/getValidatorPer", getValidatorPer)

    router.Run("localhost:" + httpPort)
}


func handleConn(conn net.Conn) {
	defer conn.Close()

	go func() {
		for {
			msg := <-announcements
			io.WriteString(conn, msg)
		}
	}()
	// validator address
	var address string

	// allow user to allocate number of tokens to stake
	// the greater the number of tokens, the greater chance to forging a new block
	io.WriteString(conn, "Enter token balance:")
	scanBalance := bufio.NewScanner(conn)
	for scanBalance.Scan() {
		balance, err := strconv.Atoi(scanBalance.Text())
		if err != nil {
			log.Printf("%v not a number: %v", scanBalance.Text(), err)
			return
		}
		t := time.Now()
		address = calculateHash(t.String())
		validators[address] = balance
		fmt.Println(validators)
		break
	}

	io.WriteString(conn, "\nEnter a new inventory count:")

	scanInventory := bufio.NewScanner(conn)

	go func() {
		for {
			// take in inventory from stdin and add it to blockchain after conducting necessary validation
			for scanInventory.Scan() {
				inventory, err := strconv.Atoi(scanInventory.Text())
				// if malicious party tries to mutate the chain with a bad input, delete them as a validator and they lose their staked tokens
				if err != nil {
					log.Printf("%v not a number: %v", scanInventory.Text(), err)
					delete(validators, address)
					conn.Close()
				}

				mutex.Lock()
				oldLastIndex := Blockchain[len(Blockchain)-1]
				mutex.Unlock()

				// create newBlock for consideration to be forged
				newBlock, err := generateBlock(oldLastIndex, inventory, address)
				if err != nil {
					log.Println(err)
					continue
				}
				if isBlockValid(newBlock, oldLastIndex) {
					candidateBlocks <- newBlock
				}
				validates[address] = inventory
				fmt.Println(validators)
				
				io.WriteString(conn, "\nEnter a new inventory count:")
			}
		}
	}()

	// simulate receiving broadcast
	for {
		time.Sleep(time.Minute)
		mutex.Lock()
		output, err := json.Marshal(Blockchain)
		mutex.Unlock()
		if err != nil {
			log.Fatal(err)
		}
		io.WriteString(conn, string(output)+"\n")
	}

}

// isBlockValid makes sure block is valid by checking index
// and comparing the hash of the previous block
func isBlockValid(newBlock, oldBlock Block) bool {
	if oldBlock.Index+1 != newBlock.Index {
		return false
	}

	if oldBlock.Hash != newBlock.PrevHash {
		return false
	}

	if calculateBlockHash(newBlock) != newBlock.Hash {
		return false
	}

	return true
}

// SHA256 hasing
// calculateHash is a simple SHA256 hashing function
func calculateHash(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

//calculateBlockHash returns the hash of all block information
func calculateBlockHash(block Block) string {
	record := string(block.Index) + block.Timestamp + string(block.inventory) + block.PrevHash
	return calculateHash(record)
}

// generateBlock creates a new block using previous block's hash
func generateBlock(oldBlock Block, inventory int, address string) (Block, error) {

	var newBlock Block

	t := time.Now()

	newBlock.Index = oldBlock.Index + 1
	newBlock.Timestamp = t.String()
	newBlock.inventory = inventory
	newBlock.PrevHash = oldBlock.Hash
	newBlock.Hash = calculateBlockHash(newBlock)
	newBlock.Validator = address

	return newBlock, nil
}
