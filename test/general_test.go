package test

import (
	"bufio"
	"crypto/ecdsa"
	"crypto/elliptic"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/altair-lab/xoreum/crypto"
)

func SavePrivateKey(addr string, priv *ecdsa.PrivateKey, filePath string) {
	D := priv.D.String()
	X := priv.PublicKey.X.String()
	Y := priv.PublicKey.Y.String()

	content := addr + " " + D + " " + X + " " + Y + "\n"
	//content_byte := []byte(content)

	file, _ := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	file.WriteString(content)
	file.Close()

	//ioutil.WriteFile(filePath, content_byte, 0644)
}

func LoadPrivateKeys(filePath string) map[string]*ecdsa.PrivateKey {
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	users := make(map[string]*ecdsa.PrivateKey)

	curve := elliptic.P256()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		contents := strings.Fields(line)

		priv := ecdsa.PrivateKey{}
		priv.Curve = curve

		addr := contents[0]

		d := new(big.Int)
		d, _ = d.SetString(contents[1], 10)
		priv.D = d

		x := new(big.Int)
		x, _ = x.SetString(contents[2], 10)
		priv.PublicKey.X = x

		y := new(big.Int)
		y, _ = y.SetString(contents[3], 10)
		priv.PublicKey.Y = y

		users[addr] = &priv
	}

	return users
}

func ExampleFunc() {

	// save 3 *big.int
	priv1, _ := crypto.GenerateKey()
	priv2, _ := crypto.GenerateKey()

	fmt.Println(priv1)
	fmt.Println(priv2)

	SavePrivateKey("address1", priv1, "privatekeys.txt")
	SavePrivateKey("address2", priv2, "privatekeys.txt")
	users := LoadPrivateKeys("privatekeys.txt")

	i := 0
	for k, v := range users {
		fmt.Println("address", i, ":", k)
		i++
		fmt.Println("\t", v)
	}

	// output: 1
}
