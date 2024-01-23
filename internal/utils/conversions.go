package utils

func GigabytesToMegabytes(gigabyte float64) float64 {
	return gigabyte * 1024
}

func BytesToMegabytes(bytes int64) float64 {
	return float64(bytes) / (1024 * 1024)
}
