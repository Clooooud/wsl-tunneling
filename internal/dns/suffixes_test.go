package dns

import (
	"reflect"
	"testing"
)

func TestNormalizeSearchSuffixes(t *testing.T) {
	got := NormalizeSearchSuffixes([]string{" cr2.st.com ", "CR2.ST.COM", "bad suffix", ".cro.st.com.", ""})
	want := []string{"cr2.st.com", "cro.st.com"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("NormalizeSearchSuffixes() = %#v, want %#v", got, want)
	}
}
