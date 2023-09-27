package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/ethereum/go-ethereum/crypto"
)

var (
	mu      sync.Mutex
	bestMax int
)

func countLeadingSymbols(str string, symbol rune) int {
	count := 0
	for _, s := range str {
		if s == symbol {
			count++
		} else {
			break
		}
	}
	return count
}

func writeToFile(n int, privateKeyHex string, addressHex string) {
	file, err := os.OpenFile("best_wallets.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	result := fmt.Sprintf("| N = %d ", n) + fmt.Sprintf("| Address | 0x%s ", addressHex) + fmt.Sprintf("| Private | %s", privateKeyHex) + "\n"

	_, err = file.WriteString(result)
	if err != nil {
		log.Fatal(err)
	}
}

func genLoop(wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		privateKey, err := crypto.GenerateKey()
		if err != nil {
			log.Fatal(err)
		}

		address := crypto.PubkeyToAddress(privateKey.PublicKey)
		privateKeyHex := fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
		addressHex := fmt.Sprintf("%x", address)

		n := countLeadingSymbols(addressHex, '0')

		mu.Lock()
		if n >= bestMax {
			bestMax = n
			writeToFile(n, privateKeyHex, addressHex)
			log.Print(fmt.Sprintf("| N = %d ", n), fmt.Sprintf("| Address | 0x%s ", addressHex), fmt.Sprintf("| Private | %s", privateKeyHex), "\r")
		}
		mu.Unlock()
	}
}

func main() {
	numCPU := runtime.NumCPU()
	runtime.GOMAXPROCS(numCPU)
	log.SetFlags(log.LstdFlags ^ log.Ldate)

	var wg sync.WaitGroup

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		go genLoop(&wg)
	}
	wg.Wait()
}
