// +build wireinject

package postgres

import (
	"context"

	"github.com/google/wire"
	"github.com/jinzhu/gorm"
	"github.com/qilin/go-core/config"
	"github.com/qilin/go-core/invoker"
	"github.com/qilin/go-core/provider"
)

// Build returns GORM instance with resolved dependencies
func Build(ctx context.Context, initial config.Initial, observer invoker.Observer) (*gorm.DB, func(), error) {
	panic(wire.Build(WireSet, provider.Set))
}

// BuildTest returns stub/mock instance GORM with resolved dependencies
func BuildTest(ctx context.Context, initial config.Initial, observer invoker.Observer) (*gorm.DB, func(), error) {
	panic(wire.Build(WireTestSet))
}
