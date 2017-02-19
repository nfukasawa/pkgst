package fixture

import (
	"net/http"
)

type Client struct {
	http.Client
}

type Color int32

const (
	None  Color = 0
	Red   Color = 1
	Green Color = 2
	Blue  Color = 3
)

var Color_name = map[int32]string{
	0: "none",
	1: "red",
	2: "green",
	3: "blue",
}
var Color_value = map[string]int32{
	"none":  0,
	"red":   1,
	"green": 2,
	"blue":  3,
}
