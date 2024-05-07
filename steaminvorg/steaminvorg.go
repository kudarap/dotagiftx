package steaminvorg

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/kudarap/dotagiftx/steam"
)

const (
	maxGetRetries = 10
	retrySleepDur = time.Second * 5
	// freshCacheDur = time.Hour
	freshCacheDur = time.Minute * 15
)

// InventoryAsset returns a compact format from all inventory data.
func InventoryAsset(steamID string) ([]steam.Asset, error) {
	inv, err := SWR(steamID, false)
	if err != nil {
		return nil, err
	}

	return inv.ToAssets(), nil
}

// SWR stale-while-re-invalidating crawled data.
func SWR(steamID string, strict bool) (*steam.AllInventory, error) {
	// check for freshly cached inventory
	m, err := GetMeta(steamID)
	if err != nil && strict {
		log.Println("STEAMINVORG META ERR", steamID, err)
		return nil, err
	}
	if m != nil && m.isCacheFresh() {
		log.Println("STEAMINVORG CACHED", steamID)
		return Get(steamID)
	}

	log.Println("STEAMINVORG CRAWL REQUEST", steamID)
	if _, err = Crawl(steamID); err != nil {
		log.Println("STEAMINVORG CRAWL REQUEST ERR", steamID, err)
		return nil, err
	}

	// check for meta until processed with a little bit of back-off
	for i := 1; i <= maxGetRetries; i++ {
		log.Println("STEAMINVORG TRY", steamID, i)
		m, err = GetMeta(steamID)
		if err != nil {
			log.Println("STEAMINVORG TRY ERR", steamID, err)
			return nil, err
		}
		if m != nil && m.Status == "success" {
			log.Println("STEAMINVORG TRY SUCCESS", steamID, i)
			break
		}

		if i == maxGetRetries {
			log.Println("STEAMINVORG TIMED OUT", steamID)
		}
		time.Sleep(retrySleepDur + time.Duration(i)*time.Second)
	}

	res, err := Get(steamID)
	if err != nil {
		log.Println("STEAMINVORG GET ERR", steamID, err)
		return nil, err
	}

	log.Println("STEAMINVORG GET DONE", steamID)
	return res, nil
}

// GetMeta https://db.steaminventory.org/SteamInventory/76561198264023028 - check queue state
func GetMeta(steamID string) (*Metadata, error) {
	url := fmt.Sprintf("https://db.steaminventory.org/SteamInventory/%s", steamID)
	var raw rawMetadata
	if err := getRequest(url, &raw); err != nil {
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

// Get https://data.steaminventory.org/SteamInventory/76561198264023028 - aggregated inventory
// https://data-gz.steaminventory.org/SteamInventory/76561198264023028 - aggregated inventory gzipped
func Get(steamID string) (*steam.AllInventory, error) {
	url := fmt.Sprintf("https://data-gz.steaminventory.org/SteamInventory/%s", steamID)
	all := &steam.AllInventory{}
	if err := getRequest(url, all); err != nil {
		return nil, err
	}

	return all, nil
}

// Crawl POST https://job.steaminventory.org/ScheduleInventoryCrawl?profile=76561198088587178
func Crawl(steamID string) (status string, err error) {
	url := fmt.Sprintf("https://job.steaminventory.org/ScheduleInventoryCrawl?profile=%s", steamID)
	res, err := http.Post(url, "", nil)
	if err != nil {
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
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

// Metadata represents response inventory metadata.
type Metadata struct {
	Status               string
	Count                int
	ScannedCount         int
	AllDescriptionLength int
	AllInventoryLength   int
	QueuedAt             time.Time
	CrawledAt            time.Time
	LastUpdated          time.Time
}

func (m Metadata) hasError() error {
	switch m.Status {
	case "error:403":
		return steam.ErrInventoryPrivate
	case "error:404":
		return fmt.Errorf("inventory not found")
	}

	return nil
}

func (m Metadata) isCacheFresh() bool {
	if m.LastUpdated.IsZero() {
		return false
	}
	return time.Now().Before(m.LastUpdated.Add(freshCacheDur))
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

	i := d.Items[0]
	m.Status = i.Status.S
	m.AllDescriptionLength, _ = strconv.Atoi(i.AllDescriptionLength.N)
	m.AllInventoryLength, _ = strconv.Atoi(i.AllInventoryLength.N)
	m.QueuedAt = parseStrTime(i.QueuedTimestamp.N)
	m.CrawledAt = parseStrTime(i.Timestamp.N)
	m.LastUpdated = parseStrTime(i.IndexTimestamp.N)
	m.Count = d.Count
	m.ScannedCount = d.ScannedCount
	return m
}

func getRequest(url string, data interface{}) error {
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return fmt.Errorf("not found")
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(body, data); err != nil {
		return err
	}

	return nil
}

func parseStrTime(ts string) time.Time {
	if ts == "" {
		return time.Time{}
	}

	sec, err := strconv.Atoi(ts)
	if err != nil {
		return time.Time{}
	}

	return time.Unix(int64(sec)/1000, 0)
}
