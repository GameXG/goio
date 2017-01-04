package mycipher

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"fmt"
	"github.com/gamexg/goio"
	"io"
	"strconv"
)

/*

加密部分


主要是 iv 的问题，是内建iv还是外部iv？

打算外部 iv ，主要的好处是能够更好的实现高性能。

内部 iv 意味着每次新的头都需要建立一个缓冲区。不过可以尝试使用联合流解决。


考虑整个结构，干脆直接完全使用系统函数解决？
直接使用标准的流转发？

不过无所谓了，这里可以之后在处理。


预期加密的地方是：

* socks5 验证部分加密
* socks5 命令部分加密
* 日志上报部分你加密

每次加密需要 2 个参数，一个是 key，一个是 iv 。
提供这两个参数外输入或输出流，提供对应的流。

所以流式实现不用在意


*/

// 加密接口，也是方便依赖注入实现的
type CipherRead struct {
	iv     []byte
	rr     io.Reader    //原始 Read
	block  cipher.Block //块加密器
	stream cipher.Stream
	reader *cipher.StreamReader
}

func NewCipherRead(key string, r io.Reader) (*CipherRead, error) {
	k := sha1.Sum([]byte(key))
	block, err := aes.NewCipher(k[:16])
	if err != nil {
		return nil, err
	}

	c := CipherRead{
		rr:    r,
		block: block,
	}
	return &c, nil
}

// 读取加密信息
// 第一次读 vi 超时会破坏加密环境。
func (c *CipherRead) Read(dst []byte) (n int, err error) {
	if c.reader == nil {
		c.iv = make([]byte, c.block.BlockSize())

		if _, err := io.ReadFull(c.rr, c.iv); err != nil {
			return 0, err
		}

		c.stream = cipher.NewOFB(c.block, c.iv[:])

		c.reader = &cipher.StreamReader{S: c.stream, R: c.rr}
	}

	return c.reader.Read(dst)
}

// 加密接口，也是方便依赖注入实现的
type CipherWrite struct {
	iv     []byte
	rw     io.Writer    //原始 Read
	block  cipher.Block //块加密器
	stream cipher.Stream
	writer *cipher.StreamWriter
}

func NewCipherWrite(key string, w io.Writer) (*CipherWrite, error) {
	k := sha1.Sum([]byte(key))
	block, err := aes.NewCipher(k[:16])
	if err != nil {
		return nil, err
	}

	iv := make([]byte, block.BlockSize())
	if _, err := rand.Read(iv); err != nil {
		return nil, err
	}

	c := CipherWrite{
		iv:    iv,
		rw:    w,
		block: block,
	}
	return &c, nil
}

// 写加密信息
// 写 vi 错误会破坏加密环境。
func (c *CipherWrite) Write(dst []byte) (n int, err error) {
	if c.writer == nil {
		c.stream = cipher.NewOFB(c.block, c.iv[:])
		c.writer = &cipher.StreamWriter{S: c.stream, W: c.rw}

		if _, err := c.rw.Write(c.iv); err != nil {
			return 0, err
		}
	}

	return c.writer.Write(dst)
}

func (c *CipherWrite) Flush() {
	f, _ := c.rw.(goio.Flusher)
	if f != nil {
		f.Flush()
	}
}
func CEn(key []byte, sKey string) error {
	if len(key) != 32 || len(sKey) != 16*2 {
		return fmt.Errorf("key或skey长度错误。")
	}
	key2 := key[16:]

	sKeyb := []byte(sKey)
	for i := 0; i < 16; i++ {
		sk, err := strconv.ParseInt(string(sKeyb[:2]), 16, 16)
		if err != nil {
			return fmt.Errorf("sKey 不是16进制数字。")
		}
		key[0] ^= uint8(sk)
		key2[0] ^= uint8(sk)

		sKeyb = sKeyb[2:]
		key = key[1:]
		key2 = key2[1:]
	}
	return nil
}

func EcbEncrypt(plaintext []byte, key []byte) error {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(plaintext)%aes.BlockSize != 0 {
		return fmt.Errorf("必须为16的倍数")
	}

	text := make([]byte, 16)
	for len(plaintext) > 0 {
		// 每次运算一个block
		cipher.Encrypt(text, plaintext)
		copy(plaintext[:16], text)
		plaintext = plaintext[aes.BlockSize:]
	}
	return nil
}

// 解密
func EcbDecrypt(ciphertext []byte, key []byte) error {
	cipher, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		return fmt.Errorf("必须为16的倍数")
	}

	text := make([]byte, 16)
	for len(ciphertext) > 0 {
		cipher.Decrypt(text, ciphertext)
		copy(ciphertext[:16], text)
		ciphertext = ciphertext[aes.BlockSize:]
	}
	return nil
}

func EncipherMem(b []byte) {
	l := len(b)
	b[0] = b[0] ^ 155
	for i := 1; i < l; i++ {
		b[i] = b[i] ^ b[0]
	}
}

func DecryptMem(b []byte) {
	l := len(b)
	for i := 1; i < l; i++ {
		b[i] = b[i] ^ b[0]
	}
	b[0] = b[0] ^ 155
}
