//made by ramsey bissex


package main 

import "fmt"
import "net/http"

func main() {
	fmt.Println("Hello, World!")
	transport := &http.Trasport{
		MaxIdleConns:		10,
		IdleConnTimeout:	20* time.Second,
		DisableCompression:	true,
	}
	client := &http.Client{Transport: transport}
	response, err := client.Get("https://test.com/api/get?name=Aggie-1")
}

