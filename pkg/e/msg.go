package e

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	ERROR:              "fail",
	INVALID_PARAMS:     "请求参数错误",
	ERROR_ESIST_MYSQL:  "mysql 实例已存在",
	CREATED:            "创建成功",
	ERROR_CREATE_MYSQL: "创建失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
