package hlsdl

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

const (
	syncByte = uint8(71) //0x47
)

func decryptAES128(crypted, key, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = pkcs5UnPadding(origData)
	return origData, nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unPadding := int(origData[length-1])
	return origData[:(length - unPadding)]
}

// decrypt descrypts a segment
func (hlsDl *HlsDl) decrypt(segment *Segment) ([]byte, error) {

	file, err := os.Open(segment.Path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if segment.Key != nil {
		key, iv, err := hlsDl.getKey(segment)
		if err != nil {
			return nil, err
		}
		data, err = decryptAES128(data, key, iv)
		if err != nil {
			return nil, err
		}
	}

	for j := 0; j < len(data); j++ {
		if data[j] == syncByte {
			data = data[j:]
			break
		}
	}

	return data, nil
}

func (hlsDl *HlsDl) getKey(segment *Segment) (key []byte, iv []byte, err error) {
	res, err := hlsDl.client.Get(segment.Key.URI)
	if err != nil {
		return nil, nil, err
	}

	if res.StatusCode != 200 {
		return nil, nil, errors.New("failed to get descryption key")
	}

	key, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	iv = []byte(segment.Key.IV)
	if len(iv) == 0 {
		iv = defaultIV(segment.SeqId)
	}
	return
}

func defaultIV(seqID uint64) []byte {
	buf := make([]byte, 16)
	binary.BigEndian.PutUint64(buf[8:], seqID)
	return buf
}

// Decrypt given segment
func (s *Segment) Decrypt(headers map[string]string) ([]byte, error) {
	var data []byte
	copy(data, s.Data)
	if s.Key != nil {
		key, iv, err := s.getKey(headers)
		if err != nil {
			return nil, fmt.Errorf("segment Decrypt: %w", err)
		}
		data, err = decryptAES128(s.Data, key, iv)
		if err != nil {
			return nil, fmt.Errorf("segment Decrypt: %w", err)
		}
	}

	for j := 0; j < len(data); j++ {
		if data[j] == syncByte {
			data = data[j:]
			break
		}
	}
	return data, nil
}

func (s *Segment) getKey(headers map[string]string) (key []byte, iv []byte, err error) {
	client := &http.Client{}
	req, err := newRequest(s.Key.URI, headers)
	if err != nil {
		return nil, nil, fmt.Errorf("segment getKey: %w", err)
	}

	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return nil, nil, errors.New("failed to get descryption key")
	}
	defer resp.Body.Close()

	key, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("segment getKey: %w", err)
	}

	iv = []byte(s.Key.IV)
	if len(iv) == 0 {
		iv = defaultIV(s.SeqId)
	}
	return key, iv, nil
}
