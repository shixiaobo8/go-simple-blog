package utils

import (
	"fmt"
	"testing"
)

type Person struct {
	Name string
	Age  int
}

type Student struct {
	Person
	Tel   uint
	Class uint8
}

var Stu = Student{
	Person: Person{
		Name: "张三",
		Age:  13,
	},
	Tel:   18878879999,
	Class: 6,
}

func TestFieldExists(t *testing.T) {
	r, err := FieldExists(Stu, "Tel")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(r)
}

func TestGetFieldValue(t *testing.T) {
	r, err := GetFieldValue(Stu, "Person")
	if err != nil {
		t.Fatal(err)
	}

	result, ok := r.(Person)
	if ok == false {
		t.Fatal("assert error")
	}

	fmt.Printf("%T %+v\n", result, result)
	t.Log(result)
}