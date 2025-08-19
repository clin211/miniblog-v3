package response

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/clin211/miniblog-v3/pkg/errorx"
	"github.com/clin211/miniblog-v3/pkg/validate"
)

// responseBody 定义统一的响应结构体（不包含 HTTP 字段）
type responseBody struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Reason  string      `json:"reason"`
}

// SuccessCtx 返回成功结果，HTTP 200。
func SuccessCtx(ctx context.Context, w http.ResponseWriter, message string, data any, reason string) {
	body := responseBody{Code: 0, Message: message, Data: data, Reason: reason}
	httpx.WriteJsonCtx(ctx, w, http.StatusOK, body)
}

// FailCtx 返回失败结果，默认 HTTP 400。
func FailCtx(ctx context.Context, w http.ResponseWriter, code int, message string, reason string) {
	body := responseBody{Code: code, Message: message, Data: nil, Reason: reason}
	httpx.WriteJsonCtx(ctx, w, http.StatusBadRequest, body)
}

// WriteResponse 统一入口：
// - 如果 v 为 error：
//   - errorx.Errno：按 Errno.HTTP 写入 HTTP 状态，业务结构取 Code/Message/Data/Reason。
//   - validate.ValidationErrors / validate.ValidationErrorsWithCode：HTTP 400，code=0，message="参数校验失败"，reason=[]string。
//   - 其他 error：HTTP 404，code=404，message="Not Found"，reason=err.Error()。
//
// - 如果 v 非 error：HTTP 200，code=0，message="ok"，data=v。
func WriteResponse(ctx context.Context, w http.ResponseWriter, v any) {
	if v == nil {
		body := responseBody{Code: 0, Message: "ok", Data: nil, Reason: ""}
		httpx.WriteJsonCtx(ctx, w, http.StatusOK, body)
		return
	}

	if err, ok := v.(error); ok {
		// errorx.Errno 分支
		if e, ok := err.(*errorx.Errno); ok {
			body := responseBody{Code: e.Code, Message: e.Message, Data: e.Data, Reason: e.Reason}
			httpx.WriteJsonCtx(ctx, w, e.HTTP, body)
			return
		}

		// validate.ValidationErrorsWithCode
		if es, ok := err.(validate.ValidationErrorsWithCode); ok {
			var reasonMsg string
			for i, ve := range es {
				if i > 0 {
					reasonMsg += "; "
				}
				reasonMsg += ve.Error()
			}
			body := responseBody{Code: errorx.ErrInvalidParameter.Code, Message: "参数校验失败", Data: nil, Reason: reasonMsg}
			httpx.WriteJsonCtx(ctx, w, http.StatusBadRequest, body)
			return
		}

		// validate.ValidationErrors（无 code 版本）
		if es, ok := err.(validate.ValidationErrors); ok {
			var reasonMsg string
			for i, ve := range es {
				if i > 0 {
					reasonMsg += "; "
				}
				reasonMsg += ve.Error()
			}
			body := responseBody{Code: errorx.ErrInvalidParameter.Code, Message: "参数校验失败", Data: nil, Reason: reasonMsg}
			httpx.WriteJsonCtx(ctx, w, http.StatusBadRequest, body)
			return
		}

		// 其他 error -> Not Found
		body := responseBody{Code: errorx.ErrResourceNotFound.Code, Message: errorx.ErrResourceNotFound.Message, Data: nil, Reason: err.Error()}
		httpx.WriteJsonCtx(ctx, w, http.StatusNotFound, body)
		return
	}

	// 成功对象
	body := responseBody{Code: 0, Message: "ok", Data: v, Reason: ""}
	httpx.WriteJsonCtx(ctx, w, http.StatusOK, body)
}
