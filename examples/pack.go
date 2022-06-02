package main

import (
	"fmt"
	mnemonic "github.com/evgeny-klyopov/go-pack-mnemonic"
)

func main() {
	phrase := "able 1644 awful      crisp collect claim educate can ball toy wealth spell window young reveal zone review kit decorate delay defy air pear brave"

	fmt.Println("Pack:")
	pack := mnemonic.New(phrase, mnemonic.English)
	convert, err := pack.Pack()
	printResult(pack, convert, err)

	fmt.Println("UnPack:")

	for _, base := range []int{mnemonic.Base10, mnemonic.Base36, mnemonic.Base62} {
		unPack := mnemonic.New(convert.Get(base), mnemonic.English)
		convert, err = pack.Pack()
		printResult(unPack, convert, err)
	}
}

func printResult(pack mnemonic.Packer, convert mnemonic.Converter, err error) {
	fmt.Println("Error:", err)
	fmt.Println("Phrase:", pack.GetPhrase())
	fmt.Println("Lang:", pack.GetLang())
	fmt.Println("MnemonicOriginal:", pack.GetMnemonicOriginal())
	fmt.Println("Base10:", convert.Get(mnemonic.Base10))
	fmt.Println("Base36:", convert.Get(mnemonic.Base36))
	fmt.Println("Base62:", convert.Get(mnemonic.Base62))
	fmt.Println("Mnemonic:", convert.GetMnemonic())
	fmt.Println("MnemonicShort (max 4 letters):", convert.GetMnemonicShort())
	fmt.Println("NumberMnemonic:", convert.GetNumberMnemonic())
	fmt.Println("")
}
