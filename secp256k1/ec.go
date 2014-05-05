package secp256k1

import (
	//"encoding/hex"
)


func ecdsa_verify(pubkey, sig, msg []byte) int {
	var m Number
	var s Signature
	m.SetBytes(msg)

	var q XY
	if !q.ParsePubkey(pubkey) {
		return -1
	}

	if !s.sig_parse(sig) {
		return -2
	}

	if !s.Verify(&q, &m) {
		return 0
	}
	return 1
}

func Verify(k, s, m []byte) bool {
	return ecdsa_verify(k, s, m)==1
}

func init() {
	init_contants()
	ecmult_start()
}

func DecompressPoint(X []byte, off bool, Y []byte) {
	var rx, ry, c, x2, x3 Field
	rx.SetB32(X)
	rx.Sqr(&x2)
	rx.Mul(&x3, &x2)
	c.SetInt(7)
	c.SetAdd(&x3)
	c.Sqrt(&ry)
	ry.Normalize()
	if ry.IsOdd() != off {
		ry.Negate(&ry, 1)
	}
	ry.Normalize()
	ry.GetB32(Y)
	return
}


func RecoverPublicKey(r, s, h []byte, recid int, X, Y []byte) bool {
	var sig Signature
	var pubkey XY
	var msg Number
	sig.R.SetBytes(r)
	if sig.R.Sign()<=0 || sig.R.Cmp(&TheCurve.Order.Int)>=0 {
		return false
	}
	sig.S.SetBytes(s)
	if sig.S.Sign()<=0 || sig.S.Cmp(&TheCurve.Order.Int)>=0 {
		return false
	}
	msg.SetBytes(h)
	if !sig.recover(&pubkey, &msg, recid) {
		return false
	}
	pubkey.X.GetB32(X)
	pubkey.Y.GetB32(Y)
	return true
}


// Standard EC multiplacation k(xy)
// xy - is the standarized public key format (33 or 65 bytes long)
// out - should be the buffer for 33 bytes (1st byte will be set to either 02 or 03)
func Multiply(xy, k, out []byte) bool {
	var pk XY
	if !pk.ParsePubkey(xy) {
		return false
	}
	if !pk.Multi(k) {
		return false
	}
	pk.GetPublicKey(out)
	return true
}


func BaseMultiply(k, out []byte) bool {
	var r XYZ
	var n Number
	var pk XY
	n.SetBytes(k)
	ECmultGen(&r, &n)
	pk.SetXYZ(&r)
	pk.GetPublicKey(out)
	return true
}