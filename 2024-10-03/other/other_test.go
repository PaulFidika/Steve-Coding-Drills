package other

import "testing"

func TestSayHello(t *testing.T) {
	want := "Hello there, whatever\n!"
	got := Greet("whatever")

	if want != got {
		t.Errorf("wanted %s, got %s", want, got)
	}
}