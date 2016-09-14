package main

import (
	"fmt"

	"time"

	"bitbucket.org/jack/jackvpn/common/zip"
)

type write struct {
	l int
}

func (w *write) Write(buf []byte) (int, error) {
	w.l += len(buf)
	return len(buf), nil
}

func testZip(zipType string) (string, error) {
	w := write{}
	oLen := 0

	z, err := zip.NewZipWrite(&w, zipType)
	if err != nil {
		return "", err
	}

	for _, data := range HttpData {
		oLen += len(data)
		if _, err := z.Write(data); err != nil {
			return "", err
		}
	}

	return fmt.Sprintf("压缩类型：%v\r\n压缩前大小：%v\r\n压缩后大小：%v\r\n压缩率：%v", zipType, oLen, w.l,
		float32(w.l)/float32(oLen)), nil
}

func main() {
	for _, zipType := range []string{"zlib", "gzip", "deflate", "zlib:1", "gzip:1", "deflate:1", "zlib:9", "gzip:9", "deflate:9"} {
		sTime := time.Now()
		s, err := testZip(zipType)
		if err != nil {
			panic(err)
		}
		e := time.Now().Sub(sTime)
		fmt.Println(s, "\r\n耗时：", e, "\r\n---------------")
	}
}

/*
压缩类型：zlib
压缩前大小：591716
压缩后大小：444668
压缩率：0.75148886
耗时： 107.1229ms
---------------
压缩类型：gzip
压缩前大小：591716
压缩后大小：444676
压缩率：0.7515024
耗时： 96.9619ms
---------------
压缩类型：deflate
压缩前大小：591716
压缩后大小：444666
压缩率：0.7514855
耗时： 92.618ms
---------------
压缩类型：zlib:1
压缩前大小：591716
压缩后大小：451669
压缩率：0.76332057
耗时： 93.5934ms
---------------
压缩类型：gzip:1
压缩前大小：591716
压缩后大小：451677
压缩率：0.7633341
耗时： 97.5769ms
---------------
压缩类型：deflate:1
压缩前大小：591716
压缩后大小：451667
压缩率：0.7633172
耗时： 89.1391ms
---------------
压缩类型：zlib:9
压缩前大小：591716
压缩后大小：444438
压缩率：0.7511002
耗时： 106.0908ms
---------------
压缩类型：gzip:9
压缩前大小：591716
压缩后大小：444446
压缩率：0.7511137
耗时： 103.0774ms
---------------
压缩类型：deflate:9
压缩前大小：591716
压缩后大小：444436
压缩率：0.7510968
耗时： 99.7128ms
---------------

*/
