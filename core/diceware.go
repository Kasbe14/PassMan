package core

import (
	"crypto/rand"
	_ "embed"
	"fmt"
	"math/big"
	"strings"
)

//go:embed eff_short_wordlist_2_0.txt
var words string

var dicewords []string

//init function -> runs before the main function
func init() {
    dicewords = strings.Split(strings.TrimSpace(words), "\n")
}

func GenerateDiceWords() (string,error) {
    uBound := big.NewInt(int64(len(dicewords)))
    dice := make([]string,4)
    for i := 0; i<4; i++ {
        n, err := rand.Int(rand.Reader, uBound); if err != nil {
            return "",err
        }
        dice[i] = dicewords[n.Int64()]
    }
    num , err := rand.Int(rand.Reader,big.NewInt(9000))
    if err != nil {
        return "",err
    }
    num_ := fmt.Sprintf("%d",num.Int64() + 1000)
    dice = append(dice, num_)
    return strings.Join(dice,"-"), nil
}
