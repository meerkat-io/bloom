package strcase

import "testing"

var cases = []struct {
	in, snake, chain, camel, pascal string
}{
	{"", "", "", "", ""},
	{"hello", "hello", "hello", "hello", "Hello"},
	{"Hello World", "hello_world", "hello-world", "helloWorld", "HelloWorld"},
	{"helloWorld", "hello_world", "hello-world", "helloWorld", "HelloWorld"},
	{"HelloWorld", "hello_world", "hello-world", "helloWorld", "HelloWorld"},
	{"hello_world", "hello_world", "hello-world", "helloWorld", "HelloWorld"},
	{"hello-world", "hello_world", "hello-world", "helloWorld", "HelloWorld"},
	{"  multiple   spaces ", "multiple_spaces", "multiple-spaces", "multipleSpaces", "MultipleSpaces"},
	{"PDFLoader", "pdf_loader", "pdf-loader", "pdfLoader", "PdfLoader"},
	{"userID42", "user_id_42", "user-id-42", "userId42", "UserId42"},
	{"__a__b--c", "a_b_c", "a-b-c", "aBC", "ABC"},
	{"éclair Über", "éclair_über", "éclair-über", "éclairÜber", "ÉclairÜber"},
}

func TestCases(t *testing.T) {
	fns := []struct {
		name string
		fn   func(string) string
		get  func(i int) string
	}{
		{"SnakeCase", SnakeCase, func(i int) string { return cases[i].snake }},
		{"ChainCase", ChainCase, func(i int) string { return cases[i].chain }},
		{"CamelCase", CamelCase, func(i int) string { return cases[i].camel }},
		{"PascalCase", PascalCase, func(i int) string { return cases[i].pascal }},
	}
	for _, f := range fns {
		for i, c := range cases {
			if got, want := f.fn(c.in), f.get(i); got != want {
				t.Errorf("%s(%q) = %q, want %q", f.name, c.in, got, want)
			}
		}
	}
}

func TestInvalidUTF8(t *testing.T) {
	in := "bad\xffInput"
	// Invalid UTF-8 is not split on case boundaries (ToLower maps the bad byte to U+FFFD).
	if got := SnakeCase(in); got != "bad\uFFFDinput" {
		t.Errorf("SnakeCase(%q) = %q", in, got)
	}
}
