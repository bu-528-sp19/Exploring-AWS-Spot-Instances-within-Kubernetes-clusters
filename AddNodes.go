package main

import (
	"fmt"
	"os/exec"
	"reflect"
)



func main() {


	out,err := exec.Command("/bin/sh","nodeaddition.sh").Output()


	if err != nil {
		//error := string(err[:])
		fmt.Println("The type of err",reflect.TypeOf(err))
		fmt.Println("error %s", err)
	}
	// as the out variable defined above is of type []byte we need to convert
	// this to a string or else we will see garbage printed out in our console
	// this is how we convert it to a string
	//fmt.Println("Command Successfully Executed")
	output := string(out[:])
	fmt.Printf(output)
	//
	//out1,err1 := exec.Command("python", "addnodes.py").Output()
	//
	//if err1 != nil {
	//	//error := string(err[:])
	//	fmt.Println("The type of err",reflect.TypeOf(err1))
	//	fmt.Println("error %s", err1)
	//}
	//
	//output1 := string(out1[:])
	//fmt.Printf(output1)
	//
	//
	//out2,err2 := exec.Command("kops","replace", "ig", "filename=nodes.yml").Output()
	//
	//if err2 != nil {
	//	//error := string(err[:])
	//	fmt.Println("The type of err",reflect.TypeOf(err2))
	//	fmt.Println("error %s", err2)
	//}
	//
	//output2 := string(out2[:])
	//fmt.Printf(output2)
	//
	//out3,err3 := exec.Command("kops", "update", "cluster", "yes").Output()
	//
	//if err3 != nil {
	//	//error := string(err[:])
	//	fmt.Println("The type of err",reflect.TypeOf(err3))
	//	fmt.Println("error %s", err3)
	//}
	//
	//output3 := string(out3[:])
	//fmt.Printf(output3)
	//
	//
	////fmt.Printf("output ",)
	////if runtime.GOOS == "windows" {
	////	cmd = exec.Command("tasklist")
	////	fmt.Println("Working windows")
	//
	//
	////}
	////fmt.Println("Working")
	////err := cmd.Run()
	//
	//
	//
	////if cmd != nil {
	//	//log.Fatalf("cmd.Run() failed with %s\n", err)
	////}
	////fmt.Printf("error ",err)
	////fmt.Printf("output ",)
}
