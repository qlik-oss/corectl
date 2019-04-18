package printer

import (
	"fmt"

	"github.com/qlik-oss/corectl/internal"
)

// PrintMetadata prints fields tables and associations to system out.
func PrintMetadata(data *internal.ModelMetadata) {

	fmt.Println("*** Fields ***")
	PrintFields(data, false)
	fmt.Println("\n*** Tables ***")
	PrintTables(data)
	fmt.Println("\n*** Associations ***")
	PrintAssociations(data)
}
