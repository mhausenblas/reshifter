package etcd

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

var storetests = []struct {
	inpath string
	inval  string
	out    string
}{
	{"", "", ""},
	{"/", "root", ContentFile},
	{"/test", "some", "test/" + ContentFile},
	{"/test/this:also", "escaped", "test/thisESC_COLONalso/" + ContentFile},
}

func TestStore(t *testing.T) {
	cwd, _ := os.Getwd()
	for _, tt := range storetests {
		got, _ := store(".", tt.inpath, tt.inval)
		got, _ = filepath.Rel(cwd, got)
		want := tt.out
		if got != want {
			t.Errorf("etcd.store(\".\", %q, %q) => %q, want %q", tt.inpath, tt.inval, got, want)
		}
		// t.Logf("etcd.store(\".\", %q, %q) => %q", tt.inpath, tt.inval, got)
		// now clean up the directories and files created as a side effect:
		if got != "" {
			if got != "." {
				err := os.Remove(got)
				if err != nil {
					fmt.Printf("%s\n", err)
				}
			}
		}
	}
	_ = os.RemoveAll("test/")
}
