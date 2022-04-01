package manticoresearch

import "github.com/manticoresoftware/go-sdk/manticore"
import "fmt"

func main() {
	cl := manticore.NewClient()
	cl.SetServer("127.0.0.1", 9313)
	cl.Open()
	res, err := cl.Sphinxql(`replace into testrt values(1,'my subject', 'my content', 15)`)
	fmt.Println(res, err)
	res, err = cl.Sphinxql(`replace into testrt values(2,'another subject', 'more content', 15)`)
	fmt.Println(res, err)
	res, err = cl.Sphinxql(`replace into testrt values(5,'again subject', 'one more content', 10)`)
	fmt.Println(res, err)
	res2, err2 := cl.Query("more|another", "testrt")
	fmt.Println(res2, err2)
}
