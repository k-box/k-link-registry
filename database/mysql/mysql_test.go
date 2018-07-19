package mysql

import (
	"fmt"
	"testing"
)

func TestDatabase(t *testing.T) {
	db, err := NewDatabase("kregistry:kregistry@tcp(127.0.0.1)/kregistry")
	if err != nil {
		panic(err)
	}

	regs, err := db.ListRegistrants()
	if err != nil {
		panic(err)
	}

	for _, reg := range regs {
		fmt.Println(reg.Name, string(reg.Password))
		fmt.Println(reg.CheckPass("1234567890"))
	}

}
