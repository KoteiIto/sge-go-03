package main

import (
	"bytes"
	"compress/gzip"
	"encoding/gob"
	"encoding/json"
	"io/ioutil"
	"os"
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

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(b)
	if err = gob.NewDecoder(buf).Decode(e); err != nil {
		return err
	}

	return nil
}

// SaveToFile はeでjsonのMarshalを行い、gzip圧縮したデータをpathに保存する
// チューニング対象
func SaveToFile(path string, e interface{}) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	b, err := json.Marshal(e)
	if err != nil {
		return err
	}

	buf := bytes.Buffer{}
	gw := gzip.NewWriter(&buf)
	defer gw.Close()

	if _, err = gw.Write(b); err != nil {
		return err
	}

	if err = gw.Flush(); err != nil {
		return err
	}

	if _, err = f.Write(buf.Bytes()); err != nil {
		return err
	}

	return nil
}

// MergeSort はマージソートでSliceを昇順にします。
// 参考 https://www.codeflow.site/ja/article/java-merge-sort
// チューニング対象
func MergeSort(list []int) {
	l := len(list)
	if l > 1 {
		m := l / 2
		n := l - m
		sub1 := make([]int, m)
		sub2 := make([]int, n)

		for i := 0; i < m; i++ {
			sub1[i] = list[i]
		}
		for i := 0; i < n; i++ {
			sub2[i] = list[i+m]
		}

		MergeSort(sub1)
		MergeSort(sub2)
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
