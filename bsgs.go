package main

import (
	"fmt"
	"math/big"
	"github.com/cznic/mathutil"
    "github.com/kelbyludwig/fnv"
)

func add_bloom(bloom uint64, value []byte) uint64 {
    bloom |= fnv.FNV1A(value)
    return bloom
}

func check_bloom(bloom uint64, value []byte) bool {
    test := (fnv.FNV1A(value)) | bloom
    if test == bloom {
        return true
    }
    return false
}

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

    //Tiny bloom filter!
    var bloom uint64
    bloom = 0

	baby_step := make(map[string]int64)
	var bs big.Int

    //Speed up modexp by saving previous result.
    var prev *big.Int
	for i := int64(0); i < little_m; i++ {
        if prev == nil {
           bs.Exp(g, big.NewInt(i), modulus)
           bloom = add_bloom(bloom, bs.Bytes())
           baby_step[bs.String()] = i
           prev = &bs
           continue
        }
		bs.Mul(prev, g)
        bs.Mod(&bs, modulus)
        bloom = add_bloom(bloom, bs.Bytes())
		baby_step[bs.String()] = i
        prev = &bs
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
        if check_bloom(bloom, res.Bytes()) {
            fmt.Println("[DEBUG] Bloom hit!")
		    if val,exists := baby_step[res.String()]; exists {
                var dl big.Int
                gs.Mul(gs,m)
                dl.Add(big.NewInt(val), gs)
                fmt.Println("[RESULT] ", dl.String())
                return
		    }
        } else {
            fmt.Println("[DEBUG] Bloom miss!")
        }
        if gs.Cmp(m) == 1 {
            fmt.Println("[RESULT] Exponent not found.")
            return
        }
		gs.Add(gs, big.NewInt(1))
	}
}
