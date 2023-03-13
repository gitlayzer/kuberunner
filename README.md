# Kuberunner

<!-- TOC -->
* [1：此项目只用于重构Kubemanager](https://github.com/gitlayzer/kubemanager)
* [2：发布详情](#4发布详情)
* [3：更新日志](#5更新日志)
<!-- TOC -->


### 2：发布详情


|    名称    |  版本  |  发布时间  |  发布人   | 是否更新UI |
| :--------: | :----: | :--------: | :-------: | :--------: |
| kuberunner | v0.0.1 | 2023/03/10 | Gitlayzer |     ×      |
| kuberunner | v0.0.2 | 2023/03/13 | Gitlayzer |     ×      |

### 3：更新日志

```shell
1：本次更新功能如下
2：API Token认证（但只是有一个认证框架，并非实际认证，因为没有还没有接入数据库，所以认证暂时还未落实）
3：CORS跨域（目前已经可以实现启动API服务之后前端调用此API并不会出现CORS跨域问题）
4：基于以前的Kubemanager配置，此次版本针对传入配置文件进行了变化，本次版本使用的配置文件为yaml格式
5：本次更新添加了deployment/daemonset/statefulset/service/ingress/configmap/secret/pv/pvc/sc等资源的部分操作
6：后续更新会推出Websocket登录Pod的操作
```

