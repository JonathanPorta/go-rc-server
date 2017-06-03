package gpio

import (

	"github.com/kidoman/embd"

	_ "github.com/kidoman/embd/host/all"
)

func Reset() {
       if err := embd.InitGPIO(); err != nil {
                panic(err)
        }
        defer embd.CloseGPIO()
	targetPins := []int{15, 18, 17}
	for _, targetPin := range targetPins {
		
        	embd.SetDirection(targetPin, embd.Out)
	        embd.DigitalWrite(targetPin, embd.Low)
	}
}

	

func WriteToPin(targetPin string) {
	if err := embd.InitGPIO(); err != nil {
		panic(err)
	}
//	defer embd.CloseGPIO()

	embd.SetDirection(targetPin, embd.Out)
	embd.DigitalWrite(targetPin, embd.High)
//	time.Sleep(1 * time.Second)
}
