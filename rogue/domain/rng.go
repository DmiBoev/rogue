// для тестов(абстракция)
package domain

import "math/rand"

type RNG interface {
	Intn(n int) int
}

type DefaultRNG struct {
	r *rand.Rand
}

func NewDefaultRNG(seed int64) *DefaultRNG {
	return &DefaultRNG{r: rand.New(rand.NewSource(seed))}
}

func (d *DefaultRNG) Intn(n int) int {
	if n <= 0 {
		return 0
	}
	return d.r.Intn(n)
}
