
// Package hcsvlabapi provides an convienience implementation of the client side of the
// HCS-Vlab API, which can be seen at http://wherever-the-api-is.com
package hcsvlabapi

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



// A representation of an item list from the HCSvLab API
type ItemList struct {
  Name string
  Num_items float64
  Items []string
}

type ApiVersion struct {
 Api_version string `json:"API version"`
}

// A representation of a document within an item from the HCSvLab API
type DocIdentifier struct {
  Size string
  Url string
  Type string
}

// A representation of the annotations associated with an item from the HCSvLab API
type AnnotationList struct {
  Annotates string
  Annotations_found float64
  Annotations []Annotation
}

// An annotation associated with a documents. AnnotationLists have more than one of these.
type Annotation struct {
  Type string
  Label string
  Start float64
  End float64
}

// An item that contains metadata about a document from the HCSvLab API
type Item struct {
 Catalog_url string
 Metadata map[string]string
 Primary_text_url string
 Annotations_url string
 Documents []DocIdentifier
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

  log.Println("Json:",string(data))
  err = json.Unmarshal([]byte(data),&ver)
  log.Println("Unmarshalled:",ver);
  return
}

// Helper function that gets the raw data from the URL specified,
// by providing the API key appropriately.
func (api *Api) Get(url string) (data []byte, err error) {
  client := &http.Client{}
  req, err := http.NewRequest("GET", url, nil)
  req.Header.Add("X-API-KEY",api.Key)
  log.Println("Requesting ",url,"with Key",api.Key)
  start := time.Now()
  resp, err := client.Do(req)
  if err != nil {
    return
  }
  if resp.StatusCode != 200 {
    err = errors.New("Status " + strconv.Itoa(resp.StatusCode) + " from " + url)
    return
  }
  data, err = ioutil.ReadAll(resp.Body)
  end := time.Now()
  log.Println("Time",url,end.Sub(start).Seconds(),resp.ContentLength)
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


// Function to return the annotations associated with a particular item
func (api *Api) GetAnnotations(item Item) (al AnnotationList, err error)  {
  data, err := api.Get(item.Annotations_url)
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&al)
  return
}


// Function to get a particular item from 
func (api *Api) GetItemFromUri(url string)  (item Item, err error) {
  data, err := api.Get(url)
  if err != nil {
    return
  }
  err = json.Unmarshal(data,&item);
  return
}
