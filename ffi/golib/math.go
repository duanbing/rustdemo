package main

import "C"
import ()

//export Add
func Add(a int64, b int64) int64 {
	return a + b
}

//export AddArray
func AddArray(a []int64, b []int64, c *[]int64) {
	println("len a", len(a))
	println("len b", len(b))
	for i := 0; i < len(a); i++ {
		(*c)[i] = a[i] + b[i]
		println(a[i], b[i], (*c)[i])
	}
}

//export AddString
func AddString(a string, b string) *C.char {
	c := a + b
	println(c)
	return C.CString(c)
}

//export AddMultiRet
func AddMultiRet() (*C.char, int64) {
	return C.CString("nihaop!"), 7
}

func main() {}
