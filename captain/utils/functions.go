package utils

import (
	"math"
	"math/rand/v2"
	"strconv"
)

func GetFilters(upstream string, country string, isSticky bool) string {

	//iproyal - username123:password321-country-dk_session-sgn34f3e_lifetime-1h@geo.iproyal.com:12321
	//netnut - USERNAME:PASSWORD-res-nl-sid-94704546@gw.netnut.net:5959

	var filter string

	switch upstream {
	case "netnut":
		filter = "-res-" + country
		if isSticky {
			filter += "-sid-" + strconv.Itoa(generateSessionNumber())
		}
	case "iproyal":
		filter = "-country-" + country
		if isSticky {
			filter += "_session-" + generateSessionNumberString() + "_lifetime-1h"
		}
	}

	return filter
}

func generateSessionNumber() int {
	return rand.IntN(89999999) + 10000000
}

func generateSessionNumberString() string {
	stringLoopRound := rand.IntN(8) + 1
	numberLoopRound := 8 - stringLoopRound
	number := rand.IntN(int(math.Pow10(numberLoopRound)))
	var str string
	for stringLoopRound >= 1 {
		str += string((rand.IntN(122-97+1) + 97))
		stringLoopRound--
	}
	runes := []rune(str + strconv.Itoa(number))
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}
