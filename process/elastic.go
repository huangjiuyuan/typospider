package process

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const mapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"file":{
			"properties": {
				"path": {
					"type": "string"
				},
				"size": {
					"type": "integer"
				},
				"sha": {
					"type": "string"
				},
				"url": {
					"type": "string"
				},
				"data": {
					"type": "text"
				},
				"typos": {
					"type": "nested",
					"properties": {}
				},
				"valid": {
					"type": "boolean"
				}
			}
		}
	}
}`

type Elastic struct {
	Endpoint string
	Version  string

	mapping string
	ctx     context.Context
	client  *elastic.Client
}

func InitClient(scheme string, host string, port string) (*Elastic, error) {
	ep := scheme + "://" + host + ":" + port
	ctx := context.Background()

	client, err := elastic.NewClient(elastic.SetURL(ep))
	if err != nil {
		return nil, err
	}

	version, err := client.ElasticsearchVersion(ep)
	if err != nil {
		return nil, err
	}

	return &Elastic{
		Endpoint: ep,
		Version:  version,

		mapping: mapping,
		ctx:     ctx,
		client:  client,
	}, nil
}

func (es *Elastic) CreateIndex(index string) error {
	exists, err := es.client.IndexExists(index).Do(es.ctx)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index %s already exists", index)
	}

	result, err := es.client.CreateIndex(index).BodyString(es.mapping).Do(es.ctx)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return fmt.Errorf("index %s creation not acknowledged", index)
	}

	return nil
}

func (es *Elastic) DeleteIndex(index string) error {
	exists, err := es.client.IndexExists(index).Do(es.ctx)
	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("index %s does not exist", index)
	}

	result, err := es.client.DeleteIndex(index).Do(es.ctx)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return fmt.Errorf("index %s deletion not acknowledged", index)
	}

	return nil
}

func (es *Elastic) IndexFile(index string, file File) (*elastic.IndexResponse, error) {
	resp, err := es.client.Index().
		Index(index).
		Type("file").
		Id(file.SHA).
		BodyJson(file).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) GetFile(index string, id string) (*elastic.GetResult, error) {
	resp, err := es.client.Get().
		Index(index).
		Type("file").
		Id(id).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) UpdateFile(index string, id string, file File) (*elastic.UpdateResponse, error) {
	resp, err := es.client.Update().
		Index(index).
		Type("file").
		Id(id).
		Upsert(file).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) DeleteFile(index string, id string, file File) (*elastic.DeleteResponse, error) {
	resp, err := es.client.Delete().
		Index(index).
		Type("file").
		Id(id).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
