package html_test

import (
	"testing"

	"github.com/unders/mockingbird/server/pkg/errs"

	"github.com/unders/mockingbird/server/pkg/testdata"

	"bytes"

	"github.com/unders/mockingbird/server/pkg/html"
)

//
// To update golden files:
//
//      go test github.com/unders/mockingbird/server/pkg/html -update
//

func TestTemplate_Execute_WhenSuccess_ReturnsHTML(t *testing.T) {
	tmpl, err := html.NewTemplate("testdata/tmpl/")
	testdata.AssertNil(t, err)

	data := map[string]string{
		"Title": "A Title",
		"Extra": "Extra fields are ignored",
	}

	tests := []struct {
		Name   string
		layout string
		file   string
		data   map[string]string
	}{
		{"block-is-overridden", "main1.html", "index1.html", data},
		{"block-is-not-overridden", "main1.html", "index2.html", data},
		{"render-partial-in-sub-block", "main1.html", "index3.html", data},
		{"renders-sub-partials", "main1.html", "a/a.html", data},
		{"renders-sub-sub-partials", "main1.html", "a/two/a.html", data},
		{"page-has-access-to-data-fields", "main2.html", "index4.html", data},
		{"part-has-access-to-data-fields", "main2.html", "index5.html", data},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			got, err := tmpl.Execute(tc.layout, tc.file, tc.data)
			testdata.AssertNil(t, err)

			want := testdata.ReadGolden(t, tc.Name)
			if !bytes.Equal(want, got) {
				t.Errorf("\nWant:\n%s\n Got:\n%s\n", want, got)
			}
		})
	}
}

func TestTemplate_Execute_WhenError_ReturnsError(t *testing.T) {
	tmpl, err := html.NewTemplate("testdata/tmpl/")
	testdata.AssertNil(t, err)

	data := map[string]string{
		"Title": "A cool title",
		"Extra": "Extra fields are ignored",
	}

	tests := []struct {
		layout string
		file   string
	}{
		{"not-found-layout.html", "index1.html"},
		{"main.html", "not-found.html"},

		// parsing errors
		{"missing-partial.html", "index.html"},
		{"missing-data-field.html", "index.html"},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			_, err := tmpl.Execute(tc.layout, tc.file, data)
			testdata.AssertErr(t, err)
		})
	}
}
func TestReloadableTemplate_Execute_WhenError_ReturnsError(t *testing.T) {
	tmpl, err := html.NewReloadableTemplate("testdata/tmpl/")
	testdata.AssertNil(t, err)

	data := map[string]string{
		"Title": "A cool title",
		"Extra": "Extra fields are ignored",
	}

	tests := []struct {
		layout       string
		file         string
		data         map[string]string
		wantNotFound bool
	}{
		// file access errors
		{"main.html", "file-not-found-error", data, true},
		{"layout-not-found.html", "index.html", data, true},

		// parsing errors
		{"missing-partial.html", "index.html", data, false},
		{"missing-data-field.html", "index.html", data, false},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			_, err := tmpl.Execute(tc.layout, tc.file, tc.data)
			if err == nil {
				t.Fatal("\nWant err\n Got nil\n")
			}

			got := errs.IsNotFound(err)
			if tc.wantNotFound != got {
				t.Errorf("\n Want: Not Found Error\n  Got: %t\nError: %s\n", got, err)
			}
		})
	}
}

func TestReloadableTemplate_Execute_WhenSuccess_ReturnsHTML(t *testing.T) {
	tmpl, err := html.NewReloadableTemplate("testdata/tmpl/")
	testdata.AssertNil(t, err)

	data := map[string]string{
		"Title": "A Title",
		"Extra": "Extra fields are ignored",
	}

	tests := []struct {
		Name   string
		layout string
		file   string
		data   map[string]string
	}{
		{"block-is-overridden", "main1.html", "index1.html", data},
		{"block-is-not-overridden", "main1.html", "index2.html", data},
		{"render-partial-in-sub-block", "main1.html", "index3.html", data},
		{"renders-sub-partials", "main1.html", "a/a.html", data},
		{"renders-sub-sub-partials", "main1.html", "a/two/a.html", data},
		{"page-has-access-to-data-fields", "main2.html", "index4.html", data},
		{"part-has-access-to-data-fields", "main2.html", "index5.html", data},
	}

	for _, tc := range tests {
		t.Run(tc.file, func(t *testing.T) {
			got, err := tmpl.Execute(tc.layout, tc.file, tc.data)
			testdata.AssertNil(t, err)

			if *testdata.Update {
				testdata.WriteGolden(t, tc.Name, got)
			}

			want := testdata.ReadGolden(t, tc.Name)
			if !bytes.Equal(want, got) {
				t.Errorf("\nWant:\n%s\n Got:\n%s\n", want, got)
			}
		})
	}
}
