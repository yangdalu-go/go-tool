package structcopy

import (
	"fmt"
	"log"
	"testing"
)

type StructA struct {
	Name string `json:"name"`

	Age int `json:"gender"`

	InUse bool `json:"in_use"`
}

type StructB struct {
	Name string `json:"name"`

	Age int `json:"gender"`

	InUse bool `json:"in_use"`
}

func TestJsonDeepCopy(t *testing.T) {

	a := StructA{
		Name:  "测试名字",
		Age:   18,
		InUse: true,
	}

	b := StructB{}

	err := JsonDeepCopy(a, &b)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(a)
	fmt.Println(b)
}
