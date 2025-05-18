
```curl
curl --location 'localhost:8080/competitions' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'start_at=1747523554' \
--data-urlencode 'end_at=1747566761' \
--data-urlencode 'rules=event_type==bet ? amount : 0'
```


```
 .schema bets
```
