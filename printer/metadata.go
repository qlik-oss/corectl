package printer

import (
	"fmt"
	"github.com/qlik-oss/corectl/internal"
)

// PrintMetadata prints fields tables and associations to system out.
func PrintMetadata(data *internal.ModelMetadata, mode PrintMode) {
	fmt.Println("*** Fields ***")
	PrintFields(data, false, mode)
	fmt.Println("\n*** Tables ***")
	PrintTables(data)
	fmt.Println("\n*** Associations ***")
	PrintAssociations(data)
}
