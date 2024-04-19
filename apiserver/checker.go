package apiserver

import (
	"context"
	"github.com/rs/zerolog/log"
	"time"
)

type Checker interface {
	Check(ctx context.Context)
}

type loginChecker struct {
	srv          *APIServer
	loopInterval time.Duration
}

func (l loginChecker) Check(ctx context.Context) {
	ticker := time.NewTicker(l.loopInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
		ok, err := l.srv.client.CheckLogin(context.Background())
		if err != nil {
			log.Ctx(ctx).Error().Err(err).Msg("check login status")
			continue
		}
		log.Ctx(ctx).Info().Bool("login", ok).Msg("check login")
		if ok {
			l.srv.login()
			continue
		}
		// 如果已经登录，并且当前失败了，那么就标记当前的server的登录的用户为退出状态
		if l.srv.IsLogin() {
			l.srv.logout()
			return
		}
	}
}
