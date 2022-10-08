package e

var MsgFlags = map[int]string{
	SUCCESS: "ok",
	ERROR:   "fail",

	INVALID_PARAMS: "请求参数错误",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "token已超时",
	ERROR_GENERATE_TOKEN:           "token生成失败",
	ERROR_AUTH:                     "token错误",
	ERROR_EXIST_NICK:               "昵称已存在",
	ERROR_EXIST_USER:               "用户已存在",
	ERROR_NOT_EXIST_USER:           "用户不存在",
	ERROR_CHECK_PASSWORD_FAIL:      "密码错误",
	ERROR_FAIL_ENCRYPTION:          "加密失败",

	//
	ERROR_CALL_API: "调取机器人api接口失败",

	//
	ERROR_DATABASE: "数据库操作错误",

	//
	ERROR_UNMARSHAL_JSON: "解码JSON失败",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
