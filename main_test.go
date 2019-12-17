package main

import (
	"compress/gzip"
	"encoding/json"
	"math/rand"
	"os"
	"reflect"
	"runtime"
	"testing"
)

const (
	listSize     = 100000
	loadFilePath = "test/data.gob"
	saveFilePath = "test/data.json.gz"
)

func TestMain(m *testing.M) {
	runtime.GOMAXPROCS(1)
	m.Run()
}

func Test_LoadFromFile(t *testing.T) {
	var list []int
	if err := LoadFromFile(loadFilePath, &list); err != nil {
		t.Error(err)
	}

	if len(list) != listSize {
		t.Errorf("expected=%d, actual=%d", listSize, len(list))
	}
}

func Benchmark_LoadFromFile(b *testing.B) {
	var list []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := LoadFromFile(loadFilePath, &list); err != nil {
			b.Error(err)
		}
	}
}

func Test_SaveToFile(t *testing.T) {
	list := makeRandomList(listSize)
	if err := SaveToFile(saveFilePath, list); err != nil {
		t.Error(err)
	}

	file, err := os.Open(saveFilePath)
	if err != nil {
		t.Error(err)
	}

	gr, err := gzip.NewReader(file)
	if err != nil {
		t.Error(err)
	}

	var actualList []int
	json.NewDecoder(gr).Decode(&actualList)

	if !reflect.DeepEqual(list, actualList) {
		t.Error()
	}
}

func Benchmark_SaveToFile(b *testing.B) {
	list := makeRandomList(listSize)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := SaveToFile(saveFilePath, list); err != nil {
			b.Error(err)
		}
	}
}

func makeRandomList(size int) []int {
	list := make([]int, size)
	for i := 0; i < size; i++ {
		list[i] = rand.Int()
	}
	return list
}

func Test_MergeSort(t *testing.T) {
	list := makeRandomList(listSize)
	MergeSort(list)
	v := list[0]
	for i := range list {
		if list[i] < v {
			t.Error(list)
		}
		v = list[i]
	}
}

func Benchmark_MergeSort(b *testing.B) {
	originalList := makeRandomList(listSize)
	list := make([]int, len(originalList))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		copy(list, originalList)
		MergeSort(list)
	}
}

func Benchmark_Total(b *testing.B) {
	var list []int
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := LoadFromFile(loadFilePath, &list); err != nil {
			b.Error(err)
		}

		MergeSort(list)

		if err := SaveToFile(saveFilePath, list); err != nil {
			b.Error(err)
		}
	}
}
