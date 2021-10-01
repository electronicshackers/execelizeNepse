package utils

import "math"

type TechnicalData struct {
	S string    `json:"s"`
	T []int     `json:"t"`
	C []float64 `json:"c"`
	O []float64 `json:"o"`
	H []float64 `json:"h"`
	L []float64 `json:"l"`
	V []int     `json:"v"`
}

func (t TechnicalData) Diff() []float64 {
	var d []float64
	for i := 1; i < len(t.C); i++ {
		d = append(d, t.C[i]-t.C[i-1])
	}
	return d
}

func (t TechnicalData) Gains(prices []float64) []float64 {
	var g []float64
	for _, v := range prices {
		if v > 0 {
			g = append(g, v)
		} else {
			g = append(g, 0)
		}
	}
	return g
}

func (t TechnicalData) Losses(prices []float64) []float64 {
	var l []float64
	for _, v := range prices {
		if v < 0 {
			l = append(l, math.Abs(v))
		} else {
			l = append(l, 0)
		}
	}
	return l
}

func (t TechnicalData) Average(prices []float64, days int) float64 {
	var sum float64
	for _, v := range prices[0:days] {
		sum += v
	}
	return sum / float64(days)
}

func (t TechnicalData) MovingAverage(prices []float64, days int, avg float64) []float64 {
	var ma []float64
	ma = append(ma, avg)
	for _, v := range prices[days:] {
		average := (13*avg + v) / 14
		avg = average
		ma = append(ma, average)
	}
	return ma
}

func (t TechnicalData) RelativeStrength(averageLoss, averageGain []float64) []float64 {
	var rs []float64
	for i := 0; i < len(averageLoss); i++ {
		rs = append(rs, averageGain[i]/averageLoss[i])
	}
	return rs
}

func (t TechnicalData) RelativeStrengthIndicator(relativeStrength []float64) []float64 {
	var rsi []float64
	for i := 0; i < len(relativeStrength); i++ {
		rsi = append(rsi, 100-100/(1+relativeStrength[i]))
	}
	return rsi
}
