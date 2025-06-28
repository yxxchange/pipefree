package main

import (
	"github.com/yxxchange/pipefree/infra/dal"
	"gorm.io/gen"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./infra/dal/dao", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	g.ApplyBasic(dal.DbModels...)

	// Execute the generator
	g.Execute()
}
