package vincireader

import (
    //"bufio"
    "fmt"
    //"io"
    //"io/ioutil"
    "os"
)
func check(e error) {
    if e != nil {
        panic(e)
    }
}

func ReadFile(){
	f, err := os.Open("../../../C++/conversation2.lca")
	check(err)
	b1 := make([]byte, 1024)
    n1, err := f.Read(b1)
    for n1 != 0 {
		fmt.Printf("%d bytes: %s\n", n1, string(b1))
		n1, err = f.Read(b1)
	}
	f.Close()
}