package misc

import (
	"fmt"
	"io"
	"os"
)

var asciiPrintableMap map[byte]bool

func init() {
	asciiPrintableMap = make(map[byte]bool)
	for i := 0; i < 256; i++ {
		asciiPrintableMap[byte(i)] = false
	}
	s := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ!\"#$%&'()*+,-./:;<=>?@[\\]^_`{|}~ "
	for _, ch := range []byte(s) {
		asciiPrintableMap[ch] = true
	}
}

func isAsciiPrintable(b byte) bool {
	return asciiPrintableMap[b]
}

func DumpHex(w io.Writer, b []byte) error {
	lines := len(b) / 16
	for i := 0; i < lines; i += 1 {
		srep := ""
		for j := 0; j < 16; j += 1 {
			it := b[16*i+j]
			_, err := w.Write([]byte(fmt.Sprintf("%02x ", it)))
			if err != nil {
				return err
			}
			if isAsciiPrintable(it) {
				srep += string([]byte{it})
			} else {
				srep += "."
			}
		}
		_, err := w.Write([]byte("| " + srep + "\n"))
		if err != nil {
			return err
		}
	}
	remain := len(b) - 16*lines
	srep := ""
	j := 0
	for ; j < remain; j += 1 {
		it := b[16*lines+j]
		_, err := w.Write([]byte(fmt.Sprintf("%02x ", it)))
		if err != nil {
			return err
		}
		if isAsciiPrintable(it) {
			srep += string([]byte{it})
		} else {
			srep += "."
		}
	}
	for ; j < 16; j++ {
		_, err := w.Write([]byte("   "))
		if err != nil {
			return err
		}
		srep += " "
	}
	_, err := w.Write([]byte("| " + srep + "\n"))
	if err != nil {
		return err
	}
	return nil
}

func PrintHex(b []byte) {
	err := DumpHex(os.Stdout, b)
	if err != nil {
		fmt.Printf("error in PrintHex() : %s\n", err.Error())
	}
}
