package e

var MsgFlags = map[int]string{
	SUCCESS:                             "ok",
	ERROR:                               "fail",
	INVALID_PARAMS:                      "请求参数错误",
	ERROR_EXIST_MYSQL:                   "mysql cluster 实例已存在",
	CREATED:                             "cluster 创建成功",
	ERROR_CREATE_MYSQL:                  "mysql cluster 创建失败",
	MYSQL_DOES_NOT_EXIST:                "mysql cluster 实例不存在",
	MYSQL_DELETED:                       "mysql cluster 实例已删除",
	ERROR_EXIST_ZOOKEEPER:               "zookeeper cluster 实例已存在",
	ERROR_CREATE_ZOOKEEPER:              "zookeeper cluster 创建失败",
	ZOOKEEPER_DOES_NOT_EXIST:            "zookeeper cluster 实例不存在",
	ZOOKEEPER_DELETED:                   "zookeeper cluster 实例已删除",
	ERROR_COUNT_ZOOKEEPER_CLUSTER_FAIL:  "zookeeper 实例总数统计失败",
	ERROR_GET_ZOOKEEPER_CLUSTER_FAIL:    "zookeeper 实例获取失败",
	ERROR_CHECK_KAFKA_EXIST_FAIL:        "kafka 实例状态检查失败",
	ERROR_CHECK_ZOOKEEPER_EXIST_FAIL:    "zookeeper 实例状态检查失败",
	MYSQL_DELETED_FAILED:                "mysql 实例删除失败",
	ERROR_CHECK_MYSQL_EXIST_FAIL:        "mysql 实例状态检查失败",
	ERROR_COUNT_MYSQL_CLUSTER_FAIL:      "mysql 实例总数统计失败",
	ERROR_GET_MYSQL_CLUSTER_FAIL:        "mysql 实例获取失败",
	ERROR_EXIST_KAFKA:                   "kafka 实例已存在",
	ERROR_CREATE_KAFKA:                  "kafka 实例创建失败",
	KAFKA_DOES_NOT_EXIST:                "kafka 实例不存在",
	KAFKA_DELETED:                       "kafka 实例已删除",
	ERROR_COUNT_KAFKA_CLUSTER_FAIL:      "kafka 实例总数统计失败",
	ERROR_GET_KAFKA_CLUSTER_FAIL:        "kafka 实例获取失败",
	ZOOKEEPER_DELETED_FAILED:            "zookeeper 实例删除失败",
	KAFKA_DELETED_FAILED:                "kafka 实例删除失败",
	CHECK_ZOOKEEPER_USAGE_STATUS_FAILED: "zookeeper 集群的使用状态检查失败",

	ERROR_EXIST_KAFKA_TOPIC:                              "kafka topic 已存在",
	ERROR_CREATE_KAFKA_TOPIC:                             "kafka topic 创建失败",
	KAFKA_TOPIC_DOES_NOT_EXIST:                           "kafka topic 不存在",
	KAFKA_TOPIC_DELETED:                                  "kafka topic 已删除",
	ERROR_COUNT_KAFKA_TOPIC_FAIL:                         "kafka topic 总数统计失败",
	ERROR_GET_KAFKA_TOPIC_FAIL:                           "kafka topic 获取失败",
	ERROR_CHECK_KAFKA_TOPIC_EXIST_FAIL:                   "kafka topic 实例状态检查失败",
	KAFKA_TOPIC_DELETED_FAILED:                           "kafka topic 实例删除失败",
	KAFKA_TOPIC_REF_CLUSTER_DOES_NOT_EXIST:               "kafka topic 指向的集群不存在",
	KAFLA_CLUSTER_IS_USING_THE_CURRENT_ZOOKEEPER_CLUSTER: "有 kafka 集群正在使用当前 zookeeper 集群",

	K8S_MYSQL_CLUSTER_CREATE_FAILED:         "k8s 中的 mysql 实例创建失败",
	K8S_ZOOKEEPER_CLUSTER_CREATE_FAILED:     "k8s 中的 zookeeper 实例创建失败",
	K8S_KAFKA_CLUSTER_CREATE_FAILED:         "k8s 中的 kafka 实例创建失败",
	K8S_MYSQL_CLUSTER_DELETE_FAILED:         "k8s 中的 mysql 实例删除失败",
	K8S_ZOOKEEPER_CLUSTER_DELETE_FAILED:     "k8s 中的 zookeeper 实例删除失败",
	K8S_KAFKA_CLUSTER_DELETE_FAILED:         "k8s 中的 kafka 实例删除失败",
	K8S_MYSQL_PASSWORD_SECRET_DELETE_FAILED: "k8s 中的 mysql 的密码 secret 删除失败",
	K8S_MYSQL_CONFIGMAP_DELETE_FAILED:       "k8s 中的 mysql 的配置 configmap 删除失败",
	K8S_MYSQL_ROUTER_DELETE_FAILED:          "k8s 中的 mysql router 删除失败",
	K8S_MYSQL_ROUTER_SERVICE_DELETE_FAILED:  "k8s 中的 mysql router 的 sevice 删除失败",
	K8S_MYSQL_PVC_DELETE_FAILED:             "k8s 中的 mysql 实例使用的 pvc 删除失败",
	K8S_MYSQL_PASSWORD_SECRET_CREATE_FAILED: "k8s 中的 mysql 的密码 secret 创建失败",
	K8S_MYSQL_CONFIGMAP_CREATE_FAILED:       "k8s 中的 mysql 的配置 configmap 创建失败",
	K8S_MYSQL_ROUTER_CREATE_FAILED:          "k8s 中的 mysql router 删除失败",
	K8S_MYSQL_ROUTER_SERVICE_CREATE_FAILED:  "k8s 中的 mysql router 的 sevice 删除失败",

	K8S_KAFKA_TOPIC_CREATE_FAILED: "k8s 中的 kafka topic 创建失败",
	K8S_KAFKA_TOPIC_DELETE_FAILED: "k8s 中的 kafka topic 删除失败",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
