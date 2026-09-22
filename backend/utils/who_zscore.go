package utils

// Struktur acuan standar WHO per bulan (Simplifikasi Median & SD Standar WHO)
type WHOMedianSD struct {
	Median float64
	SD     float64
}

// Data acuan Standar WHO Berat Badan menurut Umur (BB/U) Laki-laki 0-12 Bulan (sampel acuan)
var whoBBUBoys = map[int]WHOMedianSD{
	0:  {Median: 3.3, SD: 0.5},
	1:  {Median: 4.5, SD: 0.6},
	2:  {Median: 5.6, SD: 0.7},
	3:  {Median: 6.4, SD: 0.8},
	4:  {Median: 7.0, SD: 0.8},
	5:  {Median: 7.5, SD: 0.8},
	6:  {Median: 7.9, SD: 0.8},
	7:  {Median: 8.3, SD: 0.9},
	8:  {Median: 8.6, SD: 0.9},
	9:  {Median: 8.9, SD: 0.9},
	10: {Median: 9.2, SD: 1.0},
	11: {Median: 9.4, SD: 1.0},
	12: {Median: 9.6, SD: 1.0},
}

// Data acuan Standar WHO Tinggi/Panjang Badan menurut Umur (TB/U) Laki-laki 0-12 Bulan
var whoTBUBoys = map[int]WHOMedianSD{
	0:  {Median: 49.9, SD: 1.9},
	1:  {Median: 54.7, SD: 2.0},
	2:  {Median: 58.4, SD: 2.1},
	3:  {Median: 61.4, SD: 2.2},
	4:  {Median: 63.9, SD: 2.2},
	5:  {Median: 65.9, SD: 2.3},
	6:  {Median: 67.6, SD: 2.3},
	7:  {Median: 69.2, SD: 2.4},
	8:  {Median: 70.6, SD: 2.4},
	9:  {Median: 72.0, SD: 2.5},
	10: {Median: 73.3, SD: 2.5},
	11: {Median: 74.5, SD: 2.5},
	12: {Median: 75.7, SD: 2.6},
}

// HitungZScoreStatusGizi menghitung Z-score BB/U dan mengembalikan kategori gizi
func HitungZScoreStatusGizi(bb float64, usiaBulan int, jenisKelamin string) string {
	ref, exists := whoBBUBoys[usiaBulan]
	if !exists {
		ref = WHOMedianSD{Median: 3.3 + (float64(usiaBulan) * 0.5), SD: 1.0}
	}

	zScore := (bb - ref.Median) / ref.SD

	switch {
	case zScore < -3.0:
		return "gizi_buruk"
	case zScore >= -3.0 && zScore < -2.0:
		return "gizi_kurang"
	case zScore >= -2.0 && zScore <= 2.0:
		return "gizi_baik"
	default:
		return "gizi_lebih"
	}
}

// HitungZScoreStatusStunting menghitung Z-score TB/U dan mengembalikan indikasi stunting
func HitungZScoreStatusStunting(tb float64, usiaBulan int, jenisKelamin string) string {
	ref, exists := whoTBUBoys[usiaBulan]
	if !exists {
		ref = WHOMedianSD{Median: 49.9 + (float64(usiaBulan) * 2.0), SD: 2.5}
	}

	zScore := (tb - ref.Median) / ref.SD

	switch {
	case zScore < -3.0:
		return "sangat_pendek"
	case zScore >= -3.0 && zScore < -2.0:
		return "pendek"
	case zScore >= -2.0 && zScore <= 3.0:
		return "normal"
	default:
		return "tinggi"
	}
}