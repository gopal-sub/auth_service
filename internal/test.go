package main

import (
	"fmt"
)

type User struct {
	name string
	age int
}





func main(){
	var users = make(map[string]User)

 	users["gopal"] = User{name: "gopal", age: 12}
	fmt.Println(users["gopal"].name)


}