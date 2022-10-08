package utils

import (
	"fmt"
)

func CheckErr(err error) {
	if err != nil {
		fmt.Printf("errors : \n")
		fmt.Println(err)
	}
}
