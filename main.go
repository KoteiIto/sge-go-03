package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"os"
	"sync"
)

func main() {}

// LoadFromFile はpathのファイルを読み込んで、gobのUnmarshalを行なってeを変更します
// チューニング対象
func LoadFromFile(path string, e interface{}) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if err = gob.NewDecoder(f).Decode(e); err != nil {
		return err
	}
	return nil
}

var bufferPool = sync.Pool{
	New: func() interface{} {
		return &bytes.Buffer{}
	},
}

func GetBuffer() *bytes.Buffer {
	return bufferPool.Get().(*bytes.Buffer)
}

func PutBuffer(buf *bytes.Buffer) {
	buf.Reset()
	bufferPool.Put(buf)
}

// SaveToFile はeでjsonのMarshalを行い、gzip圧縮したデータをpathに保存する
// チューニング対象
func SaveToFile(path string, e interface{}) error {
	buf := GetBuffer()
	defer PutBuffer(buf)

	gw, err := gzip.NewWriterLevel(buf, gzip.BestSpeed)
	if err != nil {
		return err
	}
	defer gw.Close()

	if err = json.NewEncoder(gw).Encode(e); err != nil {
		return err
	}

	if err = gw.Flush(); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err = f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// MergeSort はマージソートでSliceを昇順にします。
// 参考 https://www.codeflow.site/ja/article/java-merge-sort
// チューニング対象
func MergeSort(list []int) {
	mergeSort(list, make([]int, len(list)))
}

func mergeSort(list []int, sub []int) {
	l := len(list)
	if l > 1 {
		m := l / 2
		sub1 := sub[0:m]
		sub2 := sub[m:]

		mergeSort(list[0:m], sub1)
		mergeSort(list[m:], sub2)

		copy(sub1, list[0:m])
		copy(sub2, list[m:])

		merge(sub1, sub2, list)
	}
}

// チューニング非対象
// 関数の引数のインターフェイスを変えないでください。
func merge(sub1, sub2, list []int) {
	i1, i2 := 0, 0
	l1, l2 := len(sub1), len(sub2)
	for i1 < l1 || i2 < l2 {
		if i2 >= l2 || (i1 < l1 && sub1[i1] < sub2[i2]) {
			list[i1+i2] = sub1[i1]
			i1++
		} else {
			list[i1+i2] = sub2[i2]
			i2++
		}
	}
}
