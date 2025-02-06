package storage_test

import (
	"concurrency_hw1/internal/storage"
	"testing"
)

func TestEngine(t *testing.T) {
	e := storage.NewEngine()

	tests := []struct {
		name   string
		setup  func()
		action func() string
		want   string
	}{
		{
			name: "Set then Get existing key",
			setup: func() {
				e.Set("hello", "world")
			},
			action: func() string {
				return e.Get("hello")
			},
			want: "world",
		},
		{
			name: "Get key that doesn't exist",
			setup: func() {
				// no setup => the map is empty for this test
			},
			action: func() string {
				return e.Get("no_such_key")
			},
			want: "",
		},
		{
			name: "Set a key then delete it, expect empty on get",
			setup: func() {
				e.Set("delete_me", "please")
				e.Delete("delete_me")
			},
			action: func() string {
				return e.Get("delete_me")
			},
			want: "",
		},
		{
			name: "Overwrite existing key with new value",
			setup: func() {
				e.Set("foo", "oldval")
				e.Set("foo", "newval")
			},
			action: func() string {
				return e.Get("foo")
			},
			want: "newval",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e = storage.NewEngine()

			if tt.setup != nil {
				tt.setup()
			}

			got := tt.action()
			if got != tt.want {
				t.Errorf("got %q, want %q", got, tt.want)
			}
		})
	}
}
