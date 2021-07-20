package utils

import (
	"fmt"
	"math"
	"nepse-backend/nepse"
)

func ToFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func GetColumn(column string, num int) string {
	return fmt.Sprintf("%s%d", column, num+2)
}

func CalculateGrahamValue(eps, bookValue float64) float64 {
	return math.Sqrt(22.5 * eps * bookValue)
}

func SetIntersection(s1, s2 []nepse.Ticker) (inter []nepse.Ticker) {
	hash := make(map[string]bool)
	for _, e := range s1 {
		hash[e.Ticker] = true
	}
	for _, e := range s2 {
		// If elements present in the hashmap then append intersection list.
		if hash[e.Ticker] {
			inter = append(inter, e)
		}
	}
	fmt.Println("inter", inter)
	return removeDups(inter)
}

//Remove dups from slice.
func removeDups(elements []nepse.Ticker) (nodups []nepse.Ticker) {
	encountered := make(map[string]bool)
	for _, element := range elements {
		if !encountered[element.Companyname] {
			nodups = append(nodups, element)
			encountered[element.Companyname] = true
		}
	}
	return
}
