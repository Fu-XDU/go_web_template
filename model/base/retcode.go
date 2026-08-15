package base

import "github.com/Fu-XDU/mingfu_go_common/base_response"

var (
	SUCCESS      = base_response.SUCCESS
	UnknownError = base_response.UnknownError
	Unauthorized = base_response.Unauthorized

	WrongParams           = base_response.NewRetCode(10002, "Wrong params")
	BindDataFailed        = base_response.NewRetCode(10003, "Bind data failed")
)
