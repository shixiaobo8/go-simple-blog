## 介绍
个人学习使用的blog，后端采用go语言框架beego开发API，前端使用vue + iview-admin开发，前后端分离的小型博客。代码仅做参考，不可直接用于生产环境

## 部署

    环境：
    golang （尽量使用最新版本）需要go mod支持
    mysql
    redis
    node

后端使用mod管理依赖

配置go env

    GO111MODULE=on
    GOPROXY=https://goproxy.io
goland开发软件需要开启Go Modules，填写proxy为 https://goproxy.io

后端go运行在9090端口，在web目录下查看config/index.js
```
baseUrl: {
    dev: 'http://localhost:9090',
    pro: 'https://produce.com'     # 线上环境
},
```

1、克隆代码后根目录执行
```shell
go mod vendor
go mod download
```
    
2、前端代码在web目录下，使用vue-cli3构建。

npm慢可以使用 npm install -g cnpm --registry=https://registry.npm.taobao.org
然后即可使用cnpm
```shell
cnpm install # 安装依赖，开发模式下，使用npm run dev开启项目
cnpm run build # 构建完成后会生成dist目录，可以当做web根目录
```

配合nginx部署，需要注意vue-router设置，详见[vue-router](https://router.vuejs.org/zh/guide/essentials/history-mode.html#%E5%90%8E%E7%AB%AF%E9%85%8D%E7%BD%AE%E4%BE%8B%E5%AD%90)开发手册

3、后端接口验证可用session或jwt，默认使用jwt，使用session需要配置redis；图片上传可存储到服务器或阿里云oss，需要配置oss账号