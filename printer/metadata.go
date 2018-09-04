package printer

import "github.com/qlik-oss/core-corectl/internal"

// TODO session app!

func PrintMetadata(data *internal.ModelMetadata) {

	PrintFields(data, false)
	PrintTables(data)
	PrintAssociations(data)

}
