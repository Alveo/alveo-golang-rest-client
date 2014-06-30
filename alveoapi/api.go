// Package hcsvlabapi provides an convienience implementation of the client side of the
// HCS-Vlab API, which can be seen at http://wherever-the-api-is.com
package alveoapi

import (
  "fmt"
  "net/http"
  "encoding/json"
  "log"
  "errors"
  "strconv"
  "io/ioutil"
  "time"
)


var logger *log.Logger


// A representation of an item list from the HCSvLab API
type ItemList struct {
  Name string
  Num_items float64
  Items []string
}

// The response from API version
type ApiVersion struct {
 Api_version string `json:"API version"`
}

// A representation of a document within an item from the HCSvLab API
type DocIdentifier struct {
  Size string
  Url string
  Type string
}

// A representation of all the itemlists a user has access to:
type ItemLists struct {
  Own []ItemListIdentifier `json:"own"`
  Shared []ItemListIdentifier `json:"shared"`
}

// A representation of an itemlist from the /item_lists response
type ItemListIdentifier struct {
  Name string `json:"name"`
  ItemListUrl string `json:"item_list_url"`
  NumItems int64 `json:"num_items"`
  Shared bool `json:"shared"`
}

// A representation of the annotations associated with an item from the HCSvLab API
type AnnotationList struct {
  CommonProperties AnnotationProperties `json:"commonProperties"`
  Annotations []Annotation `json:"alveo:annotations"`
}

type AnnotationProperties struct {
  Annotates string `json:"alveo:annotates"`
}

// An annotation associated with a documents. AnnotationLists have more than one of these.
type Annotation struct {
  Type string `json:"type"`
  Label string `json:"label"`
  Start string `json:"start"`
  End string `json:"end"`
}

// An item that contains metadata about a document from the HCSvLab API
type Item struct {
 Catalog_url string `json:"alveo:catalog_url"`
 Metadata map[string]string `json:"alveo:metadata"`
 Primary_text_url string `json:"alveo:primary_text_url"`
 Annotations_url string `json:"alveo:annotations_url"`
 Documents []DocIdentifier `json:"alveo:documents"`
}

type Api struct {
  Base string
  Key string
}

// Returns the API version provided on the server
func (api *Api) GetVersion() (ver ApiVersion, err error) {
  url := fmt.Sprintf("%s/version.json",api.Base)
  data, err := api.Get(url)
  if err != nil {
    return
  }

  err = json.Unmarshal(data,&ver)
  return
}

// Sets the logger to be used when making requests. Useful for debugging
func SetLogger(newlogger *log.Logger) {
  logger = newlogger
}

// Helper function that gets the raw data from the URL specified,
// by providing the API key appropriately.
func (api *Api) Get(reqUrl string) (data []byte, err error) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", reqUrl, nil)
  req.Header.Add("X-API-KEY",api.Key)
  if logger != nil {
    logger.Println("Requesting ",reqUrl,"with Key",api.Key)
  }
  start := time.Now()
  resp, err := client.Do(req)
  if err != nil {
    return
  }
  if resp.StatusCode != 200 {
    err = errors.New("Status " + strconv.Itoa(resp.StatusCode) + " from " + reqUrl)
    return
  }
  data, err = ioutil.ReadAll(resp.Body)
  end := time.Now()
  if logger != nil {
    logger.Println("Time",reqUrl,end.Sub(start).Seconds(),resp.ContentLength)
  }
  resp.Body.Close()
  return
}

// Function to return an ItemList corresponding to the numbered itemlist
// given
func (api *Api) GetItemList(list int) (il ItemList, err error)  {
  url := fmt.Sprintf("%s/item_lists/%d.json",api.Base, list)
  data, err := api.Get(url);
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&il)
  return
}

// Function to enumerate the ItemLists a user has access to
func (api *Api) GetItemLists() (il ItemLists, err error)  {
  url := fmt.Sprintf("%s/item_lists.json",api.Base)
  data, err := api.Get(url);
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&il)
  return
}

// Function to return the annotations associated with a particular item
func (api *Api) GetAnnotations(item Item) (al AnnotationList, err error)  {
  data, err := api.Get(item.Annotations_url)
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&al)
  return
}


// Function to get a particular item from the item's url
func (api *Api) GetItemFromUri(url string)  (item Item, err error) {
  data, err := api.Get(url)
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&item);
  return
}
