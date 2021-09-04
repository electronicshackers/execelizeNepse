package utils

import (
	"errors"
	"fmt"
	"math"
	"nepse-backend/nepse"
	"net/http"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
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

func CreateExcelFile(folderName, fileName string, headers map[string]string, data []map[string]interface{}) {
	f := excelize.NewFile()
	for k, v := range headers {
		f.SetCellValue("Sheet1", k, v)
	}

	for _, vals := range data {
		for k, v := range vals {
			f.SetCellValue("Sheet1", k, v)
		}
	}

	if err := f.SaveAs(fmt.Sprintf("%s/%s.xlsx", folderName, fileName)); err != nil {
		fmt.Println(err)
	}

}

func GetDateRange(w http.ResponseWriter, start, end string) ([]string, error) {
	var days []string
	if start == "" || end == "" {
		return days, errors.New("start and end date are required")
	}

	// Change String to Date
	startDate, err := StringToTime(start)
	if err != nil {
		return days, errors.New("start date is invalid")
	}

	startDay := startDate.Weekday().String()
	if startDay == "Friday" || startDay == "Saturday" {
		return days, errors.New("start date should be a weekday")
	}

	endDate, err := StringToTime(end)
	if err != nil {
		return days, errors.New("end date is invalid")
	}

	endDay := endDate.Weekday().String()

	if endDay == "Friday" || endDay == "Saturday" {
		return days, errors.New("end date should be a weekday")
	}

	// find the difference in days between start and end date
	diffDays := endDate.Sub(startDate).Hours() / 24
	if diffDays < 0 {
		return days, errors.New("start date must be before end date")
	}
	if diffDays > 191 {
		return days, errors.New("start date must be less than 65 Nepse Days before end date")
	}

	// for loop
	// declare a variable with array of string
	for i := 0; i <= int(diffDays); i++ {
		addedDate := startDate.Add(time.Hour * 24 * time.Duration(i)).Format("2006-01-02")
		if addedDate != endDate.Format("2006-01-02") {
			days = append(days, addedDate)
		}
	}
	days = append(days, end)
	return days, nil
}

func MinMax(array []float64) (float64, float64) {
	if len(array) == 0 {
		return 0, 0
	}
	var max float64 = array[0]
	var min float64 = array[0]
	for _, value := range array {
		if max < value {
			max = value
		}
		if min > value {
			min = value
		}
	}
	return min, max
}
