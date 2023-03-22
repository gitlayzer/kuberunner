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
| kuberunner | v0.0.3 | 2023/03/13 | Gitlayzer |     √      |

### 3：更新日志

```shell
本次更新已实现部分资源的所有操作功能，包括部分资源的增删改查，Pod资源的日志/终端功能已经实现，目前可以直接打包运行使用，只需要在config.yaml内配置需要的参数
```
### 4：config.yaml解释如下
```yaml
# API监听的端口
ListenAddress: ":80"
# 终端服务监听的端口
WsListenAddress: ":81"
# 前端登录验证的账号
Username: "admin"
# 前端登录验证的密码
Password: "admin"
# 可配置多个kubeconfig实现多集群管理
KubeConfigs:
  cluster1: "C:\\Users\\Administrator\\.kube\\config"
```

