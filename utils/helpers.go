package utils

import (
	"fmt"
	"math"
	"nepse-backend/nepse"
	"time"
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

func StringToTime(str string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", str)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

func MapColumns(columns []string) []string {
	sectorMap := map[string]string{
		"hydro":    "Hydro Power",
		"org":      "Organized Fund",
		"life":     "Life Insurance",
		"micro":    "Microcredit",
		"dev-bank": "Development Bank",
		"hotel":    "Hotels",
		"non-life": "Non Life Insurance",
		"finance":  "Finance",
		"bank":     "Commercial Banks",
		"trading":  "Trading",
		"manu":     "Manufacturing And Processing",
		"telecom":  "Telecom",
	}
	var result []string
	for _, column := range columns {
		result = append(result, sectorMap[column])
	}
	return result
}
