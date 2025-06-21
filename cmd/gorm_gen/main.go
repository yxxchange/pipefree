package main

import (
	"github.com/yxxchange/pipefree/infra/dal/model"
	"gorm.io/gen"
)

func main() {
	// Initialize the generator with configuration
	g := gen.NewGenerator(gen.Config{
		OutPath:       "./infra/dal/dao", // output directory, default value is ./query
		Mode:          gen.WithDefaultQuery | gen.WithQueryInterface,
		FieldNullable: true,
	})

	ApplyList := []interface{}{
		model.PipeCfg{},
		model.PipeExec{},
		model.PipeVersion{},
		model.NodeCfg{},
		model.NodeExec{},
	}
	g.ApplyBasic(ApplyList...)

	// Execute the generator
	g.Execute()
}
