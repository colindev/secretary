# todo list

* [ ] 支援 */5 表示法
* [ ] 支援 1-15 表示法
* [x] 計數處理
* [x] 排程備份還原
* [x] 處理 process commands 迭代問題
* http handler
    - [x] process.Revoke
    - [x] process.Backup(f string)
    - [x] process.Recover(f string)
* [x] 在不開 http 服務下,支援讀檔啟動
* [ ] 指定備份間隔
* [ ] 排程設定檔允許使用#當註解


### 未來重構方向
* 參考 [robfig/cron](https://github.com/robfig/cron), 評估優化方向

