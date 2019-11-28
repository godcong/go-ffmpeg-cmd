package fftool

type Scale int

// Scale ...
const (
	Scale480P  Scale = 0
	Scale720P  Scale = 1
	Scale1080P Scale = 2
)

var bitRateList = []int64{
	//Scale480P:  1000 * 1024,
	//Scale720P:  2000 * 1024,
	//Scale1080P: 4000 * 1024,
	Scale480P:  500 * 1024,
	Scale720P:  1000 * 1024,
	Scale1080P: 2000 * 1024,
}

var frameRateList = []float64{
	Scale480P:  float64(24000)/1001 - 0.005,
	Scale720P:  float64(24000)/1001 - 0.005,
	Scale1080P: float64(30000)/1001 - 0.005,
}

type Config struct {
	Scale Scale
}

func DefaultConfig() Config {
	return Config{
		Scale: Scale720P,
	}
}
