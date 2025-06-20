package gorm_gen

import (
	"gorm.io/gen"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "../infra/dal", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	// Initialize a *gorm.DB instance
	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})

	// Use the above `*gorm.DB` instance to initialize the generator,
	// which is required to generate structs from db when using `GenerateModel/GenerateModelAs`
	// g.UseDB(db)

	// Generate default DAO interface for those specified structs
	g.ApplyBasic(model.Customer{}, model.CreditCard{}, model.Bank{}, model.Passport{})

	// Execute the generator
	g.Execute()
}
