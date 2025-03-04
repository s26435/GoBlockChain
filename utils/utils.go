package utils

import (
	"strconv"
	"encoding/hex"
)

func Must(err error){
	if err != nil {
		panic(err)
	}
}

func IntToHex(n int64) []byte {
	hexStr := strconv.FormatInt(n, 16)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	bytes, _ := hex.DecodeString(hexStr)
	return bytes
}