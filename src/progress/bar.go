package progress

type Bar[T int | uint | float32 | float64] struct {
	current chan T
	total   T
	title   string
}
