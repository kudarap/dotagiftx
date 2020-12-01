package verdeliv

/*

objective: verify buyer received reserved item from seller

params:
	- seller persona name: check for sender value
	- buyer steam id: for parsing inventory
	- item name: item key name to check against sender

result:
	- detect private inventory
	- detect malformed json
	- support multiple json for large inventory
	- challenge check

process:
	- download json inventory
	- parse json file
	- search item name
	- check sender name

*/

const (
	inventEndpoint  = "https://steamcommunity.com/profiles/%s/inventory/json/570/2?count=5000"
	inventEndpoint2 = "https://steamcommunity.com/inventory/%s/570/2?count=5000"
)
