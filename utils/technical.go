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

func (t TechnicalData) ExponentialMovingAverage(prices []float64, days int, simpleAverage float64, multiplier float64) []float64 {
	var ema []float64
	ema = append(ema, simpleAverage)
	for _, v := range prices[days:] {
		average := (v-simpleAverage)*multiplier + simpleAverage
		simpleAverage = average
		ema = append(ema, average)
	}
	return ema
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

func (t TechnicalData) Multiplier(n float64) float64 {
	return (2 / (n + 1))
}

func (t TechnicalData) MovingDifference(leadingPrices []float64, trailingPrices []float64, days int) []float64 {
	var md []float64
	for index, v := range leadingPrices[days:] {
		md = append(md, v-trailingPrices[index])
	}
	return md
}

func (t TechnicalData) RSI() []float64 {
	diff := t.Diff()
	gains := t.Gains(diff)
	losses := t.Losses(diff)

	averageGain := t.Average(gains, 14)
	averageLoss := t.Average(losses, 14)

	averageGains := t.MovingAverage(gains, 14, averageGain)
	averageLosses := t.MovingAverage(losses, 14, averageLoss)

	relativeStrength := t.RelativeStrength(averageLosses, averageGains)
	rsi := t.RelativeStrengthIndicator(relativeStrength)
	return rsi
}

func (t TechnicalData) MACD() ([]float64, []float64, []float64) {
	ema12 := t.ExponentialMovingAverage(t.C, 12, t.Average(t.C, 12), t.Multiplier(12))
	ema26 := t.ExponentialMovingAverage(t.C, 26, t.Average(t.C, 26), t.Multiplier(26))

	macd := t.MovingDifference(ema12, ema26, 14)

	signal := t.ExponentialMovingAverage(macd, 9, t.Average(macd, 9), t.Multiplier(9))

	histogram := t.MovingDifference(macd, signal, 8)

	return macd, signal, histogram
}
