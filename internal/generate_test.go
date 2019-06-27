package internal

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenVal(t *testing.T) {
	t.Run("generate from single string", func(t *testing.T) {
		var into string
		from := `test`

		val, err := GenVal(reflect.TypeOf(into), from)
		require.NoError(t, err)
		require.Equal(t, "test", val.Interface())
	})

	t.Run("generate from number string", func(t *testing.T) {
		var into int
		from := `1`

		val, err := GenVal(reflect.TypeOf(into), from)
		require.NoError(t, err)
		require.Equal(t, 1, val.Interface())
	})

	t.Run("generate into pointer", func(t *testing.T) {
		var into *int
		from := `1`

		val, err := GenVal(reflect.TypeOf(into), from)
		require.NoError(t, err)

		into, ok := val.Interface().(*int)
		require.True(t, ok)
		require.NotNil(t, into)
		require.Equal(t, 1, *into)
	})

	t.Run("generate single string into slice", func(t *testing.T) {
		var into []string
		from := []string{"a"}

		val, err := GenVal(reflect.TypeOf(into), from[0])
		require.NoError(t, err)

		into, ok := val.Interface().([]string)
		require.True(t, ok)
		require.Equal(t, []string{"a"}, into)
	})

	t.Run("generate strings into slice", func(t *testing.T) {
		var into []string
		from := []string{"a", "b"}

		val, err := GenVal(reflect.TypeOf(into), from[0], from[1:]...)
		require.NoError(t, err)

		into, ok := val.Interface().([]string)
		require.True(t, ok)
		require.Equal(t, []string{"a", "b"}, into)
	})

	t.Run("generate strings into int slice", func(t *testing.T) {
		var into []int
		from := []string{"1", "2"}

		val, err := GenVal(reflect.TypeOf(into), from[0], from[1:]...)
		require.NoError(t, err)

		into, ok := val.Interface().([]int)
		require.True(t, ok)
		require.Equal(t, []int{1, 2}, into)
	})

	t.Run("generate strings into slice pointer", func(t *testing.T) {
		var into *[]string
		from := []string{"a", "b"}

		val, err := GenVal(reflect.TypeOf(into), from[0], from[1:]...)
		require.NoError(t, err)

		into, ok := val.Interface().(*[]string)
		require.True(t, ok)
		require.NotNil(t, into)
		require.Equal(t, []string{"a", "b"}, *into)
	})
}
