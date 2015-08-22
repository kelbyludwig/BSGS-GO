package main

import (
	"fmt"
	"github.com/cznic/mathutil"
	"math/big"
)

func main() {

	//The known modulus. Must be an odd prime.
	modulus := big.NewInt(29)

	//The value we want to solve the DLP for. Should not be divisble by the modulus.
	h := big.NewInt(15)

	//The generator.
	g := big.NewInt(2)

	//The ceiling of the square root of the modulus.
	m := mathutil.SqrtBig(modulus)
	m.Add(m, big.NewInt(1))

	//Little m is just used for the baby step loop.
	little_m := m.Int64()

	//TODO: For now this will be the string rep of BigInts for the map key.
	//      This is terribly inefficient and I hope to convert this to a bloom
	//      filter that leverages bit-level ops.
	baby_step := make(map[string]int64)
	var bs big.Int
	for i := int64(0); i < little_m; i++ {
		bs.Exp(g, big.NewInt(i), modulus)
		baby_step[bs.String()] = i
	}

	var inv big.Int
	var gm big.Int
	gm.Exp(g, m, modulus)
	inv.ModInverse(&gm, modulus)

	gs := big.NewInt(0)
	var res big.Int
	for {
		res.Exp(&inv, gs, modulus)
		res.Mul(&res, h)
		res.Mod(&res, modulus)
		if val,exists := baby_step[res.String()]; exists {
            var dl big.Int
            gs.Mul(gs,m)
            dl.Add(big.NewInt(val), gs)
			fmt.Println("[RESULT] ", dl.String())
			return
		}
		gs.Add(gs, big.NewInt(1))
	}
}
