# API仕様
## エンドポイント
テスト環境では
```localhost:8080/api/v1```
## データ形式
```application/json```
## 認証
認証が必要なリクエストではcookieに
```"id":(任意のuserid)```
を付与してからAPIリクエストを送る．(将来的にはJWT認証を使った認証を実装したい)
##  エラーレスポンス
エラーが発生した場合は以下のようなJSONが帰ってくる．
```
{
    "code": 404,
    "message": "user not found"
}
```

### 共通のエラーレスポンス
| code | message | 補足 |
|:---:|:---:|:---:|
| 400 | bad request | 不正なJSON |
| 401 | unauthorized | 認証エラー |
| 500 | internal server error | 不明な内部エラー |

## GET /user
### 概要
user情報を取得する
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"userid",
    "name":"username",
    "email":"example@example.com"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | user not found | userが存在しない |

## POST /user
### 概要
新規userを作成する
### 認証
必要なし
### リクエスト
```
{
    "name":"username",
    "password":"password",
    "email":"example@example.com"
}
```
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"userid",
    "name":"username",
    "email":"example@example.com"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 400 | email already exists | 同じemailのユーザーが既に存在 |

## PUT /user
### 概要
user情報を更新する
### 認証
必要あり
### リクエスト
```
{
    "name":"username",
    "password":"password",
    "email":"example@example.com"
}
```
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"userid",
    "name":"username",
    "email":"example@example.com"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | user not found | userが存在しない |
| 400 | email already exists | 同じemailのユーザーが既に存在 |

## DELETE /user
### 概要
userを削除する
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
空
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | user not found | userが存在しない |

## GET /task/:id
### 概要
taskを取得する
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| id | taskのid |
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"taskid",
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | task not found | taskが存在しない |

## CREATE /task
### 概要
新規taskを作成する
### 認証
必要あり
### リクエスト
```
{
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"taskid",
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|

## PUT /task/:id
### 概要
task情報を更新する
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| id | taskのid |
### 認証
必要あり
### リクエスト
```
{
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"taskid",
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | task not found | taskが存在しない |

## DELETE /task/:id
### 概要
taskを削除する
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| id | taskのid |
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
空
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|
| 404 | task not found | taskが存在しない |

--------------------------------------------------------------------------------
**以下は未実装**

## PUT /task/:id/comp
### 概要
taskのcompletedを切り替える
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| id | taskのid |
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "id":"taskid",
    "title":"taskname",
    "content":"I am content.",
    "iscomp":false,
    "date":"2020-12-06"
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|

## GET /task/date/:date
### 概要
特定の日付のtaskを取得する
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| date | 指定する日付(yyyy-mm-dd) |
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "tasks": [
    {
        "id":"taskid1",
        "title":"taskname1",
        "content":"I am content1.",
        "iscomp":false,
        "date":"2020-12-06",
    },
    {
        "id":"taskid2",
        "title":"taskname2",
        "content":"I am content2.",
        "iscomp":true,
        "date":"2020-12-06",
    }
    ]
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|

## GET /task/date/from/:start/to/:end
### 概要
特定の期間のtaskを取得する
### パスパラメータ
| key | 説明 |
|:---:|:---:|
| start | 期間の開始の日付(yyyy-mm-dd) |
| end | 期間の終了の日付(yyyy-mm-dd) |
### 認証
必要あり
### リクエスト
なし
### レスポンス
| code | 補足 |
|:---:|:---:|
| 200 | |
```
{
    "tasks": [
    {
        "id":"taskid1",
        "title":"taskname1",
        "content":"I am content1.",
        "iscomp":false,
        "date":"2020-12-06",
    },
    {
        "id":"taskid2",
        "title":"taskname2",
        "content":"I am content2.",
        "iscomp":true,
        "date":"2020-12-07",
    }
    ]
}
```
### エラー
| code | message | 補足 |
|:---:|:---:|:---:|

