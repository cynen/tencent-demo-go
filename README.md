# tencent-demo-go

微信相关联的代码仓,涉及公众号,小程序 等 演示代码.

Notice: 此项目自动从此 gitee 同步到 [github](https://github.com/cynen/tencent-demo-go).如果有改动需求,直接更改gitee即可.

## 企业微信对接通义千问
代码: [qywx-tyqw](qywx-tyqw)

### 安装教程

#### 宿主机运行:
1-5 步骤与容器运行不一样.后面一样的.

1. 拉取代码
2. 本地编译: 
```
cd qywx-tyqw
go build -o qywx

# 运行项目
./qywx -c config.yml

也可以省略配置
./qywx 默认加载当前目录下的config.yml文件.
```

#### 容器运行:
1. 此项目可以直接打包成docker镜像运行.
2. 拉取代码.获取Dockerfile文件
3. 修改config.yml中的配置.(非必须,后续也可以挂载进去,百度搜索容器挂载文件.)
4. 构建镜像. `docker build -t tyqw .`
5. 运行镜像
```
docker run -d --name tyqw -p 8888:8888 tyqw
这里的端口,根据实际配置的业务端口来修改.
左边的 8888 是 在宿主机上开启的监听端口,右侧是容器内部应用的端口.
```
6. 运行完成后,在服务器配置域名.建议https. 通过nginx做代理.
```
比如: 代理成 https://qywx.test.cn  --> http://127.0.0.1:8888
重点是配置ssl证书.
``` 

7. 去企业微信后台,配置回调地址. 
```
    企业微信>应用管理>自建应用>接收消息>URL(回调地址)
    回调地址: https://qywx.test.cn/qywxpush
    这里的qywxpush是我们在项目里的url路劲,main.go中可以查阅到.
```
8. 点击保存,会自动校验地址是否可用.




#### 特技

1.  使用 Readme\_XXX.md 来支持不同的语言，例如 Readme\_en.md, Readme\_zh.md
2.  Gitee 官方博客 [blog.gitee.com](https://blog.gitee.com)
3.  你可以 [https://gitee.com/explore](https://gitee.com/explore) 这个地址来了解 Gitee 上的优秀开源项目
4.  [GVP](https://gitee.com/gvp) 全称是 Gitee 最有价值开源项目，是综合评定出的优秀开源项目
5.  Gitee 官方提供的使用手册 [https://gitee.com/help](https://gitee.com/help)
6.  Gitee 封面人物是一档用来展示 Gitee 会员风采的栏目 [https://gitee.com/gitee-stars/](https://gitee.com/gitee-stars/)
