package utils

func PerThousand(curr, max int) int {
	var currF = float64(curr) + 1
	var maxF = float64(max)

	var perc = currF / maxF

	perc *= 1000

	return int(perc)
}
