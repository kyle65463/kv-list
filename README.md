# kv-list

## Setup

新增環境變數檔案 `.env`：
```
cp .env.example .env
```
更新 `.env` 中的環境變數，例如 `DATABASE_URL`。


執行 migration：
```
goose -dir migrations postgres <DATABASE_URL> up
```

## Commands

執行 api server：
```
go run .
```

執行 unit tests：
```
go test ./...
```
Notes: 因為 base requirements 只有 Get Head / Get Page API，因此沒有撰寫其他額外 api 的 unit tests。

## DB

### DBMS
DBMS 使用 Postgres，因為在考慮高流量的情況下，Postgres 的效能會較好一些（https://itnext.io/benchmark-databases-in-docker-mysql-postgresql-sql-server-7b129368eed7 ）。並且 Postgres 更嚴格的符合了 ACID，能保證資料的完整性。而在開發上，Postgres 的資料型別相對多元且彈性，像是可以儲存 array、json（雖然沒用到），因此我認為使用 Postgres 會比較合適。


### Schema
有兩個 table，pages 與 lists。假設每個 page 裡面儲存的內容都是 `BYTEA` 型別。
```sql
CREATE TABLE IF NOT EXISTS pages (
    key TEXT PRIMARY KEY,
    data BYTEA NOT NULL,
    next_page_key TEXT REFERENCES pages(key) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS pages_created_at ON pages(created_at); -- Used for deleting pages

CREATE TABLE IF NOT EXISTS lists (
    key TEXT PRIMARY KEY,
    next_page_key TEXT REFERENCES pages(key) ON DELETE SET NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```


## API

### Get list head

```go
GET /api/v1/lists/:key

Response:
type PageResponse struct {
	Key         string  `json:"key"`
	Data        []byte  `json:"data"`
	NextPageKey *string `json:"nextPageKey"`
}
```

取得 key 對應的 list head。

### Get page

```go
GET /api/v1/pages/:key

Response:
type ListHeadResponse struct {
	Key         string  `json:"key"`
	NextPageKey *string `json:"nextPageKey"`
}
```

取得 key 對應的 page。

### Set list

```go
POST /api/v1/lists/:key

Request Body:
type body struct {
	Data [][]byte `json:"data"`
}
```

為了方便使用，可以一次 set 列表的所有 page。Request body 的 data 中的每個元素皆對應到一個 page 的儲存內容。

若該列表原本已存在，會直接覆蓋掉成要 set 的列表。

Nice to have: 可以多一個 append pages 的 api，會直接把 pages 加在原本列表的後面，而不是直接覆蓋掉整個列表。

### Delete pages

```
DELETE /api/v1/pages?interval=<interval>&limit=<limit>
```

刪除一段時間之前創建的 pages。透過 interval 參數指定時間，例如 `DELETE /api/v1/pages?interval=1d` 會刪除一天以前新增的所有 pages。若要刪除的檔案量很大，可以透過 limit 參數分批進行刪除。

Notes: 一定要傳入 interval 參數，limit 參數為 optional。

其他可能的刪除方式：在 page 上多紀錄它屬於的列表的 key，每次刪除的時候以列表為單位進行刪除（`DELETE FROM pages WHERE list_key < <list_key>`）。

## Helper scripts

提供了兩個 scripts，`set_list.py` 與 `get_list.py`，方便 debug 與測試。在執行前須先執行 server。

### set_list.py
```
python scripts/set_list.py <list_key> [base_endpoint_url]
```
可以透過 set list api 去新增 / 更新一個列表，列表的內容可以在 `set_list.py` 中定義。若沒指定 `base_endpoint_url`，預設則為 `http://localhost:8080`。

### get_list.py
```
python scripts/get_list.py <list_key> [base_endpoint_url]
```
遞迴取得一個列表的所有頁面。若沒指定 `base_endpoint_url`，預設則為 `http://localhost:8080`。