package main

import "flag"

const (
	defaultPort          = 8080
	defaultBasePath      = "tmp/"
	defaultMaxRecordSize = 64 * 1024
)

var (
	httpPort      int
	basePath      string
	maxRecordSize int
)

func init() {
	flag.IntVar(&httpPort, "port", defaultPort, "http server listening port")
	flag.StringVar(&basePath, "path", defaultBasePath, "storage path")
	flag.IntVar(&maxRecordSize, "maxRecordSize", defaultMaxRecordSize, "max size of a database record")
}
func main() {
	flag.Parse()
}
