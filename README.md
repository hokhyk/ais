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
只要回傳的 status 不是 1 ，通通都是有問題， message 會說明是什麼原因

登入
```
# HTTP POST
url: http://[Path To]/users/login

# 參數 
login=[HRM使用者帳號]
password=[HRM使用者密碼]

# 回傳 
{
    status: 1
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
    status: 1
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
    status: 1
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
    status: 1
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
token=[令牌]&pay_to=[支付對象]&vendor_name=[廠商名稱]&pay_type=[入帳類別]&list_type=[類別]pay_method=[支付方式]&bank_account=[銀行帳號]&sign_at=[日期]&pay_date=[付款日]&pr_item=[應付項目]&serial[對應單號]&installment_plan=[分多少期]&pay_by=[第幾期]&memo=[備註]&name[]=[項目]&currency[]=[幣種]&unit_price[]=[單價]&quantity[]=[數量]&exchange_rate[]=[匯率]&tax[]=[稅額]&proof[]=[佐証資料]

# 參數說明
1. 支付對象，廠商與請款人二擇一，勾選廠商需填入廠商名稱
2. 廠商名稱，如勾選支付對象為請款人時，資料可以為空字串
3. 入帳類別，數字，資料從[取得入帳類別API]獲得
4. 類別，給數字代號，1=請購單，2=請款單
5. 支付方式，數字，資料從[取得支付方式API]獲得
6. 銀行帳號，如勾選支付方式為匯款時，需填入匯款帳號，否則給空字串就好
7. 日期，需給以下日期格式之字串YYYY/MM/DD
8. 付款日，需給以下日期格式之字串YYYY/MM/DD
9. 應付項目，數字，資料從[取得應付項目API]獲得
10. 對應單號，可不填，除非有指定分期設定
11. 分多少期，可不填，除非有指定分期設定
12. 第幾期，可不填，除非有指定分期設定
13. 項目，可以同時多筆，名字著取一樣叫name就好
14. 幣種，可以同時多筆，名字都取一樣叫currency就好，資料從[取得幣別API]獲得
15. 單價，可以同時多筆，名字都取一樣叫unit_price就好，需要可以有小數點
16. 數量，可以同時多筆，名字都取一樣叫quantity就好
17. 匯率，可以同時多筆，名字都取一樣叫exchange_rate就好，請給小數點，好比1:29.5，就是29.5
18. 稅額，可以同時多筆，名字都取一樣叫tax就好，請給小數點，好比5%就是0.05
19. 佐証資料，可以同時多筆，名字著取一個proof就好，到了Server會通通壓成一個zip檔

# 回傳
{
    status: 1
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
    status: 1
    message: "true"
}
```

取得請購單列表
```
# HTTP POST
url: http://[Path To]/pr/getList

# 參數
token=[令牌]&page=[第..頁]&num=[每頁筆數]&begin[起始日期]&end=[結束日期]

# 參數說明
1. 第..頁，可不給，預設從第1頁開始
2. 每頁筆數，可不給，預設一頁十筆
3. 起始日期，可不給，有給會依簽單日期範圍做抓取，格式YYYY/MM/DD
4. 結束日期，可不給，有給會依簽單日期範圍做抓取，格式YYYY/MM/DD

# 回傳
{
    status: 1,
    list: [{
        id: [ID],
        name: [名稱],
        pr_list_id: [請購單ID],
        currency: [幣別],
        unit_price: [單價],
        quantity: [數量],
        exchange_rate: [匯率],
        tax: [稅額],
        total_price: [總價]
    },{...}]
}
```

取得請購單資訊
```
# HTTP POST
url: http://[Path To]/pr/getItem

# 參數
token=[令牌]&id=[請購單ID]

# 回傳
{
    status: 1
    list: {
        id: [ID],
        pay_to: [支付對象],
        vendor_name: [廠商名稱],
        pay_type: [入帳類別],
        list_type: [類別],
        users_id: [員工ID],
        email: [員工信箱],
        identifier: [員工編號],
        lastname: [姓名],
        firstname: [名字],
        pay_method: [支付方式],
        bank_account: [銀行帳號],
        proof: [佐証資料],
        status: [請購單狀態],
        sign_at: [簽單日期],
        pay_date: [付款日],
        payMethod: [應付項目],
        serial: [對應單號],
        installment_plan: [分多少期],
        pay_by: [第幾期],
        memo: [備註],
        create_at: [新增日期]
    },
    detail: [{
        id: [ID],
        name: [名稱],
        pr_list_id: [請購單ID],
        currency: [幣別],
        unit_price: [單價],
        quantity: [數量],
        exchange_rate: [匯率],
        tax: [稅額],
        total_price: [總價]
    }, {...}]
}
```

管理員-取得請購單列表
```
# HTTP POST
url: http://[Path To]/admin/getList

# 參數
token=[令牌]&page=[第..頁]&num=[每頁筆數]&begin[起始日期]&end=[結束日期]

# 參數說明
1. 第..頁，可不給，預設從第1頁開始
2. 每頁筆數，可不給，預設一頁十筆
3. 起始日期，可不給，有給會依簽單日期範圍做抓取，格式YYYY/MM/DD
4. 結束日期，可不給，有給會依簽單日期範圍做抓取，格式YYYY/MM/DD

# 回傳
{
    status: 1,
    list: [{
        id: [ID],
        name: [名稱],
        pr_list_id: [請購單ID],
        currency: [幣別],
        unit_price: [單價],
        quantity: [數量],
        exchange_rate: [匯率],
        tax: [稅額],
        total_price: [總價]
    },{...}]
}
```

管理員-取得請購單資訊
```
# HTTP POST
url: http://[Path To]/admin/getItem

# 參數
token=[令牌]&id=[請購單ID]

# 回傳
{
    status: 1
    list: {
        id: [ID],
        organization_id: [部門ID],
        organization_name: [部門名稱],
        pay_to: [支付對象],
        vendor_name: [廠商名稱],
        pay_type: [入帳類別],
        list_type: [類別],
        users_id: [員工ID],
        email: [員工信箱],
        identifier: [員工編號],
        lastname: [姓名],
        firstname: [名字],
        pay_method: [支付方式],
        bank_account: [銀行帳號],
        proof: [佐証資料],
        status: [請購單狀態],
        sign_at: [簽單日期],
        create_at: [新增日期]
    },
    detail: [{
        id: [ID],
        name: [名稱],
        pr_list_id: [請購單ID],
        currency: [幣別],
        unit_price: [單價],
        quantity: [數量],
        exchange_rate: [匯率],
        tax: [稅額],
        total_price: [總價]
    }, {...}]
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

取得幣別
```
# HTTP POST
url: http://[Path To]/currency/getCurrency

# 參數說明
將key給HTML項目做顯示用，將id當做value回傳

# 回傳
[
    {
        "id":[ID],
        "key":[鍵]
    }, { ... }
]
```

取得應付項目
```
# HTTP POST
url: http://[Path To]/prItem/getPrItem

# 參數說明
將key給HTML項目做顯示用，將id當做value回傳

# 回傳
[
    {
        "id":[ID],
        "key":[鍵]
    }, { ... }
]
```

取得支付方式
```
# HTTP POST
url: http://[Path To]/payMethod/getPayMethod

# 參數說明
將key給HTML項目做顯示用，將id當做value回傳

# 回傳
[
    {
        "id":[ID],
        "key":[鍵]
    }, { ... }
]
```

取得入帳類別
```
# HTTP POST
url: http://[Path To]/creditType/getCreditType

# 參數說明
將key給HTML項目做顯示用，將id當做value回傳

# 回傳
[
    {
        "id":[ID],
        "key":[鍵]
    }, { ... }
]
```

取得廠商列表
```
# HTTP POST
url: http://[Path To]/company/getCompany

# 參數說明
將key給HTML項目做顯示用，將id當做value回傳

# 回傳
[
    {
        "id":[ID],
        "key":[鍵]
    }, { ... }
]
```