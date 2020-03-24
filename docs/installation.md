### Use a release

**Scribe** is packaged as a stand-alone binary, just download a [release]() for your platform.

### Building from source

You will need the [Go tool chain](https://golang.org/) (**scribe** tested with v1.14) to build from source. Then just:

```sh
git clone https://github.com/will-rowe/scribe
cd scribe
go get -d -t -v ./...
go test ./...
go install
```
