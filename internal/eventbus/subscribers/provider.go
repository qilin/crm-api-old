package subscribers

import (
	"github.com/qilin/crm-api/internal/eventbus/common"
	"github.com/qilin/crm-api/internal/eventbus/subscribers/invites"
)

// ProviderSubscribers
func ProviderSubscribers(is *invites.InviteSubscriber) (common.Subscribers, func(), error) {
	return common.Subscribers{
		is,
	}, func() {}, nil
}
