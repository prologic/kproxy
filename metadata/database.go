package metadata

import (
	"encoding/json"
	"fmt"
	"github.com/prologic/bitcask"
	"kproxy/helpers"
	"net/http"
)

var _db *bitcask.Bitcask

func GetDatabaseSingleton() *bitcask.Bitcask {
	if _db == nil {
		panic("Database not yet initialised")
	}

	return _db
}

func Init() {
	db, err := bitcask.Open(helpers.GetDatabasePath())
	if err != nil {
		panic(err)
	}

	_db = db
}

type Resource struct {
	Name     string
	Expiry   int64 // UNIX timestamp in seconds
	MimeType string
	Visits   int

	Headers http.Header

	// for background downloads via the config API
	DownloadStatus string
}

func Get(name string) *Resource {
	resource := &Resource{
		Name: name,
	}

	db := GetDatabaseSingleton()
	value, err := db.Get([]byte(name))
	if err != nil {
		return resource
	}

	fmt.Println(string(value))
	_ = json.Unmarshal(value, &resource)
	fmt.Println(resource)
	return resource
}

func (resource *Resource) Save() {
	db := GetDatabaseSingleton()
	value, err := json.Marshal(resource)
	if err != nil {
		return
	}

	_ = db.Put([]byte(resource.Name), value)
}

func (resource *Resource) IncrementVisits() {
	resource.Visits++
	resource.Save()
}

func (resource *Resource) UpdateDownload(statusMessage string) {
	if resource.DownloadStatus == statusMessage {
		return
	}

	resource.DownloadStatus = statusMessage
	resource.Save()
}
