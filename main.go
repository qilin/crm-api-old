package main

import (
	"github.com/qilin/crm-api/cmd/daemon"
	"github.com/qilin/crm-api/cmd/mailer"
	"github.com/qilin/crm-api/cmd/migrate"
	"github.com/qilin/crm-api/cmd/root"
	"github.com/qilin/crm-api/cmd/sdk"
	"github.com/qilin/crm-api/cmd/version"
)

func main() {
	root.Execute(version.Cmd, migrate.Cmd, daemon.Cmd, mailer.Cmd, sdk.Cmd)
}
