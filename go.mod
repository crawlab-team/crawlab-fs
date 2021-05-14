module github.com/crawlab-team/crawlab-fs

go 1.15

replace (
	github.com/crawlab-team/go-trace => /Users/marvzhang/projects/crawlab-team/go-trace
	github.com/linxGnu/goseaweedfs => /Users/marvzhang/projects/tikazyq/goseaweedfs
)

require (
	github.com/apex/log v1.9.0
	github.com/cenkalti/backoff/v4 v4.1.0
	github.com/crawlab-team/go-trace v0.0.0
	github.com/google/uuid v1.1.1
	github.com/linxGnu/goseaweedfs v0.1.5
	github.com/pkg/errors v0.9.1 // indirect
	github.com/stretchr/testify v1.6.1
)
