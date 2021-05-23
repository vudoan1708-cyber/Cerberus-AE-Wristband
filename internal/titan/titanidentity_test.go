package titan

import (
	"fmt"
	"testing"
)

// Test naming convention
// <ItemUnderTest> <expected behavior> <under what state or condition>
// e.g. HandleCategory trims LEADING spaces from valid category

func TestPrettyPrintFromValidIdentity(t *testing.T) {
	ti, err := NewIdentity(0xAAAAAAAA, 0xBBBBBBBB, 0xCCCCCCCC)
	if err != nil {
		t.Errorf("NewIdentity failed with %e", err)
	}

	fmt.Println(ti.PrettyPrint())
	t.Log(ti.PrettyPrint())

	//TODO Checks need to be added
}
