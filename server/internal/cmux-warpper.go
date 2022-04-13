package internal

import (
	"sync"

	"github.com/soheilhy/cmux"
	netPkg "github.com/thataway/common-lib/pkg/net"
)

// NewCMux anew CMux and wrap it
func NewCMux(endpoint *netPkg.Endpoint) (cmux.CMux, error) {
	l, e := netPkg.Listen(endpoint)
	if e != nil {
		return nil, e
	}
	var once sync.Once
	return &cMuxWrapper{
		CMux: cmux.New(NoCloseListener{Listener: l}),
		closeNativeListener: func() {
			once.Do(func() {
				_ = l.Close()
			})
		},
	}, nil
}

var (
	_           = NewCMux
	_ cmux.CMux = (*cMuxWrapper)(nil)
)

type cMuxWrapper struct {
	cmux.CMux
	closeNativeListener func()
}

func (cm *cMuxWrapper) Close() {
	cm.closeNativeListener()
	cm.CMux.Close()
}
