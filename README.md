# myca

简单的证书链管理
通过改造cryptogen
---
初始化 ca 
myca init --dir=example.com 

example.com
    ca
        ca.example.com.pem
        ca.example.com.key
    tlsca
        tlsca.example.com.pem
        tlsca.example.com.key
    msp
        admincerts
            Admin@example.com.pem
        cacerts
            ca.example.com.pem
        tlscacerts
            tlsca.example.com.pem

    users
        Admin@example.com
            msp
                admincerts
                    Admin@example.com.pem
                cacerts
                    ca.example.com.pem
                tlscacerts
                    tlsca.example.com.pem
                keystore
                    Admin@example.com.key
                signcerts
                    Admin@example.com.pem
            tls
                tlsca.example.com.pem
                Admin@tlsca.example.com.pem
                Admin@tlsca.example.com.key

myca add-user --dir=example.com --name=User1
myca add-peer --dir=example.com --name=peer0
myca add-orderer --dir=example.com --name=orderer0

---

