package custom

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/prysmaticlabs/prysm/v5/container/slice"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Parse struct {
	HeadersByHash     map[string]string // header - hash
	HeadersByHeight   map[int64]string  // height -> header
	BlockHashByHeight map[int64]string  // height -> hash
	lock              sync.RWMutex      // todo
}

func (p *Parse) GetHeaderByHash(hash string) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	header, ok := p.HeadersByHash[hash]
	if !ok {
		return "", fmt.Errorf("header not found: %v", hash)
	}
	return header, nil
}

func (p *Parse) GetHeaderByHeight(height int64) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	header, ok := p.HeadersByHeight[height]
	if !ok {
		return "", fmt.Errorf("header not found: %v", height)
	}
	return header, nil
}

func (p *Parse) GetBlockHashByHeight(height int64) (string, error) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	hash, ok := p.BlockHashByHeight[height]
	if !ok {
		return "", fmt.Errorf("header not found: %v", height)
	}
	return hash, nil
}

func NewParse(path string) (*Parse, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	HeadersByHeight := make(map[int64]string)
	HeadersByHash := make(map[string]string)
	BlockHashByHeight := make(map[int64]string)
	var preHeader string
	scanner := bufio.NewScanner(file)
	index := int64(0)
	for scanner.Scan() {
		header := scanner.Text()
		HeadersByHeight[index] = header
		headerBytes, err := hex.DecodeString(header)
		if err != nil {
			return nil, err
		}
		//982051fd
		//00e00020094b5519bc3c0386f11b527c9da4419797e5e235a3bb000000000000000000008c61ec99cc331268d135f44bc8ad7d062feb82757afb773ed79545b73fe9fcb59cd0f7633930071720e5151d
		if index > 0 {
			preHashBytes := headerBytes[4:36]
			preHash := fmt.Sprintf("%x", slice.Reverse(preHashBytes))
			HeadersByHash[preHash] = preHeader
			BlockHashByHeight[index-1] = preHash
		}
		index++
		preHeader = header

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return &Parse{
		HeadersByHash:     HeadersByHash,
		HeadersByHeight:   HeadersByHeight,
		BlockHashByHeight: BlockHashByHeight,
	}, nil
}

func traverseFile(root string) (map[string]string, error) {
	files := make(map[string]string)
	err := filepath.WalkDir(root, func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && root != path {
			fileName, err := getFileName(info.Name())
			if err != nil {
				return err
			}
			files[fileName] = path
		}
		return nil
	})
	return files, err
}
func getFileName(path string) (string, error) {
	arrs := strings.Split(path, "/")
	if len(arrs) == 0 {
		return "", fmt.Errorf("get file name error")
	}
	return arrs[len(arrs)-1], nil
}
