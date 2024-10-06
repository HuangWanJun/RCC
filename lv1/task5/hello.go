package main

import "fmt"

type rect struct {
    width, height int
}

func (r *rect) area() int {
    return r.width * r.height
}

func (r rect) perim() int {
    return 2*r.width + 2*r.height
}

func main() {
    
	s := "Yes我喜欢你"
	fmt.Println(s)
	for _, c := range []byte(s) {
	    fmt.Printf("%c ", c) // 输出unicode
		fmt.Printf("%x ", c) // 输出16进制
	}
	fmt.Println()
	for i, c := range s {
	    fmt.Printf("%d %X ",i,c)
	}
	 strings.
}