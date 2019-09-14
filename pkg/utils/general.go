package utils

import (
	"math"
	"math/rand"
	"time"
)

func GetRandomPositiveInt(maxValue int)int{
	goodInput := int(math.Abs(float64(maxValue)))
	rand.Seed(time.Now().UnixNano())
	if goodInput > 0{
		return rand.Intn(goodInput)
	}else{
		return rand.Int()
	}
}
