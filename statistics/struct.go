package statistics

type ValueInterval struct  {
	StartValue, EndValue float64
}

type Distribution struct {
	Limits ValueInterval
	Count int
}