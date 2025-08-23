package handler

import (
	"net/http"

	"github.com/clin211/miniblog-v3/apps/user/api/internal/logic"
	"github.com/clin211/miniblog-v3/apps/user/api/internal/svc"
	"github.com/clin211/miniblog-v3/apps/user/api/internal/types"
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterRequest
		if err := httpx.Parse(r, &req); err != nil {
			response.WriteResponse(r.Context(), w, err)
			return
		}

		l := logic.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		if err != nil {
			response.WriteResponse(r.Context(), w, err)
		} else {
			response.WriteResponse(r.Context(), w, resp)
		}
	}
}
