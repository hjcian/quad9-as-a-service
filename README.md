# quad9-as-a-service

## dev notes
- test/benchmark (https://my.oschina.net/solate/blog/3034188)
  - `go test -bench=. -run=none -benchmem`
    - `-bench=.` 指的是當前路徑
    - `-run=none` run 原本是用來匹配想要執行的單元測試。不去設定會全跑，若想要全部都不跑就指定一個一定不存在的 pattern (none)
    - `-benchmem` 打開記憶體配置量量測