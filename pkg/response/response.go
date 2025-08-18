package response

import (
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"github.com/clin211/miniblog-v3/pkg/validate"
)

// BaseResponse 定义统一的响应结构。
// BaseResponse 包含业务码 Code、提示信息 Message、数据 Data 以及扩展字段 Reason。
// 所有对外响应均应包含这四个字段，以保持结构一致性。
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Reason  interface{} `json:"reason"`
}

// SuccessCtx 返回统一成功响应。
// 参数 message 为提示信息，data 为业务数据，reason 用于携带扩展信息（可为 nil）。
func SuccessCtx(ctx context.Context, w http.ResponseWriter, message string, data interface{}, reason interface{}) {
	resp := BaseResponse{
		Code:    0,
		Message: message,
		Data:    data,
		Reason:  reason,
	}
	httpx.OkJsonCtx(ctx, w, resp)
}

// FailCtx 返回统一失败响应，HTTP 状态码固定 200，业务码由 code 指定。
// 为满足需求示例，若调用方希望表示资源不存在，可传入 code=404，message="Not Found"。
// reason 用于附加错误上下文信息。
func FailCtx(ctx context.Context, w http.ResponseWriter, code int, message string, reason interface{}) {
	resp := BaseResponse{
		Code:    code,
		Message: message,
		Data:    nil,
		Reason:  reason,
	}
	httpx.OkJsonCtx(ctx, w, resp)
}

// NotFoundCtx 返回示例中的未找到响应结构：{code:404, message:"Not Found", data:null, reason:any}。
func NotFoundCtx(ctx context.Context, w http.ResponseWriter, reason interface{}) {
	FailCtx(ctx, w, 404, "Not Found", reason)
}

// JsonBaseResponseCtx 模仿 go-zero 扩展用法的统一响应入口。
// 语义：
//   - 当 v 为 error：
//     1) 若为参数校验错误（validate.ValidationErrors），返回 code=0，message 由调用者不传时默认“参数校验失败”，
//     data=nil，reason 为校验错误切片（[]string）。为了简化调用，此处 message 固定为“参数校验失败”。
//     2) 其他错误，返回示例中的失败结构：code=404, message="Not Found", data=nil，reason 为错误信息字符串。
//   - 当 v 非 error：按成功结构返回：code=0, message="ok"，data=v，reason=nil。
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v interface{}) {
	switch val := v.(type) {
	case error:
		// 参数校验失败（ValidationErrors）需要以列表形式输出 reason。
		if ve, ok := val.(validate.ValidationErrors); ok {
			var reasons []string
			for _, item := range ve {
				if item == nil {
					continue
				}
				reasons = append(reasons, item.Error())
			}
			// 按需求示例：参数校验失败也使用 code=0。
			httpx.OkJsonCtx(ctx, w, BaseResponse{
				Code:    0,
				Message: "参数校验失败",
				Data:    nil,
				Reason:  reasons,
			})
			return
		}

		// 其余错误：按失败示例返回 404 + Not Found，reason 为错误字符串，data 为 nil。
		httpx.OkJsonCtx(ctx, w, BaseResponse{
			Code:    404,
			Message: "Not Found",
			Data:    nil,
			Reason:  val.Error(),
		})
		return

	default:
		// 成功：按成功示例返回 code=0, message="ok"，data 为传入值。
		httpx.OkJsonCtx(ctx, w, BaseResponse{
			Code:    0,
			Message: "ok",
			Data:    v,
			Reason:  nil,
		})
		return
	}
}
