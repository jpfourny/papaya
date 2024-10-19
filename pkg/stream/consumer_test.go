package stream

import (
	"slices"
	"testing"
)

func Test_stopSensingConsumer(t *testing.T) {
	t.Run("stream-exhausted", func(t *testing.T) {
		s := Of(1, 2, 3)
		var saw []int
		c1 := func(e int) bool {
			saw = append(saw, e)
			return true
		}
		c2, stopped := stopSensingConsumer(c1)
		s(c2)
		if *stopped {
			t.Errorf("expected stopped to be false; got true")
		}
		want := []int{1, 2, 3}
		if !slices.Equal(saw, want) {
			t.Errorf("expected to see %v; got %v", want, saw)
		}
	})

	t.Run("consumer-stopped", func(t *testing.T) {
		s := Of(1, 2, 3)
		var saw []int
		c1 := func(e int) bool {
			saw = append(saw, e)
			return false
		}
		c2, stopped := stopSensingConsumer(c1)
		s(c2)
		if !*stopped {
			t.Errorf("expected stopped to be true; got false")
		}
		want := []int{1}
		if !slices.Equal(saw, want) {
			t.Errorf("expected to see %v; got %v", want, saw)
		}
	})
}
