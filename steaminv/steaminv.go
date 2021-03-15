package steaminv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/kudarap/dotagiftx/steam"
)

const (
	maxGetRetries = 5
	retrySleepDur = time.Second * 5
	freshCacheDur = time.Hour * 24
)

// InventoryAsset returns a compact format from all inventory data.
func InventoryAsset(steamID string) ([]steam.Asset, error) {
	inv, err := SWR(steamID)
	if err != nil {
		return nil, err
	}

	return inv.ToAssets(), nil
}

// SWR stale-while-re-invalidating crawled data.
func SWR(steamID string) (*steam.AllInventory, error) {
	// check for freshly cached inventory
	//log.Println(steamID, "checking for fresh cache...")
	m, err := GetMeta(steamID)
	if err != nil {
		return nil, err
	}
	if m != nil && m.isCacheFresh() {
		//log.Println(steamID, "cache is still fresh!", m.LastUpdated)
		//defer func() {
		//	if _, err = Crawl(steamID); err != nil {
		//		log.Println("error invalidating", err)
		//	}
		//}()
		return Get(steamID)
	}

	// crawl request
	//log.Println(steamID, "sending crawl request...")
	if _, err = Crawl(steamID); err != nil {
		return nil, err
	}

	// check for meta until processed with 5 reties
	for i := 1; i <= maxGetRetries; i++ {
		//log.Println(steamID, "checking metadata. retry", i, "...")
		m, err = GetMeta(steamID)
		if err != nil {
			return nil, err
		}
		if m != nil && m.Status == "success" {
			break
		}
		time.Sleep(retrySleepDur)
	}

	// get inventory
	//log.Println(steamID, "getting inventory...")
	return Get(steamID)
}

// https://data.steaminventory.org/SteamInventory/76561198264023028 - aggregated inventory
// https://data-gz.steaminventory.org/SteamInventory/76561198264023028 - aggregated inventory gzipped
func Get(steamID string) (*steam.AllInventory, error) {
	url := fmt.Sprintf("https://data-gz.steaminventory.org/SteamInventory/%s", steamID)
	all := &steam.AllInventory{}
	if err := getRequest(url, all); err != nil {
		return nil, err
	}

	return all, nil
}

// POST https://job.steaminventory.org/ScheduleInventoryCrawl?profile=76561198088587178
func Crawl(steamID string) (status string, err error) {
	url := fmt.Sprintf("https://job.steaminventory.org/ScheduleInventoryCrawl?profile=%s", steamID)
	res, err := http.Post(url, "", nil)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}

	crawlRes := struct {
		Status string `json:"status"`
	}{}
	if err = json.Unmarshal(body, &crawlRes); err != nil {
		return
	}

	return crawlRes.Status, nil
}

// Metadata represents reposnse inventory metadata.
type Metadata struct {
	Status               string
	Count                int
	ScannedCount         int
	AllDescriptionLength int
	AllInventoryLength   int
	QueuedAt             *time.Time
	CrawledAt            *time.Time
	LastUpdated          *time.Time
}

func (d Metadata) hasError() error {
	switch d.Status {
	case "error:403":
		return steam.ErrInventoryPrivate
	case "error:404":
		return fmt.Errorf("inventory not found")
	}

	return nil
}

func (d Metadata) isCacheFresh() bool {
	return time.Now().Before(d.LastUpdated.Add(freshCacheDur))
}

type rawMetadata struct {
	Count int `json:"Count"`
	Items []struct {
		AllDescriptionLength struct {
			N string `json:"N"`
		} `json:"all_description_length"`
		AllInventoryLength struct {
			N string `json:"N"`
		} `json:"all_inventory_length"`
		QueuedTimestamp struct {
			N string `json:"N"`
		} `json:"queued_timestamp"`
		IndexTimestamp struct {
			N string `json:"N"`
		} `json:"index_timestamp"`
		Profile struct {
			S string `json:"S"`
		} `json:"profile"`
		Status struct {
			S string `json:"S"`
		} `json:"status"`
		Timestamp struct {
			N string `json:"N"`
		} `json:"timestamp"`
		ResultURL struct {
			S string `json:"S"`
		} `json:"result_url"`
		ID struct {
			NULL bool `json:"NULL"`
		} `json:"id"`
		InventoryFormat struct {
			N string `json:"N"`
		} `json:"inventory_format"`
	} `json:"Items"`
	ScannedCount int `json:"ScannedCount"`
}

func (d *rawMetadata) format() *Metadata {
	m := &Metadata{}
	if len(d.Items) == 0 {
		return nil
	}

	mdi := d.Items[0]
	m.Status = mdi.Status.S
	m.AllDescriptionLength, _ = strconv.Atoi(mdi.AllDescriptionLength.N)
	m.AllInventoryLength, _ = strconv.Atoi(mdi.AllInventoryLength.N)
	m.QueuedAt = parseStrTime(mdi.QueuedTimestamp.N)
	m.CrawledAt = parseStrTime(mdi.Timestamp.N)
	m.LastUpdated = parseStrTime(mdi.IndexTimestamp.N)
	m.Count = d.Count
	m.ScannedCount = d.ScannedCount
	return m
}

// https://db.steaminventory.org/SteamInventory/76561198264023028 - check queue state
func GetMeta(steamID string) (*Metadata, error) {
	url := fmt.Sprintf("https://db.steaminventory.org/SteamInventory/%s", steamID)
	raw := &rawMetadata{}
	if err := getRequest(url, raw); err != nil {
		return nil, err
	}

	m := raw.format()
	if m == nil {
		return nil, nil
	}
	if err := m.hasError(); err != nil {
		return nil, err
	}
	return m, nil
}

func getRequest(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}

func parseStrTime(ts string) *time.Time {
	sec, _ := strconv.ParseInt(ts, 10, 64)
	t := time.Unix(sec/1000, 0)
	return &t
}
