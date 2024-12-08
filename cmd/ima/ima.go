package main

import (
	"encoding/json"
	"fmt"
	"strconv"
)

func main() {

	var b map[string][]map[string]string
	json.Unmarshal([]byte(jsdata), &b)

	fragments := b["frags"]
	firstfra := fragments[0]
	xMin, _ := strconv.Atoi(firstfra["left"])
	xMax := xMin
	yMin, _ := strconv.Atoi(firstfra["top"])
	yMax := yMin

	for _, fra := range fragments {
		lefty, _ := strconv.Atoi(fra["left"])
		toppy, _ := strconv.Atoi(fra["top"])
		if xMin > lefty {
			xMin = lefty
		}
		if xMax < lefty {
			xMax = lefty
		}
		if yMin > toppy {
			yMin = toppy
		}
		if yMax < toppy {
			yMax = toppy
		}
	}

	fmt.Printf("X %d %d, Y %d %d\n%v", xMin, xMax, yMin, yMax, firstfra)

}
