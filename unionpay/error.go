package unionpay

import (
	"fmt"
)

// BizErr 用于判断业务逻辑是否有错误
type BizErr struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

// bizErrCheck 检查返回码是否为成功 否则返回一个BizErr
func bizErrCheck(resp RspBase) error {
	if resp.ErrCode != "0000" && resp.ErrCode != "SUCCESS" {
		return &BizErr{
			Code: resp.ErrCode,
			Msg:  resp.ErrMsg,
		}
	}
	return nil
}

func (e *BizErr) Error() string {
	return fmt.Sprintf(`{"code":"%s","msg":"%s"}`, e.Code, e.Msg)
}
