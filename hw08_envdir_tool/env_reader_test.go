package main

import (
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("test ReadDir", func(t *testing.T) {
		envs, err := ReadDir("./testdata/env")
		if err != nil {
			t.Fatalf("failed to read directory: %v", err)
		}
		if envs["BAR"].Value != "bar" || envs["BAR"].NeedRemove {
			t.Errorf("unexpected content for BAR: %+v", envs["BAR"])
		}
		if !envs["EMPTY"].NeedRemove {
			t.Errorf("EMPTY might be removed")
		}
		if envs["FOO"].Value != "   foo\nwith new line" || envs["FOO"].NeedRemove {
			t.Errorf("unexpected content for FOO: %+v", envs["FOO"])
		}
		if envs["HELLO"].Value != "\"hello\"" || envs["HELLO"].NeedRemove {
			t.Errorf("unexpected content for HELLO: %+v", envs["HELLO"])
		}
		if !envs["UNSET"].NeedRemove {
			t.Errorf("UNSET might be removed")
		}
	})
}
