package main

import (
	"bufio"
	"container/list"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"slices"
	"testing"

	"github.com/go-resty/resty/v2"
)

// https://stackoverflow.com/questions/8757389/reading-a-file-line-by-line-in-go
func Test(t *testing.T) {
	file, err := os.Open("headers.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	fixFile, err := os.Create("headers.txt.fixed")
	if err != nil {
		log.Fatal(err)
	}
	defer fixFile.Close()
	existingHeaders := list.New()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		nextHeader := scanner.Text()
		// existingHeaders = append(existingHeaders, nextHeader)
		existingHeaders.PushBack(nextHeader)
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	lastHashHex := "0000000000000000000000000000000000000000000000000000000000000000"
	lastHash, _ := decodeHex(lastHashHex)
	lastLine := ""
	for e := existingHeaders.Front(); e != nil; e = e.Next() {
		nextHeader := e.Value.(string)
		fixedValue, isLastLineValid := testAndFixHash(lastHash, nextHeader, e)
		if isLastLineValid && lastLine != "" { // first line is always valid
			fixFile.WriteString(lastLine)
			fixFile.WriteString("\n")
		}
		for i := 0; i < len(fixedValue); i++ {
			fixFile.WriteString(fixedValue[i])
			fixFile.WriteString("\n")
		}
		lastLine = nextHeader
		lastHash = doubleSha256(lastLine)
	}

	fixFile.WriteString(lastLine) // no value to validate against, assumed valid unless found otherwise by new data
	fixFile.WriteString("\n")
}

func decodeHex(h string) ([]byte, error) {
	ret := make([]byte, len(h)/2)
	count, err := hex.Decode(ret, []byte(h))
	if err != nil {
		return nil, err
	}
	if count*2 != len(h) {
		return nil, fmt.Errorf("lenth mismatch")
	}
	return ret, nil
}

func testAndFixHash(lastHash []byte, nextHeader string, e *list.Element) ([]string, bool) {
	inHeaderHashHex := nextHeader[4*2 : 36*2]
	inHeaderHashBytes, err := decodeHex(inHeaderHashHex)
	if err != nil {
		panic(err)
	}
	if len(lastHash) != len(inHeaderHashBytes) {
		panic(fmt.Errorf("length mismatch"))
	}

	if !testHashes(lastHash, inHeaderHashBytes) {
		fmt.Printf("downloading header for %s\n", inHeaderHashHex)
		// download header with this hash value
		slices.Reverse[[]byte]([]byte(inHeaderHashBytes))
		reversedHeaderHex := hex.EncodeToString(inHeaderHashBytes)
		url := fmt.Sprintf("https://blockchain.info/rawblock/%s?format=hex", reversedHeaderHex)

		client := resty.New()
		resp, err := client.R().Get(url)

		if err != nil && resp.StatusCode() != 200 {
			fmt.Println("cannot fix error, move to next line")
			return []string{nextHeader}, false
		}

		newLineHex := resp.Body()[:160]
		newLineBytes, err := decodeHex(string(newLineHex))
		if err != nil {
			fmt.Println("cannot fix error, move to next line")
			return []string{nextHeader}, false
		}

		// if newLine bridges last header and current header, it is the missing header, and lastLine is valid
		hashBytesInNewLine := newLineBytes[4:36]
		if testHashes(lastHash, hashBytesInNewLine) {
			fmt.Printf("found missing line\n")
			return []string{string(newLineHex)}, true
		}

		// if newLine bridges last-last header and current header, then lastLine is invalid
		lastLastLine := e.Prev().Prev().Value.(string) // not strict...
		if testHashes(doubleSha256(lastLastLine), hashBytesInNewLine) {
			fmt.Printf("replaced a line\n")
			return []string{string(newLineHex)}, false
		}

		// otherwise ...
		ret, isValid := testAndFixHash(lastHash, string(newLineHex), e)
		ret = append(ret, string(newLineHex))
		return ret, isValid
	} else {
		return []string{}, true
	}
}

func testHashes(lastHash, inHeaderHash []byte) bool {
	for i := 0; i < len(lastHash); i++ {
		if lastHash[i] != inHeaderHash[i] {
			// fmt.Printf("value mismatch at %d of %d", i, pos)
			return false
		}
	}
	return true
}

func doubleSha256(header string) []byte {
	v, e := decodeHex(header)
	if e != nil {
		panic(e)
	}
	sha2 := sha256.New()
	sha2.Write(v)
	vv := sha2.Sum([]byte{})
	sha2.Reset()
	sha2.Write(vv)
	return sha2.Sum([]byte{})
}
