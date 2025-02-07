package main

import (
	"fmt"
)

func myCpuid(op uint32) (eax, ebx, ecx, edx uint32)

func main() {
	eax, _, _, _ := myCpuid(0)

	fmt.Printf("Max Std Func %d\n", eax)

}
