//made by ramsey bissex


package main 

import "fmt"
import "net/http"
import "io/ioutil"

func main() {
	fmt.Println("Hello, World!")
//	transport := &http.Trasport{
//		MaxIdleConns:		10,
//		IdleConnTimeout:	20* time.Second,
//		DisableCompression:	true,
//	}
//	client := &http.Client{Transport: transport}
	response, err := http.Get("https://httpbin.org/get")
	if err != nil {
		fmt.Println("there was an error :( :")
		fmt.Println(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		bodyBytes, err2 := ioutil.ReadAll(response.Body)
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
		if err2 != nil {
			fmt.Println("oh god")
		}
	}
}

