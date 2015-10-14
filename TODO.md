# todo list

* [x] 支援 */5 表示法
* [x] 支援 1-15 表示法
* [x] 計數處理
* [x] 排程備份還原
* [x] 處理 process commands 迭代問題
* http handler
    - [x] process.Revoke
    - [x] process.Backup(f string)
    - [x] process.Recover(f string)
* [x] 在不開 http 服務下,支援讀檔啟動
* [ ] 指定備份間隔
* [x] 排程設定檔允許使用#當註解
* [x] 開放 schedule 格式可指定 "週"
* [x] 搬入 parse 與 spec test


### 下階段重構

* [ ] 開放api, 但僅允許 http hook 註冊
* [ ] 加入提示 api server 被開啟的訊息


### 未來重構方向

* 參考 [robfig/cron](https://github.com/robfig/cron), 評估優化方向
    
    1. 預計開發另一個版本, 直接引入 robfig/cron 僅包裝注入的匿名函式


### 評比效能

* 連續跑 n 條排程 m 天, 分析是否有效能問題

