package sticky

import (
	"bytes"
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestE2ESuccess(t *testing.T) {
	t.Parallel()

	t.Run("struct constructor", func(t *testing.T) {
		c := New()
		require.NoError(t, Register(c, Constructor(func() bytes.Buffer {
			var b bytes.Buffer
			b.WriteString("foo")
			return b
		})))

		b, err := Resolve[bytes.Buffer](c)
		require.NoError(t, err)
		assert.Equal(t, "foo", b.String())
	})

	t.Run("pointer constructor", func(t *testing.T) {
		c := New()
		require.NoError(t, Register(c, Constructor(func() *bytes.Buffer {
			var b = bytes.Buffer{}
			b.WriteString("foo")
			return &b
		})))

		b, err := Resolve[*bytes.Buffer](c)
		require.NoError(t, err)
		require.NotNil(t, b)
		assert.Equal(t, "foo", b.String())
	})

	t.Run("interface constructor", func(t *testing.T) {
		c := New()
		require.NoError(t, Register(c, Constructor(func() *bytes.Reader {
			s := "bar"
			r := bytes.NewReader([]byte(s))
			return r
		}, Implements[io.Reader]())))

		v, err := Resolve[io.Reader](c)
		require.NoError(t, err)
		b := make([]byte, 3)
		v.Read(b)
		assert.Equal(t, "bar", string(b))
	})

	t.Run("slice constructor", func(t *testing.T) {
		type A struct{ string }

		c := New()
		require.NoError(t, Register(c, Constructor(func() []A {
			return []A{{"a1"}, {"a2"}, {"a3"}}
		})))
		v, err := Resolve[[]A](c)
		require.NoError(t, err)
		assert.Equal(t, []A{{"a1"}, {"a2"}, {"a3"}}, v)
	})

	t.Run("array constructor", func(t *testing.T) {
		type A struct{ string }
		c := New()
		require.NoError(t, Register(c, Constructor(func() [3]A {
			return [3]A{{"a1"}, {"a2"}, {"a3"}}
		})))
		v, err := Resolve[[3]A](c)
		require.NoError(t, err)
		assert.Equal(t, [3]A{{"a1"}, {"a2"}, {"a3"}}, v)
	})

	t.Run("map constructor", func(t *testing.T) {
		type A struct{ string }
		c := New()
		require.NoError(t, Register(c, Constructor(func() map[string]A {
			return map[string]A{
				"a1": {"a1"},
				"a2": {"a2"},
				"a3": {"a3"},
			}
		})))
		v, err := Resolve[map[string]A](c)
		require.NoError(t, err)
		assert.Equal(t, map[string]A{
			"a1": {"a1"},
			"a2": {"a2"},
			"a3": {"a3"},
		}, v)
	})

	t.Run("channel constructor", func(t *testing.T) {
		c := New()
		require.NoError(t, Register(c, Constructor(func() chan int {
			return make(chan int, 3)
		})))
		v, err := Resolve[chan int](c)
		require.NoError(t, err)
		require.NotNil(t, v)
		assert.Equal(t, 3, cap(v))
	})

	t.Run("function constructor", func(t *testing.T) {
		type A struct{ string }
		c := New()
		require.NoError(t, Register(c, Constructor(func() func() A {
			fn := func() A {
				return A{"a1"}
			}
			return fn
		})))
		f, err := Resolve[func() A](c)
		require.NoError(t, err)
		assert.Equal(t, A{"a1"}, f())
	})

	t.Run("multi args constructor", func(t *testing.T) {
		type A struct{ string }
		type B struct{ string }
		type C struct {
			A
			B
		}
		c := New()
		require.NoError(t, Register(
			c,
			Constructor(func() A { return A{"a"} }),
			Constructor(func() B { return B{"b"} }),
			Constructor(func(a A, b B) C { return C{a, b} }),
		))
		v, err := Resolve[C](c)
		require.NoError(t, err)
		assert.Equal(t, C{A{"a"}, B{"b"}}, v)
	})

	t.Run("multi returns constructor", func(t *testing.T) {
		type A struct{ string }
		type B struct{ string }
		c := New()
		require.NoError(t, Register(c, Constructor(func() (A, *B, error) {
			return A{"a"}, &B{"b"}, nil
		})))
		a, err := Resolve[A](c)
		require.NoError(t, err)
		assert.Equal(t, A{"a"}, a)
		b, err := Resolve[*B](c)
		require.NoError(t, err)
		assert.Equal(t, &B{"b"}, b)

		c = New()
		require.NoError(t, Register(c, Constructor(func() (A, *B, error) {
			return A{"a"}, &B{"b"}, errors.New("dummy error")
		})))
		_, err = Resolve[A](c)
		assert.Equal(t, errors.New("dummy error"), err)
		_, err = Resolve[*B](c)
		assert.Equal(t, errors.New("dummy error"), err)
	})

	t.Run("nested register", func(t *testing.T) {
		type A struct{ string }
		type B struct{ A }
		type C struct{ B }
		c := New()
		require.NoError(t, Register(c, Constructor(func() A { return A{"test"} })))
		require.NoError(t, Register(c, Constructor(func(a A) B { return B{a} })))
		require.NoError(t, Register(c, Constructor(func(b B) C { return C{b} })))
		v, err := Resolve[C](c)
		require.NoError(t, err)
		assert.Equal(t, C{B{A{"test"}}}, v)
	})

	t.Run("resolve by tag", func(t *testing.T) {
		type A struct{ string }

		c := New()
		require.NoError(t, Register(
			c, Constructor(func() A {
				return A{"a1"}
			}),
			Constructor(func() A {
				return A{"a2"}
			}, Tag("test")),
		))
		v, err := Resolve[A](c)
		require.NoError(t, err)
		assert.Equal(t, A{"a1"}, v)
		v, err = Resolve[A](c, Tag("test"))
		require.NoError(t, err)
		assert.Equal(t, A{"a2"}, v)
	})

	t.Run("parameter", func(t *testing.T) {
		c := New()

		type A struct{ int }

		require.NoError(t, Register(
			c,
			Param("val1", "tag1"),
			Param(100, "tag2"),
			Param(true, "tag3"),
			Param(&A{100}, "tag4"),
		))
		p1, err := Resolve[string](c, Tag("tag1"))
		require.NoError(t, err)
		assert.Equal(t, "val1", p1)

		p2, err := Resolve[int](c, Tag("tag2"))
		require.NoError(t, err)
		assert.Equal(t, 100, p2)

		p3, err := Resolve[bool](c, Tag("tag3"))
		require.NoError(t, err)
		assert.Equal(t, true, p3)

		p4, err := Resolve[*A](c, Tag("tag4"))
		require.NoError(t, err)
		assert.Equal(t, &A{100}, p4)
	})
}

func TestE2EFailure(t *testing.T) {
	t.Parallel()

	t.Run("invalid constructor", func(t *testing.T) {
		c := New()

		var e *invalidConstructorError
		err := Register(c, Constructor("dummy"))
		assert.True(t, errors.As(err, &e))
	})

	t.Run("not found register", func(t *testing.T) {
		c := New()

		type A struct{}
		var e *notFoundRegisterError
		_, err := Resolve[A](c)
		assert.True(t, errors.As(err, &e))
	})

	t.Run("already registered", func(t *testing.T) {
		c := New()

		var e *alreadyRegisteredError
		err := Register(c, Param(100, "test"))
		require.NoError(t, err)
		err = Register(c, Param(100, "test"))
		assert.True(t, errors.As(err, &e))
	})

	t.Run("cycle dependency", func(t *testing.T) {
		c := New()

		type A struct{}
		type B struct{}
		type C struct{}
		type D struct{}

		require.NoError(t, Register(c,
			Constructor(func(d D) A { return A{} }),
			Constructor(func(a A) (B, C) { return B{}, C{} }),
		))
		var e *cycleDependencyError
		err := Register(c, Constructor(func(c C) D { return D{} }))
		assert.True(t, errors.As(err, &e))
	})
}

func TestWithContext(t *testing.T) {
	type A struct{ string }

	c := New()
	require.NoError(t, Register(c, Constructor(func() *A { return &A{"a"} })))

	ctx := context.Background()
	ctx = c.WithContext(ctx)

	v, err := Resolve[*A](ctx)
	require.NoError(t, err)
	assert.Equal(t, &A{"a"}, v)
}
