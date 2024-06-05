package main

import (
	"flag"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/gorilla/websocket"
	"github.com/veandco/go-sdl2/sdl"
)

var targetIpString = flag.String("addr", "localhost", "server addres")
var port = flag.Int("p", 81, "port number (default is the 81)")
var originX = flag.Int("x", -1, "origin mouse position (default is the center)")
var originY = flag.Int("y", -1, "origin mouse position (default is the center)")
var mapingType = flag.Int("mapping", 1, "mapping type of the mouse (1 position, 2 movement speed)")
var rate = flag.Int("rate", 50, "rate used for sampleing the mouse position in ms")

type mouse struct {
	current_x int
	current_y int
	old_x     int
	old_y     int
}

func (self *mouse) update() {
	self.old_x = self.current_x
	self.old_y = self.current_y
	self.current_x, self.current_y = robotgo.Location()
	fmt.Println("mouse update")
}

func (self *mouse) set(x int, y int) {
	self.old_x = self.current_x
	self.old_y = self.current_y
	robotgo.Move(x, y)
	self.current_x, self.current_y = robotgo.Location()
}

func (self *mouse) difference() (int, int) {
	return self.current_x - self.old_x, self.current_y - self.old_y
}

func GetAngelfunc() func() int {
	robotgo.MouseSleep = 0
	cursor := mouse{current_x: 0, current_y: 0}
	cursor.update()
	sx, sy := robotgo.GetScreenSize()
	if *originX == -1 || *originY == -1 {
		*originX = sx / 2
		*originY = sy / 2
	}
	cursor.set(*originX, *originY)

	if *mapingType == 1 {
		return func() int {
			//distence from origin
			cursor.update()
			distenceX := cursor.current_x - *originX
			angle := 90 * distenceX / (sx / 2)
			//fmt.Println(angle)
			return angle + 90
		}
	} else {
		maxVelocity := 0.0
		return func() int {
			cursor.update()
			diffX, diffY := cursor.difference()
			// V = s/T
			// T=50ms (rate)
			//s = diffX
			velocity := float64(diffX / *rate)
			if velocity != 0 {
				if math.Abs(velocity) > maxVelocity {
					maxVelocity = math.Abs(velocity)
				}
			}
			angle := -1
			if maxVelocity != 0 {
				angle = int(90 * velocity / maxVelocity)
			}
			//fmt.Println(angle)
			if cursor.current_x < sx/3 || cursor.current_x > sx-sx/3 {
				cursor.set(*originX, *originY)
				cursor.old_x = cursor.current_x + diffX
				cursor.old_y = cursor.current_y + diffY
			}
			return angle + 90
		}
	}
}

func main() {
	flag.Parse()
	fmt.Println("hello word!")
	//	u, err := url.Parse(fmt.Sprintf("ws://%s:%d", *targetIpString, *port))
	u, err := url.Parse("wss://echo.websocket.org/.sse")
	if err != nil {
		panic(err)
	}
	fmt.Println("server adders: ", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		panic(err)
	}
	defer c.Close()
	sdl.Init(sdl.INIT_EVERYTHING)
	window, err := sdl.CreateWindow("title", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, 300, 300, sdl.WINDOW_SHOWN)
	renderer, err := sdl.CreateRenderer(window, -1, 0)
	err = renderer.SetDrawColor(0, 0, 0, 255)
	err = renderer.Clear()
	renderer.Present()
	if err != nil {
		panic(err)
	}
	angelUpdater := GetAngelfunc()
	for {
		sdl.PumpEvents()
		keyboard := sdl.GetKeyboardState()
		direction := "s"
		if keyboard[sdl.SCANCODE_W] || keyboard[sdl.SCANCODE_S] == true {
			if keyboard[sdl.SCANCODE_W] == true {
				direction = "f"
			} else {
				direction = "b"
			}
		}
		err := c.WriteMessage(websocket.TextMessage, []byte(strconv.Itoa(angelUpdater())+";"+direction+"\r\n"))
		if err != nil {
			panic(err)
		}
		time.Sleep(time.Duration(*rate) * time.Millisecond)
	}
}
