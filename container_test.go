package sticky

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveCache(t *testing.T) {
	t.Parallel()

	type A struct{ int }

	tests := []struct {
		name   string
		input  register
		assert func(assert.TestingT, bool, ...any) bool
	}{
		{
			name:   "no option",
			input:  Constructor(func() *A { return &A{} }),
			assert: assert.True,
		},
		{
			name:   "cache true",
			input:  Constructor(func() *A { return &A{} }, Cache(true)),
			assert: assert.True,
		},
		{
			name:   "cache false",
			input:  Constructor(func() *A { return &A{} }, Cache(false)),
			assert: assert.False,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newContainer()
			require.NoError(t, c.Register(tt.input))
			keys, err := tt.input.Keys()
			require.NoError(t, err)
			v1, err := c.Resolve(keys[0])
			require.NoError(t, err)
			keys, err = tt.input.Keys()
			require.NoError(t, err)
			v2, err := c.Resolve(keys[0])
			require.NoError(t, err)
			t.Logf("v1=%p\n", v1)
			t.Logf("v2=%p\n", v2)
			tt.assert(t, fmt.Sprintf("%p", v1) == fmt.Sprintf("%p", v2))
		})
	}

}

func TestExtract(t *testing.T) {
	c := newContainer()

	type A struct{ string }
	type B struct{ A }
	require.NoError(t, c.Register(Constructor(func() A { return A{"test"} })))
	require.NoError(t, c.Register(Constructor(func(a A) *B { return &B{a} })))

	require.NoError(t, c.Extract(func(a A, b *B) {
		assert.Equal(t, A{"test"}, a)
		assert.Equal(t, &B{A{"test"}}, b)
	}))
}

func TestValidate(t *testing.T) {
	t.Parallel()

	type A struct{ string }
	type B struct{ A }
	type C struct{ *B }

	t.Run("success", func(t *testing.T) {
		c := newContainer()

		for _, f := range []any{
			func() A { return A{"test"} },
			func(a A) *B { return &B{a} },
			func(b *B) C { return C{b} },
		} {
			require.NoError(t, c.Register(Constructor(f)))
		}
		require.NoError(t, c.Validate())
	})

	t.Run("failure", func(t *testing.T) {
		c := newContainer()

		for _, f := range []any{
			func() A { return A{"test"} },
			func(a A) *B { return &B{a} },
			func(b B) C { return C{&b} },
		} {
			require.NoError(t, c.Register(Constructor(f)))
		}
		require.Error(t, c.Validate())
	})
}

func TestDecorate(t *testing.T) {
	c := newContainer()

	type A struct{ string }

	ini := Constructor(func() A { return A{"value"} })
	require.NoError(t, c.Register(ini))

	keys, err := ini.Keys()
	require.NoError(t, err)
	v, err := c.Resolve(keys[0])
	require.NoError(t, err)
	assert.Equal(t, "value", v.(A).string)

	keys, err = ini.Keys()
	require.NoError(t, err)
	require.NoError(t, c.Decorate(keys[0], func(v any) (any, error) {
		a, ok := v.(A)
		require.True(t, ok)
		a.string = "decorated"
		return any(a), nil
	}))
	keys, err = ini.Keys()
	require.NoError(t, err)
	v, err = c.Resolve(keys[0])
	require.NoError(t, err)
	assert.Equal(t, "decorated", v.(A).string)
}
