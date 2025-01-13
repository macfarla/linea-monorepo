// Copyright 2020 Consensys Software Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by consensys/gnark-crypto DO NOT EDIT

package fext

import (
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fp"
	"github.com/consensys/gnark-crypto/ecc/bls12-377/fr"
	"github.com/consensys/linea-monorepo/prover/maths/field"
	"math/big"
	"math/rand"
)

const noQNR = 11

// Element is a degree two finite field extension of fr.Element
type Element struct {
	A0, A1 fr.Element
}

// Equal returns true if z equals x, false otherwise
func (z *Element) Equal(x *Element) bool {
	return z.A0.Equal(&x.A0) && z.A1.Equal(&x.A1)
}

// Bits
// TODO @gbotrel fixme this shouldn't return a Element
func (z *Element) Bits() Element {
	r := Element{}
	r.A0 = z.A0.Bits()
	r.A1 = z.A1.Bits()
	return r
}

// Cmp compares (lexicographic order) z and x and returns:
//
//	-1 if z <  x
//	 0 if z == x
//	+1 if z >  x
func (z *Element) Cmp(x *Element) int {
	if a1 := z.A1.Cmp(&x.A1); a1 != 0 {
		return a1
	}
	return z.A0.Cmp(&x.A0)
}

// LexicographicallyLargest returns true if this element is strictly lexicographically
// larger than its negation, false otherwise
func (z *Element) LexicographicallyLargest() bool {
	// adapted from github.com/zkcrypto/bls12_381
	if z.A1.IsZero() {
		return z.A0.LexicographicallyLargest()
	}
	return z.A1.LexicographicallyLargest()
}

// SetString sets a Element element from strings
func (z *Element) SetString(s1, s2 string) (*Element, error) {
	_, err := z.A0.SetString(s1)
	if err != nil {
		return z, err
	}

	_, err = z.A1.SetString(s2)
	if err != nil {
		return z, err
	}
	return z, nil
}

// SetZero sets an Element elmt to zero
func (z *Element) SetZero() *Element {
	z.A0.SetZero()
	z.A1.SetZero()
	return z
}

// Set sets an Element from x
func (z *Element) Set(x *Element) *Element {
	z.A0 = x.A0
	z.A1 = x.A1
	return z
}

// SetOne sets z to 1 in Montgomery form and returns z
func (z *Element) SetOne() *Element {
	z.A0.SetOne()
	z.A1.SetZero()
	return z
}

// SetRandom sets a0 and a1 to random values
func (z *Element) SetRandom() (*Element, error) {
	if _, err := z.A0.SetRandom(); err != nil {
		return nil, err
	}
	if _, err := z.A1.SetRandom(); err != nil {
		return nil, err
	}
	return z, nil
}

// IsZero returns true if z is zero, false otherwise
func (z *Element) IsZero() bool {
	return z.A0.IsZero() && z.A1.IsZero()
}

// IsOne returns true if z is one, false otherwise
func (z *Element) IsOne() bool {
	return z.A0.IsOne() && z.A1.IsZero()
}

// Add adds two elements of Element
func (z *Element) Add(x, y *Element) *Element {
	addE2New(z, x, y)
	return z
}

// Sub subtracts two elements of Element
func (z *Element) Sub(x, y *Element) *Element {
	subE2New(z, x, y)
	return z
}

// Double doubles an Element element
func (z *Element) Double(x *Element) *Element {
	doubleE2New(z, x)
	return z
}

// Neg negates an Element element
func (z *Element) Neg(x *Element) *Element {
	negE2New(z, x)
	return z
}

// String implements Stringer interface for fancy printing
func (z *Element) String() string {
	return z.A0.String() + "+" + z.A1.String() + "*u"
}

// MulByElement multiplies an element in Element by an element in fp
func (z *Element) MulByElement(x *Element, y *fr.Element) *Element {
	var yCopy fr.Element
	yCopy.Set(y)
	z.A0.Mul(&x.A0, &yCopy)
	z.A1.Mul(&x.A1, &yCopy)
	return z
}

// Conjugate conjugates an element in Element
func (z *Element) Conjugate(x *Element) *Element {
	z.A0 = x.A0
	z.A1.Neg(&x.A1)
	return z
}

// Halve sets z to z / 2
func (z *Element) Halve() {
	z.A0.Halve()
	z.A1.Halve()
}

// Legendre returns the Legendre symbol of z
func (z *Element) Legendre() int {
	var n fr.Element
	z.norm(&n)
	return n.Legendre()
}

// Exp sets z=xᵏ (mod q²) and returns it
func (z *Element) Exp(x Element, k *big.Int) *Element {
	if k.IsUint64() && k.Uint64() == 0 {
		return z.SetOne()
	}

	e := k
	if k.Sign() == -1 {
		// negative k, we invert
		// if k < 0: xᵏ (mod q²) == (x⁻¹)ᵏ (mod q²)
		x.Inverse(&x)

		// we negate k in a temp big.Int since
		// Int.Bit(_) of k and -k is different
		e = bigIntPool.Get().(*big.Int)
		defer bigIntPool.Put(e)
		e.Neg(k)
	}

	z.SetOne()
	b := e.Bytes()
	for i := 0; i < len(b); i++ {
		w := b[i]
		for j := 0; j < 8; j++ {
			z.Square(z)
			if (w & (0b10000000 >> j)) != 0 {
				z.Mul(z, &x)
			}
		}
	}

	return z
}

// Sqrt sets z to the square root of and returns z
// The function does not test whether the square root
// exists or not, it's up to the caller to call
// Legendre beforehand.
// cf https://eprint.iacr.org/2012/685.pdf (algo 10)
func (z *Element) Sqrt(x *Element) *Element {

	// precomputation
	var b, c, d, e, f, x0 Element
	var _b, o fr.Element

	// c must be a non square (works for p=1 mod 12 hence 1 mod 4, only bls377 has such a p currently)
	c.A1.SetOne()

	q := fp.Modulus()
	var exp, one big.Int
	one.SetUint64(1)
	exp.Set(q).Sub(&exp, &one).Rsh(&exp, 1)
	d.Exp(c, &exp)
	e.Mul(&d, &c).Inverse(&e)
	f.Mul(&d, &c).Square(&f)

	// computation
	exp.Rsh(&exp, 1)
	b.Exp(*x, &exp)
	b.norm(&_b)
	o.SetOne()
	if _b.Equal(&o) {
		x0.Square(&b).Mul(&x0, x)
		_b.Set(&x0.A0).Sqrt(&_b)
		z.Conjugate(&b).MulByElement(z, &_b)
		return z
	}
	x0.Square(&b).Mul(&x0, x).Mul(&x0, &f)
	_b.Set(&x0.A0).Sqrt(&_b)
	z.Conjugate(&b).MulByElement(z, &_b).Mul(z, &e)

	return z
}

// BatchInvertE2New returns a new slice with every element in a inverted.
// It uses Montgomery batch inversion trick.
//
// if a[i] == 0, returns result[i] = a[i]
func BatchInvertE2New(a []Element) []Element {
	res := make([]Element, len(a))
	if len(a) == 0 {
		return res
	}

	zeroes := make([]bool, len(a))
	var accumulator Element
	accumulator.SetOne()

	for i := 0; i < len(a); i++ {
		if a[i].IsZero() {
			zeroes[i] = true
			continue
		}
		res[i].Set(&accumulator)
		accumulator.Mul(&accumulator, &a[i])
	}

	accumulator.Inverse(&accumulator)

	for i := len(a) - 1; i >= 0; i-- {
		if zeroes[i] {
			continue
		}
		res[i].Mul(&res[i], &accumulator)
		accumulator.Mul(&accumulator, &a[i])
	}

	return res
}

// Select is conditional move.
// If cond = 0, it sets z to caseZ and returns it. otherwise caseNz.
func (z *Element) Select(cond int, caseZ *Element, caseNz *Element) *Element {
	//Might be able to save a nanosecond or two by an aggregate implementation

	z.A0.Select(cond, &caseZ.A0, &caseNz.A0)
	z.A1.Select(cond, &caseZ.A1, &caseNz.A1)

	return z
}

// Div divides an element in Element by an element in Element
func (z *Element) Div(x *Element, y *Element) *Element {
	var r Element
	r.Inverse(y).Mul(x, &r)
	return z.Set(&r)
}

func PseudoRand(rng *rand.Rand) Element {
	x := field.PseudoRand(rng)
	y := field.PseudoRand(rng)
	result := new(Element).SetZero()
	return *result.Add(result, &Element{x, y})
}
