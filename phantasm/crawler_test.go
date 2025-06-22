package phantasm

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	jsoniter "github.com/json-iterator/go"
)

var fastjson = jsoniter.ConfigFastest

func Test_merge_verify(t *testing.T) {
	invs, err := parseInventoryFiles("./testdata/inventory_raw.json", "./testdata/inventory_reduced.json")
	if err != nil {
		t.Fatal(err)
	}
	raw, reduced := invs[0], invs[1]

	merged := merge(raw)
	if diff := cmp.Diff(merged, reduced); diff != "" {
		t.Fatalf("merge mismatch (-want +got):\n%s", diff)
	}
	missing := checkMissingAssetDesc(raw)
	if len(missing) > 0 {
		t.Fatalf("raw missing asset desc: %d/%d", len(missing), raw.TotalInventoryCount)
	}
	missing = checkMissingAssetDesc(reduced)
	if len(missing) > 0 {
		t.Fatalf("reduced missing asset desc: %d/%d", len(missing), raw.TotalInventoryCount)
	}
}

func Test_merge_pagination(t *testing.T) {
	invs, err := parseInventoryFiles("./testdata/inventory_raw.json", "./testdata/inventory_reduced.json")
	if err != nil {
		t.Fatal(err)
	}
	reduced := invs[0]

	merged := merge(invs[1:]...)
	missing := checkMissingAssetDesc(merged)
	if len(missing) > 0 {
		t.Fatalf("raw missing asset desc: %d/%d", len(missing), merged.TotalInventoryCount)
	}
	missing = checkMissingAssetDesc(reduced)
	if len(missing) > 0 {
		t.Fatalf("reduced missing asset desc: %d/%d", len(missing), merged.TotalInventoryCount)
	}
}

func checkMissingAssetDesc(inv *inventory) []string {
	var missing []string
	for _, ass := range inv.Assets {
		var found bool
		for _, desc := range inv.Descriptions {
			if ass.Classid == desc.Classid {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, ass.Assetid)
			for _, desc := range inv.Descriptions {
				if ass.Classid == desc.Classid {
					fmt.Println(ass.Assetid, desc.Name)
					break
				}
			}
		}
	}
	return missing
}

func parseInventory(r io.Reader) (*inventory, error) {
	var raw inventory
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err = fastjson.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func parseInventoryFiles(paths ...string) ([]*inventory, error) {
	var invs []*inventory
	for _, path := range paths {
		f, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer f.Close()
		i, err := parseInventory(f)
		if err != nil {
			return nil, err
		}
		invs = append(invs, i)
	}
	return invs, nil
}
