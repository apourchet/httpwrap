package httpwrap

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUtils(t *testing.T) {
	t.Run("are types unique", func(t *testing.T) {
		fn := func(err error) {}
		types, _ := typesOf(fn)
		err := areTypesUnique(types)
		require.NoError(t, err)

		fn1 := func(err1, err2 error) {}
		types, _ = typesOf(fn1)
		err = areTypesUnique(types)
		require.Error(t, err)
	})
}
