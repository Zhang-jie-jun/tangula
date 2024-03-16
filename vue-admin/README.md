## Install
```bush
// install dependencies
npm install
```
## Run
### Development
```bush
npm run dev
```

### 本地开发代理配置

修改根目录下的`proxyPath.js`中的proxyPath对应的代理服务器地址

```
module.exports = {
  proxyPath: 'http://10.10.5.242:8080'
}
```
启动命令使用`npm run proxy`


### 生产部署
```
1. 生产用nginx代理
   server {
        listen       8088         #前端的端口号，不是后端接口端口号;
        server_name  10.4.108.47  #后端服务地址;
        location / {
            root   /tangula/dist;   #前端资源目录
            index  index.html index.htm;
            try_files $uri $uri/ /index.html;
        }
  }

2. config目录下index.js文件中，baseUrl.pro路径配置为后端服务地址

3. 本地 npm run build

4.tar -czvf xxxxxx.tar.gz dist/

5. 将打完的包上传nginx前端资源目录

6.删除原dist/目录，解压
```

