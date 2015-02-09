package tools

import (
	"testing"
)

func TestIso88591ToUtf8(t *testing.T) {

	check := func(b []byte, s string) {
		res := iso88591ToUtf8(b)
		if res != s {
			t.Errorf("Expected '%s', got '%s'", s, res)
		}
	}

	check([]byte{0x55}, "U")
	check([]byte{0x46, 0x6F, 0x6F}, "Foo")
	check([]byte{0x61, 0x73, 0x74, 0xe9, 0x72, 0x6F, 0xEF, 0x64, 0x65}, "astéroïde")
}

func TestInArray(t *testing.T) {

	check := func(a []string, e string, b bool) {
		res := InArray(a, e)
		if res != b {
			t.Errorf("Expected '%t', got '%t'", b, res)
		}
	}

	check(nil, "", false)
	check([]string{"a", "b", "c"}, "", false)
	check([]string{"a", "b", "c"}, "d", false)
	check([]string{"a", "b", "c"}, "abc", false)
	check([]string{"a", "b", "c"}, "a", true)
	check([]string{"a", "b", "c"}, "b", true)
	check([]string{"a", "b", "c"}, "c", true)
	check([]string{"a"}, "a", true)
	check([]string{"a"}, "b", false)
}

func TestCheckFilePath(t *testing.T) {

	check := func(p string, b bool) {
		res := CheckFilePath(p)
		if res != b {
			t.Errorf("Expected '%t', got '%t'", b, res)
		}
	}

	check("/dev/null", true)
	check("", false)
}

func TestGetSlugFromString(t *testing.T) {

	check := func(f string, s string) {
		res := GetSlugFromString(f)
		if res != s {
			t.Errorf("Expected '%s', got '%s'", s, res)
		}
	}

	check("", "")
	check("The Big Bang Theory S01E01", "the_big_bang_theory_s01e01")
	check("Mary à tout prix", "mary_a_tout_prix")
	check("L'Alabama", "l_alabama")
	check("Laurel & Hardy", "laurel_hardy")
}
