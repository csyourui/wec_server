package server

import (
	"github.com/csyourui/wechat_server/pkg/log"
	"net/http"
)

func Server(server *http.Server) error {
	log.Logger.Debug("start server", server.Addr)
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Logger.Error("start fail", err)
		return err
	}
	return nil
}
