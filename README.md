# casbin-playground

[Casbin](https://casbin.org/) を理解するために作成したサンプルプログラムです

## 基本仕様

### Model

ACL ([`basic_model.conf`](https://github.com/casbin/casbin/blob/master/examples/basic_model.conf))を設定しています

### Adapter

ビルトインの File Adapter を用いており、`policy.csv` が用いられます

### Watcher

以下のリクエストを送信すると更新が通知される簡易的な `Watcher` を実装しています

```
POST /update
```

## Watcher の動作確認

```console
# bob は policy.csv に存在しておらず deny される
$ http :8080/check query==bob,data1,read
HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Mon, 09 Jun 2025 06:37:25 GMT

{
    "ok": false
}

# policy.csv に p, bob, data1, read を追記

# そのままリクエストすると、更新前のまま
$ http :8080/check query==bob,data1,read
HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Mon, 09 Jun 2025 06:38:25 GMT

{
    "ok": false
}

# Watcher 経由でコールバックを発火
$ http POST :8080/update
HTTP/1.1 204 No Content
Content-Type: application/json
Date: Mon, 09 Jun 2025 06:39:25 GMT

# 更新が反映される
$ http :8080/check query==bob,data1,read
HTTP/1.1 200 OK
Content-Length: 13
Content-Type: application/json
Date: Mon, 09 Jun 2025 06:40:25 GMT

{
    "ok": true
}
```
