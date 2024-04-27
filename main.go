package main

import(
	"github.com/gorilla/websocket"
	"github.com/go-vgo/robotgo"
	"fmt"
	"flag"
	"net/url"
	"math"
	"time"
)

var targetIpString = flag.String("addr", "localhost", "server addres")
var originX = flag.Int("x", -1, "origin mouse position")
var originY = flag.Int("y", -1, "origin mouse position")
var mapingType = flag.Int("mapping", 1, "mapping type of the mouse (1 position, 2 movement speed)")
var maxSpeed = flag.Int("maxSpeed", 100, "used for mapping type 2")
var rate = flag.Int("rate", 50, "rate used for sampleing the mouse position in ms")

type mouse struct{
	current_x int
	current_y int
	old_x int
	old_y int
}

func (self *mouse) update(){
	self.old_x = self.current_x
	self.old_y = self.current_y
	self.current_x, self.current_y = robotgo.Location()
	fmt.Println("mouse update")
}

func (self *mouse) set(x int, y int){
	self.old_x = self.current_x
	self.old_y = self.current_y
	robotgo.Move(x, y)
}

type vector struct{
	x int
	y int
};

func (v *vector) vectorFromPoint(fromX, fromY, toX, toY int){
	v.x = toX - fromX
	v.y = toY - fromY
}
func (v *vector) length() int {
	return int(math.Floor(math.Sqrt(math.Pow(float64(v.x), 2) + math.Pow(float64(v.y), 2))))
}

func main(){
	flag.Parse()
	fmt.Println("hello word!")
//	u, err := url.Parse(fmt.Sprintf("ws://%s:81", *targetIpString))
	u, err := url.Parse("wss://echo.websocket.org/.sse")
	if err != nil {
		panic(err)
	}
	fmt.Println("server adders: ",u.String())
	
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	robotgo.MouseSleep = 0
	cursor :=  mouse{current_x :0, current_y:0}
	cursor.update()
	sx, sy := robotgo.GetScreenSize()
	if *originX == -1 || *originY == -1 {
		*originX = sx/2
		*originY = sy/2
	}
	cursor.set(*originX, *originY)
	if *mapingType == 1 {
		for{
			//distence from origine
			cursor.update()
			distenceX := cursor.current_x - *originX
			angel := 90*distenceX/(sx/2)
			fmt.Println(angel)
			time.Sleep(time.Duration(*rate) * time.Millisecond)
		}
	}else {
		cursorVec := vector{x:0,y:0}
		/*for{
			cursor.update()
			cursorVec.vectorFromPoint(cursor.current_x, cursor.current_y, cursor.old_x, cursor.old_y)
			var angle int
			if cursor.current_x - cursor.old_x >0 {
				angle = 90 * ((cursorVec.length()/ *maxSpeed))
			} else{
				angle = 90 * ((cursorVec.length()/ *maxSpeed)*-1)
			}
			fmt.Println(angle)
			time.Sleep(time.Duration(*rate) * time.Millisecond)
			cursor.set(*originX, *originY)
		}
		*/
		fmt.Println("work in progress")
	}
}
