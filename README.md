# BaaS

> 此项目用于管理 CaaS 平台上Mysql, Redis, Kafka 的集群实例，包含集群的创建，配置，服务暴露等，在部署此项目之前，请先部署好对应的 Operator

## Requirements

- openshift 3.11 +
- kubernetes 1.11 +
- Go 1.13.3

## mysql operator deployment ([mysql-operator](https://github.com/oracle/mysql-operator))

- 部署架构图

![avatar](images/mysql.jpg)

- 操作步骤如下：
	1. 创建 mysql-operator namespace

	        oc create ns mysql-operator
	2. 部署 mysql operator

	        oc create -f ./deployments/mysql-operator

## kafka operator deployment（[kafka-operator](https://github.com/banzaicloud/kafka-operator)）

- 架构图

![avatar](images/kafka.jpg)

- 具体部署步骤如下
	1. cert-manager 部署

			oc create -f ./deployments/cert-manager-openshift.yaml
	1. zookeeper-operator 部署

			oc new-project zookeeper-operator
			oc create -f ./deployments/zookeeper-operator
	1. kafka-operator 部署
		1. 执行下列命令获取当前集群配置的 caBundle

				oc config view --raw | grep certificate-authority-data | head -n 1 | awk -F: '{print $2}' | sed 's/\s//'
		1. 将 `deployments/kafka-operator/operator-validating-webhook.yaml` 文件中的 `caBundle` 属性值替换成上一步所获取的 caBunle
		1. 部署 operator

				oc new-project kafka-operator
				oc create -f ./deployments/kafka-operator

## build management image

1. 在本地编译

        export GO111MODULE="on"
        export GOPROXY="https://goproxy.cn,direct"
        
        go build -o ./app ./cmd/main.go
2. build docker 镜像

        docker build -t dbaas:v1.0 .

## deploy management

	oc adm policy add-cluster-role-to-user cluster-admin -z default -n default
	oc create -f ./deployments/dbaas.yaml

## Todo List

- [x] MySQL 集群的创建，删除，详情，列表，状态同步
- [x] Zookeeper 集群的创建，删除，详情，列表，状态同步
- [x] Kafka 集群的创建，删除，详情，列表，状态同步
- [x] Kafka topic 的创建，删除，详情，列表
- [ ] Redis 集群的创建，删除，详情，列表，状态同步