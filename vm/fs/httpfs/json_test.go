package httpfs

import (
	"fmt"
	"reflect"
	"testing"
)

type testResult[T any] struct {
	value T
	error error
}

func okResult[T any](value T) testResult[T] {
	return testResult[T]{value: value}
}

func TestFileTree_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name string
		in   string
		out  testResult[FileTreeRoot]
	}{
		{
			name: "v0-short",
			in:   `{"a": {"size": 1}}`,
			out: okResult(FileTreeRoot{
				BaseURL: "",
				Tree: map[string]FileTreeValue{
					"a": FileInfo{Size: 1},
				},
			}),
		},
		{
			name: "v0-full",
			in:   `{"$version": 0, "tree": {"a": {"size": 1}}}`,
			out: testResult[FileTreeRoot]{
				error: fmt.Errorf("json: cannot unmarshal number into Go value of type httpfs.FileInfo"),
			},
		},
		{
			name: "v1",
			in:   `{"$version": 1}`,
			out: testResult[FileTreeRoot]{
				error: fmt.Errorf("file tree version 1 is not used"),
			},
		},
		{
			name: "v2",
			in:   `{"$version": 2, "base_url": "http://example.com", "tree": {"a": {"size": 1}}}`,
			out: okResult(FileTreeRoot{
				BaseURL: "http://example.com",
				Tree: map[string]FileTreeValue{
					"a": FileInfo{Size: 1},
				},
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tree FileTreeRoot
			if err := tree.UnmarshalJSON([]byte(test.in)); err != nil {
				if test.out.error == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if err.Error() != test.out.error.Error() {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}

			if test.out.error != nil {
				t.Fatalf("expected error: %v", test.out.error)
			}

			if !reflect.DeepEqual(tree, test.out.value) {
				t.Fatalf("did not get expected value:\n"+
					"expect: %+v\n"+
					"got:    %+v", test.out.value, tree)
			}
		})
	}
}
