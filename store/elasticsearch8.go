//go:build !windows

package store

import (
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/whj1990/go-core/handler"
)

func NewElasticsearch8() *elasticsearch.Client {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"https://localhost:9200",
			"https://localhost:9201",
		},
		// ...
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		handler.HandleError(err)
	}
	return es
}
