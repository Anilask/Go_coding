// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"io/ioutil"
// 	"net/http"

// 	"github.com/gorilla/mux"
// )

// type ResponseStruct struct{
// 	Headers interface{} `json:"headers"`
// }

// func main() {
// 	r:=mux.NewRouter()
// 	r.HandleFunc("/test",func(w http.ResponseWriter,r *http.Request){
// 		fmt.Fprintf(w,"hello World")

// 	url:="https://httpbin.org/get"
// 	req,err:= http.NewRequest("GET",url,nil)
// 	if err!=nil{
// 		fmt.Println("error while creating request %v",err)
// 	}
// 	res,err:= http.DefaultClient.Do(req)
// 	if err!=nil{
// 		fmt.Println("error while reading hitting API %v",err)
// 	}

// 	body,readErr:= ioutil.ReadAll(res.Body)
// 	if readErr!=nil{
// 		fmt.Println("error while reading response %v",err)
// 	}

// 	var bodyResp ResponseStruct
// 	err=json.Unmarshal(body,&bodyResp)

// 	fmt.Println(bodyResp)
// 	fmt.Fprintf(w,bodyResp.Headers.(string))
// 	})
// 	http.ListenAndServe(":8080",r)

// }

// package main

// import (
// 	"fmt"
// )

// func test(ch chan int, quit chan bool) {
// 	x, y := 0, 1
// 	for {
// 		select {
// 		case ch <- x:
// 			x, y = y, x+y
// 		case <-quit:
// 			fmt.Println("quit")
// 			return
// 		}
// 	}
// }
// func main() {
// 	ch := make(chan int)
// 	quit := make(chan bool)
// 	n := 10
// 	go func(n int) {
// 		for i := 0; i < n; i++ {
// 			fmt.Println(<-ch)
// 		}
// 		quit <- false
// 	}(n)
// 	test(ch, quit)
// }
// //0,1,1,2,3,5,8,15,

package main

import (
	"fmt"
)

func main() {
	fmt.Println("Started Main")
	ch := make(chan int)
	go sampleRoutine(ch)
	fmt.Println(<-ch)
	
	fmt.Println("Finished Main")
}
func sampleRoutine(ch chan int) {
 close(ch)

	fmt.Println("Inside Sample Goroutine")
	ch <- 1
	 
}
