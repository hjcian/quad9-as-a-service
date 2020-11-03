# quad9-as-a-service - q9aas

**What is q9aas (this repo)?**
- `Go` 的個人小練習
- 做一個 REST API server 供其他服務方便利用 Quad9 的情資
- 直接向 `9.9.9.9` 請求一個惡意網域，會得到 `NXDOMAIN` 的結果
- 故此 server 實際上是非同步地向 `9.9.9.9` 與 `9.9.9.10` 蒐集結果並判斷是否真的是惡意、而非真的不存在

**What is Quad9?**
> [Quad9](https://www.quad9.net/) is a free DNS query data platform.
- a DNS server: `9.9.9.9`
- 它的特色是會幫你做 security protection，直接在 DNS query 後就幫你 block 惡意網域
- 更多的疑問可在官方 [FAQ](https://www.quad9.net/faq/) 上得到解答


## Install, Test, Build and Run
> - test on go 1.14 and latter

```shell
git clone https://github.com/hjcian/quad9-as-a-service.git
cd quad9-as-a-service
make install
make test   # run basic unit tests
make build  # build the executable
make run    # execute the executable
```

## API
- 簡單使用 GET 向 server 請求

***Request***
```shell
curl --request GET \
  --url 'http://localhost:12345/checkBlocklist?hostname=aaa.com'
```

***Response***
- 回傳一個 JSON body，告訴你這個網域是否有被 blocked
```
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Tue, 03 Nov 2020 05:24:00 GMT
Content-Length: 17
Connection: close

{
  "blocked": false
}
```
