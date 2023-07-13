package main

import "fmt"
func Rev(str string) string {
	runes :=[]rune(str)

	for i,j :=0,len(runes)-1;i<j;i,j=i+1,j-1{
		runes[i],runes[j] = runes[j],runes[i]
	}

	return string(runes)
}
func rev(str string) (result string) {
	for _,v :=range str{
        result = string(v) + result
	}
	return
}
func main(){
	str:="Anil"
	rever :=Rev(str)
	reverse :=rev(str)
	fmt.Println(rever)
	fmt.Println(reverse)
}
