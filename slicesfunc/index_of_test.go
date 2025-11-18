package slicesfunc

import "testing"

func TestIndexOf(t *testing.T) {
	t.Run("should return the index of the first element that satisfies the condition", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		index := IndexOf(slice, func(item int) bool {
			return item == 3
		})
		if index != 2 {
			t.Errorf("expected index 2, got %d", index)
		}
	})

	t.Run("should return -1 if no element satisfies the condition", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		index := IndexOf(slice, func(item int) bool {
			return item == 6
		})
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})

	t.Run("should return -1 for an empty slice", func(t *testing.T) {
		slice := []int{}
		index := IndexOf(slice, func(item int) bool {
			return item == 1
		})
		if index != -1 {
			t.Errorf("expected index -1, got %d", index)
		}
	})

	t.Run("should work with different slice types", func(t *testing.T) {
		slice := []string{"a", "b", "c", "d", "e"}
		index := IndexOf(slice, func(item string) bool {
			return item == "d"
		})
		if index != 3 {
			t.Errorf("expected index 3, got %d", index)
		}
	})

	t.Run("should return the index of the first element if it satisfies the condition", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		index := IndexOf(slice, func(item int) bool {
			return item == 1
		})
		if index != 0 {
			t.Errorf("expected index 0, got %d", index)
		}
	})

	t.Run("should return the index of the last element if it satisfies the condition", func(t *testing.T) {
		slice := []int{1, 2, 3, 4, 5}
		index := IndexOf(slice, func(item int) bool {
			return item == 5
		})
		if index != 4 {
			t.Errorf("expected index 4, got %d", index)
		}
	})
}
