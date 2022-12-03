# traceSysBackend

## 说明

这是 Leo 的毕业设计项目，所有的代码均未经过严格测试，请谨慎使用。

本项目为前后端分离，前端需要单独部署，请详见 [tracesys_vue](https://github.com/sjlleo/tracesys_vue)

## 安装指南

对于 TraceSystem 的后端，我们不推荐自动安装，建议您手动安装（其实是懒得写脚本了）

首先，您应该安装 MySQL（>=5.7），然后创建一个数据库

接下来，将本项目 Clone 下来

```bash
git clone https://github.com/sjlleo/traceSysBackend.git
cd traceSysBackend
```

然后，您应该修改 database 文件夹下的 `db.go` 文件

```Go
username := "tracesys"         //账号
password := "hZanGEC8FjbTh2x8" //密码
host := "localhost"            //数据库地址，可以是IP或者域名
port := 3306                   //数据库端口
Dbname := "tracesys"           //数据库名
timeout := "10s"               //连接超时，10秒
```

安装上述提示完成对 `db.go` 的数据库配置信息修改

修改完成后，您可以将此项目编译为二进制文件了

```bash
go build .
```

您可以将项目放入系统的 Systemctl 中管理，下面是一个示范

```
# /etc/systemd/system/traceSysBackend.service

[Unit]
Description=Trace System Backend
After=network.target
[Service]
Type=simple
Restart=always
WoringDirectory=Your_Working_Path_Here
ExecStart=Your_Working_Path_Here/traceSysBackend
[Install]
WantedBy=multi-user.target
```

```bash
systemctl daemon-reload
systemctl enable traceSysBackend.service
systemctl start traceSysBackend.service
systemctl status traceSysBackend.service
```

然后项目将运行在 `50888` 端口，请您使用 Nginx 进行转发。

```nginx
# 在 nginx.conf 的 Server 配置项中添加
location ^~ /
{
    proxy_pass http://localhost:50888;
    proxy_set_header Host localhost;
    proxy_set_header X-Real-IP $remote_addr;
    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    proxy_set_header REMOTE-HOST $remote_addr;
}
```

本项目强制使用 HTTPS 进行节点和后端之间的数据传输，请务必为 Nginx 部署 SSL 证书~

> PS: 可以直接使用 `acme.sh` 快捷为域名颁发证书。


## 项目截图

![p2](https://user-images.githubusercontent.com/13616352/205425006-00911fff-7d63-46d5-b95d-7073149a8224.jpeg)
![p1](https://user-images.githubusercontent.com/13616352/205425004-b9e26492-1d71-421f-9881-68b3edb89087.jpeg)
![p3](https://user-images.githubusercontent.com/13616352/205424993-58337531-f663-42c3-aacc-f93402c77a0b.jpeg)
![p4](https://user-images.githubusercontent.com/13616352/205424997-a5aed927-5b76-4943-a413-db15ad4a1baf.jpeg)
![p5](https://user-images.githubusercontent.com/13616352/205425002-d27ea472-478a-4f88-bd97-3d49c937a99e.jpeg)
