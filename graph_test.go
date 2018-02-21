package studio_statistics

import (
	"os"
	"testing"
)

func TestMakeGraph(t *testing.T) {
	in, _ := os.Open("activity.log")
	out, _ := os.Create("graph.png")
	if err := MakeGraph(in, out); err != nil {
		t.Error(err)
	}
}
