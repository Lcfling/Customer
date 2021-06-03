package utils

import (
	"fmt"
	"math"
	"strconv"
)

func GetDistance(lat1, lng1, lat2, lng2 string) string {
	s := GetDistanceNone(lat1, lng1, lat2, lng2)
	s = s / 10000
	var dist string
	if s < 1 {
		s = Decimal(s * 1000)
		dist = fmt.Sprintf("%.2fm", s)
	} else {
		dist = fmt.Sprintf("%.2fkm", s)
	}

	return dist
}
func GetDistanceNone(lat11, lng11, lat22, lng22 string) float64 {
	lat1, _ := strconv.ParseFloat(lat11, 64)
	lng1, _ := strconv.ParseFloat(lng11, 64)
	lat2, _ := strconv.ParseFloat(lat22, 64)
	lng2, _ := strconv.ParseFloat(lng22, 64)
	radius := 6378137.00 // 6378137
	rad := math.Pi / 180.0

	lat1 = lat1 * rad
	lng1 = lng1 * rad
	lat2 = lat2 * rad
	lng2 = lng2 * rad

	theta := lng2 - lng1
	dist := math.Acos(math.Sin(lat1)*math.Sin(lat2) + math.Cos(lat1)*math.Cos(lat2)*math.Cos(theta))

	return dist * radius
}
func Decimal(value float64) float64 {
	value, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", value), 64)
	return value
}
