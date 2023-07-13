### 前言
目前我们K8s正在使用ingress-nginx作为流量入口，为了更好的查看到微服务的流量，我们需要尝试迁移至istio。由于我们再ingress-nginx中实现了多种需求，目前我们需要将这些需求迁移至istio。主要的需求列表如下：  
1. 同一个域名，需要转发至前后端的Pod上。
2. 转发至后端多个Pod时，需要配置会话亲和性
3. 需要配置跨域

基于以上需求，我们将使用一个简单的应用模拟下流量，以下是这个应用的相关消息：  
项目地址：https://github.com/shaxiaozz/istio-gin-test  
模拟域名：istio-gin-test.test.com  
后端API接口：/api/v1/version，/api/v2/version，/api/v3/version  

**以下操作都是基于istio已安装至K8s中为前提**  
### 一、部署istio-gin-test应用
#### 1.1 创建命名空间
```
[root@k8s-master ~]# kubectl create ns istio-gin-test
[root@k8s-master ~]# kubectl label namespace istio-gin-test istio-injection=enabled  # 打上标签 istio-injection=enabled，自动注入 Sidecar
```

1.2 创建前端应用

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-fe-deploy.yaml
[root@k8s-master istio-gin-test]# kubectl get pods -n istio-gin-test


1.3 创建后端应用

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-deploy.yaml
[root@k8s-master istio-gin-test]# kubectl get pods -n istio-gin-test -o wide

测试后端接口是否正常
[root@k8s-master istio-gin-test]# curl 10.244.235.230:8080/api/v1/version
{"message":"Current API version: v1, Current hostname is: istio-gin-test-84cb778764-dfvff"}
[root@k8s-master istio-gin-test]# curl 10.244.235.230:8080/api/v2/version
{"message":"Current API version: v2, Current hostname is: istio-gin-test-84cb778764-dfvff"}
[root@k8s-master istio-gin-test]# curl 10.244.235.230:8080/api/v3/version
{"message":"Current API version: v3, Current hostname is: istio-gin-test-84cb778764-dfvff"}


二、配置istio路由
基本概念
Gateway（Gateway）： 
servers：定义入口点列表
selector：选择器，用于通过label选择集群中Istio网关的Pod
虚拟服务（Virtual Service）：
定义路由规则，匹配请求
描述满足条件的请求去哪里
目标规则（Destination Rule）：
定义子集、策略
描述到达目标的请求怎么处理
2.1 创建Gateway

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-gateway.yaml
[root@k8s-master istio-gin-test]# kubectl get gateway -n istio-gin-test


2.2 创建VirtualService

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-route.yaml
[root@k8s-master istio-gin-test]# kubectl get virtualservices.networking.istio.io -n istio-gin-test


2.4 浏览器测试访问
由于我们内部没有LoadBalancer，因此需要将istio-ingressgateway servicename的类型修改为NodePort，参考链接：https://istio.io/latest/zh/docs/tasks/traffic-management/ingress/ingress-control/#determining-the-ingress-ip-and-ports
[root@k8s-master istio-gin-test]# kubectl get svc -n istio-ingress


前端页面访问正常：

后端V1接口访问正常：


后端V2接口访问正常：

后端V3接口访问正常：

v1，v2，v3按钮接口都能正常返回，下面我们将测试下会话亲和性的需求了

三、istio会话亲和性
3.1 创建DestinationRule
官网文档：https://istio.io/latest/docs/reference/config/networking/destination-rule/#LoadBalancerSettings

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-sticky.yaml
[root@k8s-master istio-gin-test]# kubectl get destinationrules.networking.istio.io -n istio-gin-test


3.2 测试
浏览器疯狂访问第四个按钮，观测主机名称是否发生变化


四、istio 跨域配置
4.1 新建VirtualService
由于目前后端的域名为：istio-gin-test.test.com，因此我们要模拟前端跨域的话，得新建一个新域名给到前端：istio-gin-test-cross.test.com

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-cross-route.yaml
[root@k8s-master istio-gin-test]# kubectl get virtualservices.networking.istio.io -n istio-gin-test

浏览器访问新域名，并触发后端请求，存在跨域错误

4.2 VirtualService配置跨域
将后端域名的VirtualService添加跨域配置（官网文档：https://istio.io/latest/docs/reference/config/networking/virtual-service/#CorsPolicy）

[root@k8s-master istio-gin-test]# kubectl apply -f istio-gin-test-cross-route-fix.yaml
[root@k8s-master istio-gin-test]# kubectl get virtualservices.networking.istio.io -n istio-gin-test

浏览器重新访问跨域的域名，跨域问题已解决

