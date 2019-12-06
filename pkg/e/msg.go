package e

var MsgFlags = map[int]string{
	SUCCESS:              "ok",
	ERROR:                "fail",
	INVALID_PARAMS:       "请求参数错误",
	ERROR_EXIST_MYSQL:    "mysql 实例已存在",
	CREATED:              "mysql 创建成功",
	ERROR_CREATE_MYSQL:   "mysql 创建失败",
	MYSQL_DOES_NOT_EXIST: "mysql 实例不存在",
	MYSQL_DELETED:        "mysql 实例已删除",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
