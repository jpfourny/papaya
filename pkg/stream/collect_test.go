package stream

import (
	"context"
	"testing"

	"github.com/jpfourny/papaya/internal/assert"
	"github.com/jpfourny/papaya/pkg/pair"
)

func TestCollectSlice(t *testing.T) {
	got := CollectSlice(Of(1, 2, 3))
	want := []int{1, 2, 3}
	assert.ElementsMatch(t, got, want)
}

func TestCollectMap(t *testing.T) {
	got := CollectMap(Of(
		pair.Of(1, "one"),
		pair.Of(2, "two"),
		pair.Of(3, "three"),
	))
	want := map[int]string{1: "one", 2: "two", 3: "three"}
	if len(got) != len(want) {
		t.Fatalf("got %#v, want %#v", got, want)
	}
	for k, v := range want {
		if got[k] != v {
			t.Fatalf("got[%d] %#v, want %#v", k, got, want)
		}
	}
}

func TestCollectChannel(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannel(Empty[int](), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannel(Of(1, 2, 3), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})
}

func TestCollectChannelCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannelCtx(context.Background(), Empty[int](), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		ch := make(chan int)
		go func() {
			CollectChannelCtx(context.Background(), Of(1, 2, 3), ch)
			close(ch)
		}()
		got := CollectSlice(FromChannel(ch))
		want := []int{1, 2, 3}
		assert.ElementsMatch(t, got, want)
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		ch := make(chan int)
		cancel()
		go func() {
			CollectChannelCtx(ctx, Of(1, 2, 3), ch)
			close(ch)
		}()
		// Due to race condition, we may get 0 or 1 element before cancel is seen.
		got := CollectSlice(FromChannel(ch))
		if len(got) > 1 {
			t.Errorf("expected no more than 1 element, got %#v", got)
		}
	})
}

func TestCollectChannelAsync(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := CollectChannelAsync(Empty[int](), 0) // ch closed at end of stream.
		got := CollectSlice(FromChannel(ch))
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := CollectChannelAsync(Of(1, 2, 3), 0) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := CollectChannelAsync(Of(1, 2, 3), 3) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
	})
}

func TestCollectChannelAsyncCtx(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		ch := CollectChannelAsyncCtx(context.Background(), Empty[int](), 0) // ch closed at end of stream.
		got := CollectSlice(FromChannel(ch))
		var want []int
		assert.ElementsMatch(t, got, want)
	})

	t.Run("non-empty", func(t *testing.T) {
		t.Run("unbuffered", func(t *testing.T) {
			ch := CollectChannelAsyncCtx(context.Background(), Of(1, 2, 3), 0) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
		t.Run("buffered", func(t *testing.T) {
			ch := CollectChannelAsyncCtx(context.Background(), Of(1, 2, 3), 3) // ch closed at end of stream.
			got := CollectSlice(FromChannel(ch))
			want := []int{1, 2, 3}
			assert.ElementsMatch(t, got, want)
		})
	})

	t.Run("cancelled", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		ch := CollectChannelAsyncCtx(ctx, Of(1, 2, 3), 0) // ch closed at end of stream.
		// Due to race condition, we may get 0 or 1 element before cancel is seen.
		got := CollectSlice(FromChannel(ch))
		if len(got) > 1 {
			t.Errorf("expected no more than 1 element, got %#v", got)
		}
	})
}
