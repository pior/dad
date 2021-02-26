package taskapi

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsString(t *testing.T) {
	s, err := asString("poipoi")

	require.NoError(t, err, "should not err for a string")
	require.Equal(t, s, "poipoi")
}

func TestAsStringWithBoolean(t *testing.T) {
	_, err := asString(false)

	require.Error(t, err, "should err for a boolean")
}

func TestAsListOfStrings(t *testing.T) {
	val, err := asListOfStrings([]interface{}{})
	require.NoError(t, err)
	require.Equal(t, []string{}, val)

	val, err = asListOfStrings([]string{})
	require.NoError(t, err)
	require.Equal(t, []string{}, val)

	val, err = asListOfStrings([]interface{}{"one", "two"})
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two"}, val)

	val, err = asListOfStrings([]string{"one", "two"})
	require.NoError(t, err)
	require.Equal(t, []string{"one", "two"}, val)
}

func TestAsListOfStringsInvalid(t *testing.T) {
	_, err := asListOfStrings(nil)
	require.Error(t, err)
	require.Equal(t, "expecting a list of strings, found a <nil> (<nil>)", err.Error())

	_, err = asListOfStrings(false)
	require.Error(t, err)
	require.Equal(t, "expecting a list of strings, found a bool (false)", err.Error())
}
