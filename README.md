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
hrm.url=[HRM 系統網址]
resis.addr=[Redis 位址]
redis.password=[Redis 密碼]
redis.db=[Redis 資料庫名稱]
salt=[鹽值，供SHA256加密用]
```
3. 此系統架構於Jorani這套HRM系統，所有的員工資料都由HRM調用，使用其他HRM系統與LDAP，請自行撰寫API，與修改使用者驗証相關

## Todo
1. ~~使用者驗証~~
2. 請款/請購單上傳
3. 請款/請購上級主管
4. 請款/請購作廢
5. 管理員檢閱當月請款/請購細項
