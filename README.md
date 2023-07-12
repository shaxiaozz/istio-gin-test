### 前言
目前我们K8s正在使用ingress-nginx作为流量入口，为了更好的查看到微服务的流量，我们需要尝试迁移至istio。由于我们再ingress-nginx中实现了多种需求，目前我们需要将这些需求迁移至istio。主要的需求列表如下：
1、同一个域名，需要转发至前后端的Pod上。
2、转发至后端多个Pod时，需要配置会话亲和性
基于以上需求，我们将使用一个简单的应用模拟下流量，以下是这个应用的相关消息：
项目地址：https://github.com/shaxiaozz/istio-gin-test
模拟域名：istio-gin-test.test.com
后端API接口：/api/v1/version，/api/v2/version，/api/v3/version

### 以下操作都是基于istio已安装至K8s中为前提
### 一、部署istio-gin-test应用
#### 1.1 创建命名空间
```
[root@k8s-master ~]# kubectl create ns istio-gin-test
[root@k8s-master bookinfo]# kubectl label namespace istio-gin-test istio-injection=enabled  # 打上标签 istio-injection=enabled，自动注入 Sidecar
```

#### 1.2 创建前端应用
