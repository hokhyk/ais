# AIS
應收應付系統API

## 資料夾結構
./controller 控制器

./libs 函式庫

./service 服務器

./.env.swp 設定置換檔

./.go.mod Golang Module 設定檔

./main.go 主程式

## 資料流
User -> Borwser -> HTTP Request -> main.go -> 分配到相應之控制器 -> 初步檢查 HTTP Request 參數都有帶齊 -> 傳給服務器 -> 執行商業邏輯 -> 調用運行商業邏輯所需之函式庫 -> 回傳結果 -> 控制器 -> 輸出 HTTP Response -> Borwser -> User

## 使用說明
1. 將 .env.swp 置換成 .env
2. .env內參數對應如下
```
port=[本系統運行的Port]
hrm.url=[HRM 系統網址]
resis.addr=[Redis 位址]
redis.password=[Redis 密碼]
redis.db=[Redis 資料庫名稱]
mysql.host=[MySQL 主機位置]
mysql.user=[MySQL 帳號]
mysql.password=[MySQL 密碼]
mysql.database=[MySQL 資料庫]
mysql.charset=[MySQL 預設編碼]
mysql.parsetime=[MySQL 使用本機時間]
mysql.loc=[MySQL 預設時區]
salt=[鹽值，供SHA256加密用]
```
3. 此系統架構於Jorani這套HRM系統，所有的員工資料都由HRM調用，使用其他HRM系統與LDAP，請自行撰寫API，與修改使用者驗証相關

## Todo
1. 管理員檢閱當月請款/請購細項

## API
只要回傳的 code 不是 200 ，通通都是有問題， message 會說明是什麼原因

登入
```
# HTTP POST
url: http://[Path To]/users/login

# 參數 
login=[HRM使用者帳號]
password=[HRM使用者密碼]

# 回傳 
{
    code: 1
    message: [令牌]
}
```

驗証使用者是否登入
```
# HTTP POST
url: http://[Path To]/users/checkLogin

# 參數
token=[令牌]

# 回傳
{
    code: 1
    message: "true"
}
```

登出
```
# HTTP GET
url: http://[Path To]/users/logout

# 參數
token=[令牌]

# 回傳
{
    code: 1
    message: "true"
}
```

取得使用者資訊
```
# HTTP POST
url: http://[Path To]/users/getUser

# 參數
token=[令牌]

# 回傳
{
    code: 1
    message: [使用者資訊JSON]
}
```

新增訂購單
```
# HTTP POST
url: http://[Path To]/pr/add

# HTTP Header
Content-Type: multipart/form-data

# 參數
token=[令牌]&pay_to=[支付對象]&vendor_name=[廠商名稱]&pay_type=[入帳類別]&list_type=[類別]pay_method=[支付方式]&bank_account=[銀行帳號]&sign_at=[日期]&name[]=[項目]&currency[]=[幣種]&unit_price[]=[單價]&quantity[]=[數量]&exchange_rate[]=[匯率]&tax[]=[稅額]&proof[]=[佐証資料]

# 參數說明
1. 支付對象，廠商與請款人二擇一，勾選廠商需填入廠商名稱
2. 廠商名稱，如勾選支付對象為請款人時，資料可以為空字串
3. 入帳類別，給數字代號，1=庶務消耗性商品，2=業務轉售類商品，3=設備類固定資產，4=原物料類商品，5=其他
4. 類別，給數字代號，1=請購單，2=請款單
5. 支付方式，1=支票，2=現金，3=匯款，4=零用金
6. 銀行帳號，如勾選支付方式為匯款時，需填入匯款帳號，否則給空字串就好
7. 日期，需給以下日期格式之字串YYYY/MM/DD
8. 項目，可以同時多筆，名字著取一樣叫name就好
9. 幣種，可以同時多筆，名字都取一樣叫currency就好，直接給UI上文字的中文就好
10. 單價，可以同時多筆，名字都取一樣叫unit_price就好，需要可以有小數點
11. 數量，可以同時多筆，名字都取一樣叫quantity就好
12. 匯率，可以同時多筆，名字都取一樣叫exchange_rate就好，請給小數點，好比1:29.5，就是29.5
13. 稅額，可以同時多筆，名字都取一樣叫tax就好，請給小數點，好比5%就是0.05
14. 佐証資料，可以同時多筆，名字著取一個proof就好，到了Server會通通壓成一個zip檔

# 回傳
{
    code: 1
    message: "true"
}
```

作廢訂購單
```
# HTTP POST
url: http://[Path To]/pr/setCancel

# 參數
token=[令牌]&id=[請購單號]

# 回傳
{
    code: 1
    message: "true"
}
```

下載佐証檔
```
# HTTP POST
url: http://[Path To]/download/getFile?proof=[佐証檔名]

# 參數
token=[令牌]

# 參數說明
1. 佐証檔名，透過訂購單列表可取得佐証檔名與整段下載網址

# 回傳
[下載檔案]
```