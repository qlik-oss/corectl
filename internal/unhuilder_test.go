package internal

import (
	"fmt"
	"testing"
)

func TestBuildName(t *testing.T) {
	fmt.Println(buildEntityFilename("wefwef", "mastesrobject", "table", "'='Halleluljah moment'"))
}
