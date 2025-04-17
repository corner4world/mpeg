module github.com/lkmio/mpeg

require github.com/lkmio/avformat v0.0.0

replace (
	github.com/lkmio/avformat => ../avformat
)

go 1.19
