package config

import (
	"io"
	"io/ioutil"
	"testing"
)

func TestDefaultLoader(t *testing.T) {
	cases := []struct {
		Name     string
		Instance func(data []byte) *Loader
		Validate func(reader io.Reader)
	}{
		{
			Name: "Should create a valid default loader",
			Instance: func(data []byte) *Loader {
				return DefaultLoader(data)
			},
			Validate: func(reader io.Reader) {
				content, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatalf("Failed to read data: %v", err)
				}

				if string(content[:]) != "reader" {
					t.Fatalf("Expected content: \"content\" but got \"%v\"", string(content[:]))
				}
			},
		},
	}

	for _, tc := range cases {
		b, err := ioutil.ReadFile("_testdata/dummy")
		if err != nil {
			t.Fatalf("Failed to open file: %v", err)
		}

		loader := tc.Instance(b)
		tc.Validate(loader)
	}
}

func TestFileLoader(t *testing.T) {
	cases := []struct {
		Name        string
		FilePath    string
		Instance    func(path string) (*Loader, error)
		Validate    func(reader io.Reader)
		ExpectError bool
	}{
		{
			Name:     "Should create a loader from a file path",
			FilePath: "_testdata/dummy",
			Instance: func(path string) (*Loader, error) {
				return FileLoader(path)
			},
			Validate: func(reader io.Reader) {
				content, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatalf("Failed to read data: %v", err)
				}

				if string(content[:]) != "reader" {
					t.Fatalf("Expected content: \"content\" but got \"%v\"", string(content[:]))
				}
			},
			ExpectError: false,
		},
		{
			Name:     "Should fail to create a loader from an invalid file path",
			FilePath: "_testdata/_dummy",
			Instance: func(path string) (*Loader, error) {
				return FileLoader(path)
			},
			Validate: func(reader io.Reader) {
				content, err := ioutil.ReadAll(reader)
				if err != nil {
					t.Fatalf("Failed to read data: %v", err)
				}

				if string(content[:]) != "reader" {
					t.Fatalf("Expected content: \"content\" but got \"%v\"", string(content[:]))
				}
			},
			ExpectError: true,
		},
	}

	for _, tc := range cases {
		loader, err := tc.Instance(tc.FilePath)

		if tc.ExpectError && err == nil {
			t.Fatalf("Test case was expecting an error but got none")
		}

		if !tc.ExpectError && err != nil {
			t.Fatalf("Test case failed: %v", err)
		}

		if err == nil {
			tc.Validate(loader)
		}
	}
}
