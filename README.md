# check-ssl
通过接口的方式检查ssl证书过期时间

##  使用方法
```shell
go build -o check-ssl main.go #会生成一个check-ssl文件
./check-ssl #运行
```


## 域名接口监测

http://localhost:8080/check/域名

在线接口，可不用部署服务直接使用
https://ssl.kubehan.cn/check/www.kubehan.cn

## 返回类型案例
```json
{
    "domain": "www.kubehan.cn",
    "start_date": "2025-02-06",
    "end_date": "2025-05-07",
    "remaining_days": 85,
    "subject_cn": "kubehan.cn",
    "issuer_cn": "R10"
}
```
