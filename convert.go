package mnemonic

import (
	"encoding/json"
	"sort"
)

type convert struct {
	mnemonic       []string
	mnemonicShort  []string
	numberMnemonic []string
	base           map[int]string
}

type Converter interface {
	toString() string
	GetMnemonic() []string
	GetMnemonicShort() []string
	GetNumberMnemonic() []string
	Get(base int) string
}

func (c *convert) toString() string {
	str := struct {
		M []string  `json:"m"`
		N []string  `json:"n"`
		B [3]string `json:"b"`
	}{
		M: c.mnemonic,
		N: c.numberMnemonic,
	}
	keys := make([]int, 0, len(c.base))
	for k := range c.base {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	for i, k := range keys {
		str.B[i] = c.base[k]
	}
	j, _ := json.Marshal(str)
	return string(j)
}

func (c *convert) GetMnemonic() []string {
	return c.mnemonic
}
func (c *convert) GetMnemonicShort() []string {
	return c.mnemonicShort
}
func (c *convert) GetNumberMnemonic() []string {
	return c.numberMnemonic
}
func (c *convert) Get(base int) string {
	return c.base[base]
}
