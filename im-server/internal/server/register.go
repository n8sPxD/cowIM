package server

import (
	"fmt"
	"time"

	"github.com/n8sPxD/cowIM/common/utils"
)

// 续租
func (s *Server) register() {
	c := s.svcCtx.Config
	if len(c.Etcd.Hosts) > 0 {
		ip, err := utils.GetLocalIP()
		if err != nil {
			panic(err)
		}
		addr := fmt.Sprintf("%s:%d", ip, c.Port)
		service := fmt.Sprintf("server%d", s.svcCtx.Config.WorkID)
		// 先注册一遍
		leaseID, err := s.svcCtx.RegisterHub.Register(s.ctx, service, addr, 0)
		if err != nil {
			panic(err)
		}
		go func() {
			// 异步持续续租
			for {
				s.registerHub.Register(s.ctx, service, addr, leaseID)
				time.Sleep(time.Duration(3)*time.Second - 100*time.Millisecond)
			}
		}()
	}
}
