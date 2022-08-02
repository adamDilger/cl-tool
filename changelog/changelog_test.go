package changelog

import (
	"bytes"
	"testing"
)

func TestTitle(t *testing.T) {
	cv := ChangelogVersion{Version: "Unreleased"}

	if cv.Title() != "## [Unreleased]\n" {
		t.Error("Unreleased version title not working")
	}

	cv = ChangelogVersion{Version: "1.2.3", Date: "2006-01-02"}
	if cv.Title() != "## [1.2.3] - 2006-01-02\n" {
		t.Error("Version title not working")
	}
}

func TestIsUnreleased(t *testing.T) {
	cv := ChangelogVersion{Version: "Unreleased"}

	if !cv.IsUnreleased() {
		t.Error("IsUnreleased should be true")
	}

	cv = ChangelogVersion{Version: "Date is Empty"}
	if !cv.IsUnreleased() {
		t.Error("IsUnreleased should be true")
	}
}

func TestEntryListRender(t *testing.T) {
	var el EntryList
	el = append(el, "hello")
	el = append(el, "world")

	var buf bytes.Buffer
	el.Render(&buf, "HelloWorld")

	got := buf.String()
	wanted := "### HelloWorld\n- hello\n- world\n"
	if got != wanted {
		t.Errorf("EntryList render failed, got %s wanted %s\n", got, wanted)
	}
}

func TestChangelogRenderHeadAndTail(t *testing.T) {
	cl := Changelog{
		head: "### head\n",
		tail: "### tail\n",
		Versions: []ChangelogVersion{{
			Version: "1.1.1", Date: "2001-01-01",
			Files: []ChanglogFile{{Added: []string{"hello"}}},
		}},
	}

	var buf bytes.Buffer
	cl.Render(&buf)

	wanted := `### head

## [1.1.1] - 2001-01-01
### Added
- hello

### tail
`

	got := buf.String()
	if wanted != got {
		t.Errorf("EntryList render failed, got %s wanted %s\n", got, wanted)
	}
}

func TestChangelogRender(t *testing.T) {
	cf1 := ChanglogFile{Added: []string{"Hello", "World"}}
	cf2 := ChanglogFile{Added: []string{"added"}, Changed: []string{"Changed"}}
	cf3 := ChanglogFile{Removed: []string{"ok"}}

	cl := Changelog{Versions: []ChangelogVersion{
		{Version: "1.1.1", Date: "2001-01-01", Files: []ChanglogFile{cf2}},            // last
		{Version: "2.2.2", Date: "2002-02-02", Files: []ChanglogFile{cf1}},            // first
		{Version: "Unreleased", Files: []ChanglogFile{{Deprecated: []string{"dep"}}}}, // first
		{Version: "1.1.2", Date: "2001-01-02", Files: []ChanglogFile{cf3}},            // middle
	}}

	var buf bytes.Buffer
	cl.Render(&buf)

	wanted := `
## [Unreleased]
### Deprecated
- dep

## [2.2.2] - 2002-02-02
### Added
- Hello
- World

## [1.1.2] - 2001-01-02
### Removed
- ok

## [1.1.1] - 2001-01-01
### Added
- added
### Changed
- Changed
`

	got := buf.String()
	if wanted != got {
		t.Errorf("EntryList render failed, got %s wanted %s\n", got, wanted)
	}
}
