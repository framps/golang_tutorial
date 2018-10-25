package classes

import (
	"time"
)

func main() {
	// func TestLeds(t *testing.T) {

	lc := NewLEDController()
	//lc.Enable()
	lc.Open()

	for i := 0; i < 8; i++ {
		lc.On(i)
		time.Sleep(time.Second)
		lc.Off(i)
		time.Sleep(time.Second)
	}
	lc.Close()
}
