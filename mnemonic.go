package mnemonic

import (
	"errors"
	"fmt"
	"github.com/tyler-smith/go-bip39/wordlists"
	"golang.org/x/exp/slices"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type Packer interface {
	Pack() (Converter, error)
	GetPhrase() string
	GetMnemonicOriginal() []string
	GetLang() string
	UnPack(base int) (Converter, error)
	disableCheck() Packer
}

type mnemonic struct {
	isCheck          bool
	wordlist         []string
	mnemonicOriginal []string
	lang             string
	phrase           string
	convert          convert
	supportBases     []int
}

func New(phrase string, lang string) Packer {
	return &mnemonic{
		phrase:       phrase,
		lang:         lang,
		supportBases: []int{Base10, Base36, Base62},
		isCheck:      true,
		convert: convert{
			base: make(map[int]string),
		},
	}
}

func (m *mnemonic) GetPhrase() string {
	return m.phrase
}
func (m *mnemonic) GetMnemonicOriginal() []string {
	return m.mnemonicOriginal
}
func (m *mnemonic) GetLang() string {
	return m.lang
}
func (m *mnemonic) parseOriginal() error {
	var err error
	lengthWordLists := len(m.wordlist)
	lengthMnemonic := len(m.mnemonicOriginal)

	m.convert.numberMnemonic = make([]string, 0, lengthMnemonic)
	m.convert.mnemonic = make([]string, 0, lengthMnemonic)
	m.convert.mnemonicShort = make([]string, 0, lengthMnemonic)

	for _, word := range m.mnemonicOriginal {
		word = strings.ToLower(word)

		matched, _ := regexp.MatchString("^\\d+$", word)

		var index int
		if i, errInt := strconv.Atoi(word); matched && errInt == nil && lengthWordLists >= i {
			index = i
		} else {
			index = slices.IndexFunc(m.wordlist, func(val string) bool {
				return m.getShort(val) == m.getShort(word)
			})
		}

		if index == -1 {
			err = errors.New(fmt.Sprintf("not found word - %s", word))
			break
		}

		m.convert.numberMnemonic = append(m.convert.numberMnemonic, m.getIndexFormat(index))
		m.convert.mnemonic = append(m.convert.mnemonic, m.wordlist[index])
		m.convert.mnemonicShort = append(m.convert.mnemonicShort, m.getShort(m.wordlist[index]))
	}
	return err
}

func (m *mnemonic) getIndexFormat(index int) string {
	format := "%0" + fmt.Sprintf("%d", numberOfLettersForIdentification) + "d"
	return fmt.Sprintf(format, index)
}
func (m *mnemonic) getShort(word string) string {
	asRunes := []rune(word)
	return string(asRunes[0:numberOfLettersForIdentification])
}
func (m *mnemonic) UnPack(base int) (Converter, error) {
	err := m.getWordList()
	if err != nil {
		return nil, err
	}

	if slices.Index(m.supportBases, base) == -1 {
		return nil, errors.New(fmt.Sprintf("not support base - %d", base))
	}

	distance := new(big.Int)
	distance.SetString(m.phrase, base)

	for _, b := range m.supportBases {
		m.convert.base[b] = distance.Text(b)
	}

	asRunes := []rune(m.convert.base[Base10])

	countSymbols := numberOfLettersForIdentification - 1

	var reverse []int
	for i := len(asRunes) - 1; i >= 0; i-- {
		start := i - countSymbols
		if start < 0 {
			start = 0
		}
		wordIndex := string(asRunes[start : i+1])
		index, _ := strconv.Atoi(wordIndex)
		reverse = append(reverse, index)
		i = start
	}
	m.convert.numberMnemonic = make([]string, 0, len(reverse))
	m.convert.mnemonic = make([]string, 0, len(reverse))
	for i := len(reverse) - 1; i >= 0; i-- {
		m.convert.numberMnemonic = append(m.convert.numberMnemonic, m.getIndexFormat(reverse[i]))
		m.convert.mnemonic = append(m.convert.mnemonic, m.wordlist[reverse[i]])
	}

	return &m.convert, m.check()
}
func (m *mnemonic) setMnemonicOriginal() *mnemonic {
	r := regexp.MustCompile("\\S+")
	m.mnemonicOriginal = r.FindAllString(m.phrase, -1)

	return m
}
func (m *mnemonic) Pack() (Converter, error) {
	err := m.getWordList()
	if err != nil {
		return nil, err
	}

	err = m.setMnemonicOriginal().parseOriginal()
	if err != nil {
		return nil, err
	}

	var re = regexp.MustCompile("^(0+)")
	base10 := strings.Join(m.convert.numberMnemonic, "")
	base10 = re.ReplaceAllString(base10, "")
	distance := new(big.Int)
	distance.SetString(base10, Base10)

	for _, base := range m.supportBases {
		m.convert.base[base] = distance.Text(base)
	}

	return &m.convert, m.check()
}

func (m *mnemonic) getWordList() error {
	var w []string
	var err error
	switch m.lang {
	case English:
		w = wordlists.English
	case Czech:
		w = wordlists.Czech
	case French:
		w = wordlists.French
	case Italian:
		w = wordlists.Italian
	case Japanese:
		w = wordlists.Japanese
	case Korean:
		w = wordlists.Korean
	case Spanish:
		w = wordlists.Spanish
	case ChineseTraditional:
		w = wordlists.ChineseTraditional
	case ChineseSimplified:
		w = wordlists.ChineseSimplified
	default:
		err = errors.New("not support lang")
	}
	m.wordlist = w

	return err
}
func (m *mnemonic) disableCheck() Packer {
	m.isCheck = false
	return m
}
func (m *mnemonic) check() error {
	if !m.isCheck {
		return nil
	}

	convertToString := m.convert.toString()

	type toCheck struct {
		phrase string
		method string
		args   []reflect.Value
	}

	phrases := []toCheck{
		{strings.Join(m.convert.mnemonic, " "), "Pack", nil},
		{strings.Join(m.convert.numberMnemonic, " "), "Pack", nil},
	}

	for _, b := range m.supportBases {
		phrases = append(phrases, toCheck{m.convert.Get(b), "UnPack", []reflect.Value{reflect.ValueOf(b)}})
	}

	var err error
	for _, c := range phrases {
		object := New(c.phrase, m.lang).disableCheck()
		method := reflect.ValueOf(object).MethodByName(c.method)

		result := method.Call(c.args)

		if e, ok := result[1].Interface().(error); ok && e != nil || convertToString != result[0].Interface().(Converter).toString() {
			err = errors.New("failed check")
			break
		}
	}

	return err
}
