module github.com/crawlab-team/crawlab-fs

go 1.15

replace (
	github.com/linxGnu/goseaweedfs => /Users/marvzhang/projects/tikazyq/goseaweedfs
	github.com/crawlab-team/crawlab-core => /Users/marvzhang/projects/crawlab-team/crawlab-core
)

require (
	github.com/linxGnu/goseaweedfs v0.1.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/crawlab-team/crawlab-core v0.0.1
)
