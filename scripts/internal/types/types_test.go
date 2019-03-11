package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRepo_SetTags(t *testing.T) {
	repo := &Repo{Host: "abc", Name: "xyz", Language: "Go"}

	cases := []struct {
		in  []string
		out []string
	}{
		{nil, []string{"abc", "code", "go"}},
		{[]string{"A"}, []string{"a", "abc", "code", "go"}},
		{[]string{"ABC"}, []string{"abc", "code", "go"}},
		{[]string{"D e F", " \t\n", "ABC"}, []string{"abc", "code", "d e f", "go"}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.in), func(t *testing.T) {
			repo.Tags = nil
			repo.SetTags(c.in...)

			require.Equal(t, c.out, repo.Tags)
		})
	}
}

func TestPost_SetTags(t *testing.T) {
	post := &Post{Website: "abc", Title: "xyz"}

	cases := []struct {
		in  []string
		out []string
	}{
		{nil, []string{"abc", "post"}},
		{[]string{"A"}, []string{"a", "abc", "post"}},
		{[]string{"ABC"}, []string{"abc", "post"}},
		{[]string{"D e F", "", "ABC"}, []string{"abc", "d e f", "post"}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.in), func(t *testing.T) {
			post.Tags = nil
			post.SetTags(c.in...)

			require.Equal(t, c.out, post.Tags)
		})
	}
}

func TestMicroPost_SetTags(t *testing.T) {
	post := &MicroPost{Website: "abc", Text: "xyz"}

	cases := []struct {
		in  []string
		out []string
	}{
		{nil, []string{"abc", "micro", "post"}},
		{[]string{"A"}, []string{"a", "abc", "micro", "post"}},
		{[]string{"ABC"}, []string{"abc", "micro", "post"}},
		{[]string{"D e F", "", "ABC"}, []string{"abc", "d e f", "micro", "post"}},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.in), func(t *testing.T) {
			post.Tags = nil
			post.SetTags(c.in...)

			require.Equal(t, c.out, post.Tags)
		})
	}
}
