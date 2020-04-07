package main

import (
	"Tests_shared/malibutest"
	"log"
)

//"Tests_shared/"
//"Tests_shared/teamsbrowsertest"

// func check(err error) {
// 	if err != nil {
// 		panic(err)
// 	}
// }

func main() {
	theTest := malibutest.NewMalibuTest()
	log.Println(theTest.Params.HostName.Value)
}
