package hashids

import (
	"math"
	"reflect"
	"testing"
	"math/big"
)

func TestBigEncodeDecode(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 30
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := makeNumbers()
	hash, err := hid.EncodeBigInt(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := hid.DecodeBigInt(hash)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestBigEncodeWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := makeNumbers()

	hash, err := hid.EncodeBigInt(numbers)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v", numbers, hash)

	if hash != "1ZZhw9wwZE8aj5LByEFWNipkHK2loxbgn38bK" {
		t.Error("hash does not match expected one")
	}
}

func TestBigDecodeWithKnownHash(t *testing.T) {
	hdata := NewData()
	hdata.MinLength = 0
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	hash := "1ZZhw9wwZE8aj5LByEFWNipkHK2loxbgn38bK"
	numbers, err := hid.DecodeBigInt(hash)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v", hash, numbers)

	expected := makeNumbers()
	if !reflect.DeepEqual(numbers, expected) {
		t.Errorf("Decoded numbers `%v` did not match with expected `%v`", numbers, expected)
	}
}

func TestBigMinLength(t *testing.T) {
	hdata := NewData()
	hdata.Salt = "salt1"
	hdata.MinLength = 10
	hid, _ := NewWithData(hdata)
	hid.EncodeBigInt([]*big.Int{big.NewInt(0)})
}

func TestBigCustomAlphabet(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)

	numbers := makeNumbers()
	hash, err := hid.EncodeBigInt(numbers)
	if err != nil {
		t.Fatal(err)
	}
	dec, err := hid.DecodeBigInt(hash)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	if !reflect.DeepEqual(dec, numbers) {
		t.Errorf("Decoded numbers `%v` did not match with original `%v`", dec, numbers)
	}
}

func TestBigDecodeWithError(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "this is my salt"

	hid, _ := NewWithData(hdata)
	// hash now contains a letter not in the alphabet
	dec, err := hid.DecodeBigInt("MAkhkloFAxAoskaZ")

	if dec != nil {
		t.Errorf("Expected `nil` but got `%v`", dec)
	}
	expected := "alphabet used for hash was different"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

// tests issue #28
func TestBigDecodeWithWrongSalt(t *testing.T) {
	hdata := NewData()
	hdata.Alphabet = "PleasAkMEFoThStx"
	hdata.Salt = "temp"

	hidEncode, _ := NewWithData(hdata)

	numbers := makeNumbers()
	hash, _ := hidEncode.EncodeBigInt(numbers)

	hdata.Salt = "test"
	hidDecode, _ := NewWithData(hdata)
	dec, err := hidDecode.DecodeWithError(hash)

	t.Logf("%v -> %v -> %v", numbers, hash, dec)

	expected := "mismatch between encode and decode: lAPTelMEEeaMaPEPPEPkaEexMPMskMklteMxTMkPaPAxElkAlakxllxAk start olkhA re-encoded. result: [73 0]"
	if err == nil || err.Error() != expected {
		t.Errorf("Expected error `%s` but got `%s`", expected, err)
	}
}

func makeNumbers() []*big.Int {
	numbers := []*big.Int{
		big.NewInt(45),
		big.NewInt(434),
		big.NewInt(1313),
		big.NewInt(99),
		big.NewInt(math.MaxInt64),
	}
	numbers[1].Mul(numbers[1], numbers[4])
	return numbers
}
