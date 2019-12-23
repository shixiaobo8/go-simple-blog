package utils

import (
	"github.com/astaxie/beego/logs"
	"log"
	"testing"
)

func TestSliceRemoveDupString(t *testing.T) {
	s := []string{"John", "Bill", "Bill", "Gary", "Alice"}
	r := SliceRemoveDupString(s)
	logs.Info(r)
}

func TestSliceStructColumn(t *testing.T) {
	s := []Person{
		{
			Name: "John",
			Age:  10,
		},
		{
			Name: "Daniel",
			Age:  20,
		},
		{
			Name: "Smith",
			Age: 30,
		},
	}

	si := make([]interface{}, len(s))
	for i, v := range s {
		si[i] = interface{}(v)
	}

	columns, err := SliceStructColumn(si, "Name")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%v\n", columns)
}

func TestSliceIntersection(t *testing.T) {
	s := []string{
		"John",
		"Daniel",
		"Smith",
	}
	si := make([]interface{}, len(s))
	for i, v := range s {
		si[i] = interface{}(v)
	}

	ts := []string{
		"Daniel",
		"Alice",
		"Ross",
	}
	ti := make([]interface{}, len(ts))
	for i, v := range ts {
		ti[i] = interface{}(v)
	}

	logs.Info(SliceIntersection(si, ti))


	m1 := []string{"a", "b", "c", "d", "e"}
	m2 := []string{"b", "h", "a", "c", "m", "q", "s", "e"}

	var res []string
	var count int
	for _, item := range m1 {
		for k, value := range m2 {
			count++
			if item == value {
				res = append(res, item)
				m2 = append(m2[0:k], m2[k+1:]...)
				break
			}
		}
	}
	log.Println("循环：", count)
	log.Printf("%+v\n", res)
}
