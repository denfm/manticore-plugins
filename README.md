# Manticore search engine plugins

## Official example and docs
1. [Documentation](https://manual.manticoresearch.com/Extensions/UDFs_and_Plugins/UDFs_and_Plugins)
1. [Example go](https://github.com/manticoresoftware/plugins/tree/master/curl)

## BUILD 
Requires [Go](https://golang.org/doc/install). Tested with Go 1.15.

Clone this repo locally and run test, build:
```
mkdir -p $HOME/manticotre-plugins && \
cd $HOME/manticotre-plugins && \
git clone https://github.com/denfm/manticore-plugins ./ && \
git clone https://github.com/manticoresoftware/manticoresearch ./manticore && \
cd plugins/group-sort && make test && make build && \
cd ../../build && ls -la
```

LICENSE
========

See [LICENSE](./LICENSE)