package main

import "math/big"

func (app *application) getFactorial(x *big.Int) *big.Int {
	n := big.NewInt(1)
	if x.Cmp(big.NewInt(0)) == 0 {
		return n
	}
	return n.Mul(x, app.getFactorial(n.Sub(x, n)))
}
