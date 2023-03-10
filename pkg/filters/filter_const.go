package filters

const RedisDeduplicatorPrefix = "yr:bloom"
const RedisDeduplicatorCount = "yr:bloom:count"

type FilterType = int

const (
	None FilterType = iota
	Bloom
	Exact
)

// FilterParams
type FilterParams struct {
	Kind   FilterType `json:"kind"`
	N      uint       `json:"n"`
	P      float64    `json:"p"`
	Expire uint       `json:"expire"`
}

type FilterRequest struct {
	Jobid string `json:"jobid"`
	Data  string `json:"data"`
}
