# Go REST API Server
Goで作ったTODOアプリのREST API Server  
Go, Gin, Gorm, MySQL, Docker, Clean Architecture etc.
## サーバー起動方法
(VSCodeの場合)
1. このレポジトリをclone
1. VScodeの拡張プラグイン[Remote Container](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)をインストール
1. VSCodeのコマンドパレットから```Remote-Containers: Open Folder in Container...```を選択
1. ```go run main.go```
## API仕様
### エンドポイント
```localhost:8080/api/v1```
### データ形式
```application/json```
### 認証
認証が必要なリクエストではcookieに
```"id":(任意のuserid)```
を付与してからAPIリクエストを送る．(将来的にはJWT認証を使った認証を実装したい)
###  エラーレスポンス
エラーが発生した場合は以下のようなJSONが帰ってくる．
```
{
    "code": 404,
    "message": "user not found"
}
```

#### 共通のエラーレスポンス
| code | message | 補足 |
|:--- |:---:| ---:|
| 400 | bad request | 不正なJSON |
| 500 | internal server error | 不明な内部エラー |

### 
