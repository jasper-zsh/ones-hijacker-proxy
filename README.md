# ONES Hijacker Proxy

通过HTTPS中间人劫持所有指向DEV环境的API请求来实现通过DEV环境任意版本前端访问位于任意位置的后端API

## Quick start
### Proxy
本项目使用go mod处理依赖

```bash
git clone https://github.com/jasper-zsh/ones-hijacker-proxy
cd ones-hijacker-proxy
go get
go build
./ones-hijacker-proxy
```

目前监听端口是固定的，如果出现端口冲突先请调整占用端口的服务
```
HTTP Proxy: 6789
REST API(浏览器插件使用): 9090
```

### Chrome extension
搭配浏览器插件食用
```bash
git clone https://github.com/jasper-zsh/ones-hijacker
cd ones-hijacher
yarn
yarn build
```
构建好的插件在dist目录下，添加到Chrome即可使用
注意：要把Chrome的代理设置为本机6789端口才会生效！

### CA证书
与Fiddler/Charles拦截HTTPS流量类似，我们也需要安装CA证书并设置信任。
CA证书在goproxy包的根目录下面（`ca.pem`），不随本仓库提供。