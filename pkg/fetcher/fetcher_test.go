package fetcher

import (
	"reflect"
	"strings"
	"testing"
)

var page = `
<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8" />
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <title>Test Page</title>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" type="text/css" media="screen" href="main.css" />
  <script src="main.js"></script>
</head>
<body>
	<div>
		<p>Hello World</p>
			<a href="https://github.com">GitHub</a>
		</p>

		<a href="/about">About us</a>
	</div>
</body>
</html>
`

func TestFetcher(t *testing.T) {
	t.Run("Returns all links", func(t *testing.T) {
		f := New()
		r := strings.NewReader(page)
		base := "https://test.page.com"
		filter := ""

		got := f.FindLinks(r, base, filter)
		want := []string{"https://github.com", "https://test.page.com/about"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("Returns filtered links", func(t *testing.T) {
		f := New()
		r := strings.NewReader(page)
		base := "https://test.page.com"
		filter := "/"

		got := f.FindLinks(r, base, filter)
		want := []string{"https://test.page.com/about"}

		if !reflect.DeepEqual(got, want) {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
