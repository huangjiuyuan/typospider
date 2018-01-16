package process

import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
)

const fileMapping = `
{
    "mappings":{
        "file":{
            "properties":{
                "path":{
                    "type":"keyword"
                },
                "size":{
                    "type":"integer"
                },
                "sha":{
                    "type":"keyword"
                },
                "url":{
                    "type":"text"
                },
                "data":{
                    "type":"text"
                },
                "typos":{
                    "type":"text"
                },
                "valid":{
                    "type":"boolean"
                }
            }
        }
    }
}`

const typoMapping = `
{
    "mappings":{
        "typo":{
            "properties":{
                "sha":{
                    "type":"keyword"
                },
                "match":{
                    "type":"object",
                    "properties":{
                        "message":{
                            "type":"text"
                        },
                        "shortMessage":{
                            "type":"text"
                        },
                        "offset":{
                            "type":"integer"
                        },
                        "length":{
                            "type":"integer"
                        },
                        "replacements":{
                            "type":"nested",
                            "properties":{
                                "value":{
                                    "type":"text"
                                }
                            }
                        },
                        "context":{
                            "type":"object",
                            "properties":{
                                "text":{
                                    "type":"text"
                                },
                                "offset":{
                                    "type":"integer"
                                },
                                "length":{
                                    "type":"integer"
                                }
                            }
                        },
                        "sentence":{
                            "type":"text"
                        },
                        "rule":{
                            "type":"object",
                            "properties":{
                                "id":{
                                    "type":"text"
                                },
                                "subId":{
                                    "type":"text"
                                },
                                "description":{
                                    "type":"text"
                                },
                                "urls":{
                                    "type":"nested",
                                    "properties":{
                                        "value":{
                                            "type":"text"
                                        }
                                    }
                                },
                                "issueType":{
                                    "type":"text"
                                },
                                "category":{
                                    "type":"object",
                                    "properties":{
                                        "id":{
                                            "type":"text"
                                        },
                                        "name":{
                                            "type":"text"
                                        }
                                    }
                                }
                            }
                        }
                    }
                },
                "valid":{
                    "type":"boolean"
                }
            }
        }
    }
}`

type Elastic struct {
	Endpoint string
	Version  string

	fileMapping string
	typoMapping string
	ctx         context.Context
	client      *elastic.Client
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

		fileMapping: fileMapping,
		typoMapping: typoMapping,
		ctx:         ctx,
		client:      client,
	}, nil
}

func (es *Elastic) CreateFileIndex(index string) error {
	exists, err := es.client.IndexExists(index).Do(es.ctx)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index %s already exists", index)
	}

	result, err := es.client.CreateIndex(index).BodyString(es.fileMapping).Do(es.ctx)
	if err != nil {
		return err
	}
	if !result.Acknowledged {
		return fmt.Errorf("index %s creation not acknowledged", index)
	}

	return nil
}

func (es *Elastic) CreateTypoIndex(index string) error {
	exists, err := es.client.IndexExists(index).Do(es.ctx)
	if err != nil {
		return err
	}

	if exists {
		return fmt.Errorf("index %s already exists", index)
	}

	result, err := es.client.CreateIndex(index).BodyString(es.typoMapping).Do(es.ctx)
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

func (es *Elastic) IndexTypo(index string, typo Typo) (*elastic.IndexResponse, error) {
	resp, err := es.client.Index().
		Index(index).
		Type("typo").
		Id(typo.SHA).
		BodyJson(typo).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) GetTypo(index string, id string) (*elastic.GetResult, error) {
	resp, err := es.client.Get().
		Index(index).
		Type("typo").
		Id(id).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) UpdateTypo(index string, id string, typo Typo) (*elastic.UpdateResponse, error) {
	resp, err := es.client.Update().
		Index(index).
		Type("typo").
		Id(id).
		Upsert(typo).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (es *Elastic) DeleteTypo(index string, id string, typo Typo) (*elastic.DeleteResponse, error) {
	resp, err := es.client.Delete().
		Index(index).
		Type("typo").
		Id(id).
		Do(es.ctx)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
