package daemon

import (
	"context"

	"github.com/ProtocolONE/go-core/v2/pkg/entrypoint"
	"github.com/ProtocolONE/go-core/v2/pkg/logger"
	"github.com/qilin/crm-api/cmd"
	"github.com/qilin/crm-api/internal/daemon"
	"github.com/qilin/crm-api/pkg/http"
	"github.com/spf13/cobra"
)

const Prefix = "cmd.daemon"

var (
	Cmd = &cobra.Command{
		Use:           "daemon",
		Short:         "GraphQL API daemon",
		SilenceUsage:  true,
		SilenceErrors: true,
		Run: func(_ *cobra.Command, _ []string) {
			log := cmd.Slave.Logger().WithFields(logger.Fields{"service": Prefix})
			var (
				s *http.HTTP
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
				s, c, e = daemon.BuildHTTP(ctx, initial, cmd.Observer)
				return e
			}, func(ctx context.Context) error {
				return s.ListenAndServe(ctx)
			})
		},
	}
)

func init() {
	// pflags
	Cmd.PersistentFlags().StringP("http.bind", "b", ":8081", "bind address")
}
