package crawler

import "testing"

// TODO: test processURL.

func TestAbsURL(t *testing.T) {
	base := "http://test.com/foo/a"

	testCases := []struct {
		candidate string
		wantURL   string
		wantOK    bool
	}{
		{"/b", "http://test.com/b", true},
		{"c", "http://test.com/foo/c", true},
		{"../c", "http://test.com/c", true},
		{"http://test.com/bar", "http://test.com/bar", true},
		{"http://test.net/foo", "http://test.net/foo", false},
	}

	for _, tc := range testCases {
		gotURL, gotOK, err := absURL(base, tc.candidate)
		if err != nil {
			t.Errorf("got unexpected error: %v", err)
			continue
		}
		if gotURL != tc.wantURL {
			t.Errorf("url %q, want url %q, got %q", tc.candidate, tc.wantURL, gotURL)
		}
		if gotOK != tc.wantOK {
			t.Errorf("url %s, want ok %t, got %t", tc.candidate, tc.wantOK, gotOK)
		}
	}
}
