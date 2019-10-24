package subscribers

import (
	"github.com/qilin/crm-api/internal/eventbus/common"
)

// ProviderSubscribers
func ProviderSubscribers() (common.Subscribers, func(), error) {
	return common.Subscribers{
		NewInviteSubscriber(),
	}, func() {}, nil
}
