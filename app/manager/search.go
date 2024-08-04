package manager

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"gin-demo/app/util"

	"github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

var (
	searchManager *SearchManager
)

func GetSearch() *SearchManager {
	return searchManager
}

type SearchManager struct {
	client  *elasticsearch.Client
	config  SearchConfig
	context context.Context
}

type SearchConfig struct {
	Url         string
	IndexPrefix string
}

func (manager *SearchManager) Setup(config SearchConfig) (err error) {

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return err
	}

	res, err := client.Info()
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	searchManager = &SearchManager{client: client, config: config, context: context.Background()}

	return nil
}

func (manager *SearchManager) Close() {

}

type Info struct {
	took       int
	timed_out  bool
	total      int
	successful int
	skipped    int
	failed     int
}

func (manager *SearchManager) CreateData(index string, key string, data interface{}) error {

	index = manager.config.IndexPrefix + index

	byteData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: key,
		Body:       bytes.NewReader(byteData),
	}

	res, doErr := req.Do(manager.context, manager.client)
	if doErr != nil {
		return doErr
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	var r map[string]interface{}
	if decoderErr := json.NewDecoder(res.Body).Decode(&r); decoderErr != nil {
		return decoderErr
	}

	return nil
}

func (manager *SearchManager) UpsertData(index string, id string, data interface{}) error {

	index = manager.config.IndexPrefix + index

	byteData, jsonErr := json.Marshal(data)
	if jsonErr != nil {
		return jsonErr
	}

	req := esapi.IndexRequest{
		Index:      index,
		DocumentID: id,
		Body:       bytes.NewReader(byteData),
	}

	res, doErr := req.Do(manager.context, manager.client)
	if doErr != nil {
		return doErr
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	var r map[string]interface{}
	if decoderErr := json.NewDecoder(res.Body).Decode(&r); decoderErr != nil {
		return decoderErr
	}

	return nil
}

func (manager *SearchManager) SearchData(index string, query map[string]interface{}, from int, size int) (int, []map[string]interface{}, error) {

	index = manager.config.IndexPrefix + index

	var buf bytes.Buffer
	if jsonErr := json.NewEncoder(&buf).Encode(query); jsonErr != nil {
		return 0, nil, jsonErr
	}

	res, searchErr := manager.client.Search(
		manager.client.Search.WithContext(manager.context),
		manager.client.Search.WithIndex(index),
		manager.client.Search.WithBody(&buf),
		manager.client.Search.WithFrom(from),
		manager.client.Search.WithSize(size),
		manager.client.Search.WithTrackTotalHits(true),
		manager.client.Search.WithPretty(),
	)
	if searchErr != nil {
		return 0, nil, searchErr
	}

	defer res.Body.Close()

	if res.IsError() {
		return 0, nil, errors.New(res.String())
	}

	var r map[string]interface{}
	if decoderErr := json.NewDecoder(res.Body).Decode(&r); decoderErr != nil {
		return 0, nil, decoderErr
	}

	repShards := r["_shards"].(map[string]interface{})
	shardtotal := int(repShards["total"].(float64))
	shardsuccessful := int(repShards["successful"].(float64))
	shardskipped := int(repShards["skipped"].(float64))
	shardfailed := int(repShards["failed"].(float64))

	timeout := r["timed_out"].(bool)
	took := int(r["took"].(float64))

	info := Info{
		took:       took,
		timed_out:  timeout,
		total:      shardtotal,
		successful: shardsuccessful,
		skipped:    shardskipped,
		failed:     shardfailed,
	}

	repHits := r["hits"].(map[string]interface{})
	total := repHits["total"].(map[string]interface{})
	hits := repHits["hits"].([]interface{})
	count := int(total["value"].(float64))

	var sourceList []map[string]interface{}
	for _, hit := range hits {
		obj := hit.(map[string]interface{})
		source := obj["_source"].(map[string]interface{})
		sourceList = append(sourceList, source)
	}

	fmt.Printf("SearchData, req:" + util.StructToJsonString(query))
	fmt.Printf("SearchData, res:" + util.StructToJsonString(sourceList))
	fmt.Printf("SearchData, count:" + strconv.Itoa(count))
	fmt.Printf("SearchData, info:" + fmt.Sprintf("%+v", info))

	return count, sourceList, nil
}

func (manager *SearchManager) BulkData(index string, data bytes.Buffer) error {

	index = manager.config.IndexPrefix + index

	req := esapi.BulkRequest{
		Index:   index,
		Refresh: "false",
		Body:    strings.NewReader(data.String() + "\n"),
	}

	res, doErr := req.Do(manager.context, manager.client)
	if doErr != nil {
		return doErr
	}
	defer res.Body.Close()

	if res.IsError() {
		return errors.New(res.String())
	}

	var r map[string]interface{}
	if decoderErr := json.NewDecoder(res.Body).Decode(&r); decoderErr != nil {
		return decoderErr
	}

	return nil
}

/*
{
   "hits" : {
      "total" :       14,
      "hits" : [
        {
          "_index":   "user",
          "_id":      "7",
          "_score":   1,
          "_source": {
             "date":    "2024-08-04",
             "name":    "PinYi",
          }
       }
      ],
      "max_score" :   1
   },
   "took" :           4,
   "_shards" : {
      "failed" :      0,
      "successful" :  10,
      "total" :       10
   },
   "timed_out" :      false
}
*/
