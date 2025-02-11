package main

import (
        "crypto/tls"
        "fmt"
        "net/http"
        "time"
)

func main() {
        http.HandleFunc("/check/", handleSSLCheck)
        fmt.Println("服务启动在 :8080 端口...")
        http.ListenAndServe(":8080", nil)
}

func handleSSLCheck(w http.ResponseWriter, r *http.Request) {
        // 从URL中获取域名
        domain := r.URL.Path[len("/check/"):]
        if domain == "" {
                fmt.Printf("[Error] 未提供域名\n")
                http.Error(w, "请提供域名", http.StatusBadRequest)
                return
        }

        fmt.Printf("[Info] 正在检查域名: %s\n", domain)

        // 获取证书信息
        conn, err := tls.Dial("tcp", domain+":443", &tls.Config{
                InsecureSkipVerify: true,
        })
        if err != nil {
                fmt.Printf("[Error] 连接域名 %s 失败: %v\n", domain, err)
                http.Error(w, fmt.Sprintf("连接错误: %v", err), http.StatusInternalServerError)
                return
        }
        defer conn.Close()

        // 获取证书
        cert := conn.ConnectionState().PeerCertificates[0]

        // 计算剩余天数
        remainingDays := int(cert.NotAfter.Sub(time.Now()).Hours() / 24)

        // 打印证书详细信息
        fmt.Printf("[Info] 域名: %s, 证书有效期: %s 至 %s, 剩余天数: %d\n",
                domain,
                cert.NotBefore.Format("2006-01-02"),
                cert.NotAfter.Format("2006-01-02"),
                remainingDays,
        )

        fmt.Printf("[Info] 证书详细信息:\n"+
                "  主题: %s\n"+
                "  颁发者: %s\n"+
                "  序列号: %x\n"+
                "  版本: %d\n"+
                "  签名算法: %s\n"+
                "  DNS名称: %v\n",
                cert.Subject,
                cert.Issuer,
                cert.SerialNumber,
                cert.Version,
                cert.SignatureAlgorithm,
                cert.DNSNames,
        )

        // 构建响应
        response := fmt.Sprintf(`{
    "domain": "%s",
    "start_date": "%s",
    "end_date": "%s",
    "remaining_days": %d,
    "subject_cn": "%s",
    "issuer_cn": "%s"
}`, domain, cert.NotBefore.Format("2006-01-02"), cert.NotAfter.Format("2006-01-02"), remainingDays, cert.Subject.CommonName, cert.Issuer.CommonName)

        // 设置响应头
        w.Header().Set("Content-Type", "application/json")
        w.Write([]byte(response))
}
