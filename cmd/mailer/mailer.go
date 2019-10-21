package mailer

import (
	"context"

	"github.com/qilin/crm-api/internal/mailer"

	"github.com/ProtocolONE/go-core/v2/pkg/entrypoint"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/cmd"
	"github.com/spf13/cobra"
)

const Prefix = "cmd.mailer"

var (
	Cmd = &cobra.Command{
		Use:           "mailer",
		Short:         "Mailer daemon",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			log := cmd.Slave.Logger().WithFields(logger.Fields{"service": Prefix})
			var (
				//s *eventbus.EventBus
				c func()
				e error
			)
			defer func() {
				if c != nil {
					c()
				}
			}()
			cmd.Slave.Executor(func(ctx context.Context) error {
				cmd.Slave.OnReload(func(ctx context.Context) {
					initial, ok := entrypoint.CtxExtractInitial(ctx)
					log.Info("catch reload in %s, debug: %v, ok: %v", logger.Args(Prefix, initial.WorkDir, ok))
				})
				initial, _ := entrypoint.CtxExtractInitial(ctx)
				_, c, e = mailer.BuildMailer(ctx, initial, cmd.Observer)
				if e != nil {
					return e
				}
				return nil
			}, func(ctx context.Context) error {
				<-ctx.Done()
				return nil
			})
		},
	}
)
