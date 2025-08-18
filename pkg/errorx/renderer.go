package errorx

import (
	"github.com/clin211/miniblog-v3/pkg/response"
	"github.com/clin211/miniblog-v3/pkg/validate"
)

// Render 将任意值渲染为统一的 BaseResponse（不负责写出 HTTP）。
//
// 约定保持与 docs/03 响应结构 一致：
//   - 成功：code=0, message="ok", data=v, reason=nil
//   - 参数校验失败：code=0, message="参数校验失败", data=nil, reason=[]string
//   - 业务错误 errorx.Error：code=e.Code, message=e.Message, data=nil, reason=e.Reason 或 e.Cause.Error()
//   - 其他 error：按 NotFound 语义示例返回 code=404, message="Not Found", data=nil, reason=err.Error()
func Render(v interface{}) response.BaseResponse {
	switch val := v.(type) {
	case nil:
		return response.BaseResponse{Code: int(CodeOK), Message: "ok", Data: nil, Reason: nil}

	case validate.ValidationErrors:
		var reasons []string
		for _, item := range val {
			if item == nil {
				continue
			}
			reasons = append(reasons, item.Error())
		}
		return response.BaseResponse{Code: int(CodeOK), Message: "参数校验失败", Data: nil, Reason: reasons}

	case *Error:
		if val == nil {
			return response.BaseResponse{Code: int(CodeOK), Message: "ok", Data: nil, Reason: nil}
		}
		reason := val.Reason
		if reason == nil && val.Cause != nil {
			reason = val.Cause.Error()
		}
		return response.BaseResponse{Code: int(val.Code), Message: val.Message, Data: nil, Reason: reason}

	case error:
		// 其他错误按文档中的示例 Not Found 处理
		return response.BaseResponse{Code: int(CodeNotFound), Message: "Not Found", Data: nil, Reason: val.Error()}

	default:
		return response.BaseResponse{Code: int(CodeOK), Message: "ok", Data: v, Reason: nil}
	}
}
