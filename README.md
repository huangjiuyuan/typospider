# Typospider

---

## Introduction

Typospider is a tool for crawling a GitHub project and finding out typos in the comments.

## How to Start

Typospider relies on [Elasticsearch](https://www.elastic.co/) 6.0 or above and [LanguageTool](https://languagetool.org/) 4.0 to work properly. Elasticsearch is a distributed, RESTful search and analytics engine capable of solving a growing number of use cases. As the heart of the Elastic Stack, it centrally stores your data so you can discover the expected and uncover the unexpected. LanguageTool is a proof­reading service for English, German, Polish, Russian, and more than 20 other languages. All files and typos can be searched and viewed with [Kibana](https://www.elastic.co/products/kibana), which lets you visualize your Elasticsearch data.

To start an Elasticsearch instance, Please download Elasticsearch from [here](https://www.elastic.co/downloads/elasticsearch). Elasticsearch will be serving on port 9200.

To start an Kibana instance, Please download Kibana from [here](https://www.elastic.co/downloads/kibana). Kibana will be serving on port 5601.

LanguageTool can be served locally on your machine. To start an instance of LanguageTool, please run:

```
$ curl -L https://git.io/vNqdP | bash -
$ cd LanguageTool-*
$ java -cp languagetool-server.jar org.languagetool.server.HTTPServer --port 6066
```

The LanguageTool will serve on port 6066 now.

Run `go run cmd/main.go` to start processing GitHub project. Visit Kibana on `localhost:5601` to check the result.
