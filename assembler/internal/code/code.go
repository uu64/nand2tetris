package code

var compMap = map[string][]byte{
	// a = 0
	"0":   []byte("0101010"),
	"1":   []byte("0111111"),
	"-1":  []byte("0111010"),
	"D":   []byte("0001100"),
	"A":   []byte("0110000"),
	"!D":  []byte("0001101"),
	"!A":  []byte("0110001"),
	"-D":  []byte("0001111"),
	"-A":  []byte("0110011"),
	"D+1": []byte("0011111"),
	"A+1": []byte("0110111"),
	"D-1": []byte("0001110"),
	"A-1": []byte("0110010"),
	"D+A": []byte("0000010"),
	"D-A": []byte("0010011"),
	"A-D": []byte("0000111"),
	"D&A": []byte("0000000"),
	"D|A": []byte("0010101"),
	// a = 1
	"M":   []byte("1110000"),
	"!M":  []byte("1110001"),
	"-M":  []byte("1110011"),
	"M+1": []byte("1110111"),
	"M-1": []byte("1110010"),
	"D+M": []byte("1000010"),
	"D-M": []byte("1010011"),
	"M-D": []byte("1000111"),
	"D&M": []byte("1000000"),
	"D|M": []byte("1010101"),
}

var destMap = map[string][]byte{
	"null": []byte("000"),
	"M":    []byte("001"),
	"D":    []byte("010"),
	"MD":   []byte("011"),
	"A":    []byte("100"),
	"AM":   []byte("101"),
	"AD":   []byte("110"),
	"AMD":  []byte("111"),
}

var jumpMap = map[string][]byte{
	"null": []byte("000"),
	"JGT":  []byte("001"),
	"JEQ":  []byte("010"),
	"JGE":  []byte("011"),
	"JLT":  []byte("100"),
	"JNE":  []byte("101"),
	"JLE":  []byte("110"),
	"JMP":  []byte("111"),
}

func Dest(mnemonic string) []byte {
	return destMap[mnemonic]
}

func Comp(mnemonic string) []byte {
	return compMap[mnemonic]
}

func Jump(mnemonic string) []byte {
	return jumpMap[mnemonic]
}
