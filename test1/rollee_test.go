package rollee

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func sum(a int, b int) int { return a + b }
func mul(a int, b int) int { return a * b }

var refValues = []int{1, 2, 3, 4}

func TestFold(t *testing.T) {
	l := List{1, refValues}
	res := Fold(0, sum, l)
	assert.Equal(t, 10, res[1])
	res = Fold(1, sum, l)
	assert.Equal(t, 11, res[1])

	res = Fold(0, mul, l)
	assert.Equal(t, 0, res[1])
	res = Fold(1, mul, l)
	assert.Equal(t, 24, res[1])
}

func TestFoldChan(t *testing.T) {
	t.Run("diff ids", func(t *testing.T) {
		ch := make(chan List)

		go func() {
			ch <- List{1, refValues}
			ch <- List{2, refValues}
			close(ch)
		}()

		res := FoldChan(0, sum, ch)
		assert.Equal(t, 10, res[1])
		assert.Equal(t, 10, res[2])
	})
	t.Run("same ids", func(t *testing.T) {
		ch := make(chan List)

		go func() {
			ch <- List{1, refValues}
			ch <- List{1, refValues}
			close(ch)
		}()
		res := FoldChan(1, mul, ch)
		assert.Equal(t, 576, res[1])
	})
}

func TestFoldChanX(t *testing.T) {
	t.Run("diff ids", func(t *testing.T) {
		chs := []chan List{
			make(chan List),
			make(chan List),
			make(chan List),
		}
		go func() {
			for index, c := range chs {
				c <- List{index, refValues}
				c <- List{index, refValues}
			}
			for _, c := range chs {
				close(c)
			}
		}()
		res := FoldChanX(0, sum, chs...)
		assert.Equal(t, 20, res[0])
		assert.Equal(t, 20, res[1])
		assert.Equal(t, 20, res[2])
	})
	t.Run("same ids", func(t *testing.T) {
		chs := []chan List{
			make(chan List),
			make(chan List),
			make(chan List),
		}
		go func() {
			for _, c := range chs {
				c <- List{0, refValues}
				c <- List{0, refValues}
			}
			for _, c := range chs {
				close(c)
			}
		}()
		res := FoldChanX(0, sum, chs...)
		assert.Equal(t, 60, res[0])
	})
}
