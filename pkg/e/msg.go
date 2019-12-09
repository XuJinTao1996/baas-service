package e

var MsgFlags = map[int]string{
	SUCCESS:                                 "ok",
	ERROR:                                   "fail",
	INVALID_PARAMS:                          "请求参数错误",
	ERROR_EXIST_MYSQL:                       "mysql cluster 实例已存在",
	CREATED:                                 "mysql cluster 创建成功",
	ERROR_CREATE_MYSQL:                      "mysql cluster 创建失败",
	MYSQL_DOES_NOT_EXIST:                    "mysql cluster 实例不存在",
	MYSQL_DELETED:                           "mysql cluster 实例已删除",
	MYSQL_DELETED_FAILED:                    "mysql 实例删除失败",
	ERROR_CHECK_MYSQL_EXIST_FAIL:            "mysql 实例状态检查失败",
	ERROR_COUNT_MYSQL_CLUSTER_FAIL:          "MYSQL 实例总数统计失败",
	ERROR_GET_MYSQL_CLUSTER_FAIL:            "mysql 实例获取失败",
	K8S_MYSQL_CLUSTER_CREATE_FAILED:         "k8s 中的 mysql 实例创建失败",
	K8S_MYSQL_CLUSTER_DELETE_FAILED:         "k8s 中的 mysql 实例删除失败",
	K8S_MYSQL_PASSWORD_SECRET_DELETE_FAILED: "k8s 中的 mysql 的密码 secret 删除失败",
	K8S_MYSQL_CONFIGMAP_DELETE_FAILED:       "k8s 中的 mysql 的配置 configmap 删除失败",
	K8S_MYSQL_ROUTER_DELETE_FAILED:          "k8s 中的 mysql router 删除失败",
	K8S_MYSQL_ROUTER_SERVICE_DELETE_FAILED:  "k8s 中的 mysql router 的 sevice 删除失败",
	K8S_MYSQL_PVC_DELETE_FAILED:             "k8s 中的 mysql 实例使用的 pvc 删除失败",
	K8S_MYSQL_PASSWORD_SECRET_CREATE_FAILED: "k8s 中的 mysql 的密码 secret 创建失败",
	K8S_MYSQL_CONFIGMAP_CREATE_FAILED:       "k8s 中的 mysql 的配置 configmap 创建失败",
	K8S_MYSQL_ROUTER_CREATE_FAILED:          "k8s 中的 mysql router 删除失败",
	K8S_MYSQL_ROUTER_SERVICE_CREATE_FAILED:  "k8s 中的 mysql router 的 sevice 删除失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
