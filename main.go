package main

import(
	"github.com/gorilla/websocket"
	"github.com/go-vgo/robotgo"
	"fmt"
	"flag"
	"net/url"
)
// #cgo CFLAGS: -O3 -D LINUX=1
//#include "osSpecific.c"
import "C"

var targetIpString = flag.String("addr", "localhost", "server addres")

type mouse struct{
	current_x int
	current_y int
	old_x int
	old_y int
}

func (mouse) update(){

}

func main(){
	flag.Parse()
	fmt.Println("hello word!")
	u, err := url.Parse(fmt.Sprintf("ws://%s:81", *targetIpString))
	if err != nil {
		panic(err)
	}
	fmt.Println("server adders: ",u.String())
	
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
}
