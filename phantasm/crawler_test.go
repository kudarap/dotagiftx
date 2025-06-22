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
		t.Fatalf("raw missing Asset desc: %d/%d", len(missing), raw.TotalInventoryCount)
	}
	missing = checkMissingAssetDesc(reduced)
	if len(missing) > 0 {
		t.Fatalf("reduced missing Asset desc: %d/%d", len(missing), raw.TotalInventoryCount)
	}
}

func Test_merge_pagination(t *testing.T) {
	invs, err := parseInventoryFiles(
		"./testdata/inventory_paginated_reduced.json",
		"./testdata/inventory_paginated_1.json",
		"./testdata/inventory_paginated_2.json",
		"./testdata/inventory_paginated_3.json",
		"./testdata/inventory_paginated_4.json",
	)
	if err != nil {
		t.Fatal(err)
	}
	reduced := invs[0]

	var merged *Inventory
	for _, inv := range invs[1:] {
		merged = merge(merged, inv)
	}

	missing := checkMissingAssetDesc(merged)
	if len(missing) > 0 {
		t.Fatalf("merged missing Asset desc: %d/%d", len(missing), merged.TotalInventoryCount)
	}

	missing = checkMissingAssetDesc(reduced)
	if len(missing) > 0 {
		t.Fatalf("reduced missing Asset desc: %d/%d", len(missing), merged.TotalInventoryCount)
	}
}

func checkMissingAssetDesc(inv *Inventory) []string {
	var missing []string
	for _, ass := range inv.Assets {
		var found bool
		for _, desc := range inv.Descriptions {
			if ass.Classid == desc.Classid && ass.Instanceid == desc.Instanceid {
				found = true
				break
			}
		}
		if !found {
			missing = append(missing, ass.Assetid)
			for _, desc := range inv.Descriptions {
				if ass.Classid == desc.Classid {
					fmt.Println(ass.Assetid, desc.Name, ass.Classid, ass.Instanceid)
					break
				}
			}
		}
	}
	return missing
}

func parseInventory(r io.Reader) (*Inventory, error) {
	var raw Inventory
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	if err = fastjson.Unmarshal(b, &raw); err != nil {
		return nil, err
	}
	return &raw, nil
}

func parseInventoryFiles(paths ...string) ([]*Inventory, error) {
	var invs []*Inventory
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
