package mycipher

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"testing"
)

type A int

func (a A) Write(buf []byte) (int, error) {
	fmt.Println(len(buf))
	return len(buf), nil
}

/*
输出结果：
16
1048576

// 实际看源码，可以看到内部是创建了一个和输入相同的缓冲区，
// 然后将加密数据保存到这个缓冲区，之后直接写的这个缓冲区。
// 所以尺寸注定是和输入一致。
// 不过还有另外一个独立输出，即首次写时的随机 VI
*/
func TestSize(t *testing.T) {
	d := make([]byte, 1024*1024)
	rand.Read(d)

	a := A(0)
	zw, err := NewCipherWrite("key", a)
	if err != nil {
		t.Fatal(err)
	}
	zw.Write(d)
	zw.Flush()
}

func TestA(t *testing.T) {
	key := "sdygvetse"
	buf := bytes.Buffer{}

	d0 := make([]byte, 1024)
	d1 := []byte{1, 2, 3}
	d2 := []byte{4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23}
	d3 := []byte{24, 25, 26, 27}

	cw, err := NewCipherWrite(key, &buf)
	if err != nil {
		t.Fatal(err)
	}

	cr, err := NewCipherRead(key, &buf)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := cw.Write(d1); err != nil {
		t.Fatal(err)
	}

	if n, err := cr.Read(d0); err != nil {
		t.Fatal(err)
	} else if bytes.Equal(d0[:n], d1) == false {
		t.Fatal(d0[:n], "!=", d1)
	}

	if _, err := cw.Write(d2); err != nil {
		t.Fatal(err)
	}

	if _, err := cw.Write(d3); err != nil {
		t.Fatal(err)
	}

	if n, err := cr.Read(d0); err != nil {
		t.Fatal(err)
	} else {
		d23 := append(d2, d3...)
		if bytes.Equal(d0[:n], d23) == false {
			t.Fatal(d0[:n], "!=", d23)
		}
	}
}

func TestCEn(t *testing.T) {
	key := []byte{0xC2, 0x3E, 0x26, 0xA3, 0x47, 0x09, 0x09, 0x66, 0x0E, 0xF8, 0x0C, 0xB7, 0x01, 0x50, 0xFE, 0x1A, 0x0A, 0xD7, 0xD9, 0x48, 0x33, 0xE7, 0x02, 0xDF, 0x3E, 0x0A, 0xF9, 0x6B, 0x3B, 0x71, 0xC5, 0x27}
	sKey := "4CB7F4EA57154A5798D8AA8760227010"
	eKey := []byte{142, 137, 210, 73, 16, 28, 67, 49, 150, 32, 166, 48, 97, 114, 142, 10, 70, 96, 45, 162, 100, 242, 72, 136, 166, 210, 83, 236, 91, 83, 181, 55}

	if err := CEn(key, sKey); err != nil {
		t.Fatal(err)
	} else if bytes.Equal(key, eKey) != true {
		t.Fatal(key, "!=", eKey)
	}
}

func TestEcb(t *testing.T) {
	r := make([]byte, 32*4)
	k := make([]byte, 32)
	e := make([]byte, 32*4)
	if _, err := rand.Read(r); err != nil {
		t.Fatal(err)
	}
	if _, err := rand.Read(k); err != nil {
		t.Fatal(err)
	}
	copy(r, e)

	if err := EcbEncrypt(e, k); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(r, e) == true {
		t.Fatal(r, "==", e)
	}

	if err := EcbDecrypt(e, k); err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(r, e) != true {
		t.Fatal(r, "!=", e)
	}

}

func TestC(t *testing.T) {
	b1 := []byte("2016-08-01 01:00:00")
	b2 := make([]byte, len(b1))
	copy(b2, b1)

	EncipherMem(b1)

	fmt.Printf("%#v", b1)

	if bytes.Equal(b1, b2) == true {
		t.Error(b1, "==", b2)
	}

	DecryptMem(b1)

	if bytes.Equal(b1, b2) != true {
		t.Error(b1, "!=", b2)
	}
}
