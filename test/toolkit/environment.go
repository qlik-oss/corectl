package toolkit

import (
	"flag"
)

var update = true // *flag.Bool("update", false, "update golden files")

var testAppName = "corectl_test_app.qvf"

var Engine1IP = flag.String("engineIP", "localhost:9076", "URL to first engine instance in docker-compose.yml i.e qix-engine-1")
var Engine2IP = flag.String("Engine2IP", "localhost:9176", "URL to second engine instance in docker-compose.yml i.e qix-engine-2")
var Engine3IP = flag.String("Engine3IP", "localhost:9276", "URL to third engine instance in docker-compose.yml i.e qix-engine-3")
