package util

import (
	"math"
	"math/rand"
	"time"
)

//GetRandomPositiveInt is a wrapper for rand.Intn(int), generating a seed each time it's called.
//If 0 is passed as input, rand.Int() is called instead.
func GetRandomPositiveInt(maxValue int)int{
	goodInput := int(math.Abs(float64(maxValue)))
	rand.Seed(time.Now().UnixNano())
	if goodInput > 0{
		return rand.Intn(goodInput)
	}else{
		return rand.Int()
	}
}
