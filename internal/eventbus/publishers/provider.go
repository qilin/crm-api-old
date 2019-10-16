package publishers

import (
	"github.com/qilin/crm-api/internal/eventbus/common"
)

// ProviderPublishers
func ProviderPublishers() (common.Publishers, func(), error) {
	return []common.Publisher{
		//
	}, func() {}, nil
}
