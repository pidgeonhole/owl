package owl

import (
	"testing"
	"strings"
	"log"
	"fmt"
)

func TestCheckTrue(t *testing.T) {
	output := strings.NewReader(`1
4
9
`)
	ans := strings.NewReader(`1
4
9
`)
	types := []int{INT}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	if !result {
		t.Fail()
	}
}

func TestCheckFalse(t *testing.T) {
	output := strings.NewReader(`1
4
10
`)
	ans := strings.NewReader(`1
4
9
`)
	types := []int{INT}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	if result {
		t.Fail()
	}
}

func TestCheckTwoTokensPass(t *testing.T) {
	output := strings.NewReader(`1 4
4 9
9 16
`)
	ans := strings.NewReader(`1 4
4 9
9 16
`)
	types := []int{INT, INT}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	if !result {
		t.Fail()
	}
}

func TestCheckTwoTokensFail(t *testing.T) {
	output := strings.NewReader(`1 5
4 10
9 17
`)
	ans := strings.NewReader(`1 4
4 9
9 16
`)
	types := []int{INT, INT}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	if result {
		t.Fail()
	}
}

func TestCheckIntFloatPass(t *testing.T) {
	output := strings.NewReader(`1 4.0
4 9.0
9 16.0
`)
	ans := strings.NewReader(`1 4.0
4 9.0
9 16.0
`)
	types := []int{INT, FP}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}

func TestCheckIntFloatTooShort(t *testing.T) {
	output := strings.NewReader(`1 4.0
4 9.0
`)
	ans := strings.NewReader(`1 4.0
4 9.0
9 16.0
`)
	types := []int{INT, FP}

	result, err := Check(ans, output, types)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(result)
}