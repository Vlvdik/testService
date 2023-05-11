---
# Task #5

---

### migrations

download ```https://github.com/golang-migrate/migrate```

**postgres**:
```
migrate -path ./schema/pg -database "your_path" up
```
**clickhouse**:
```
migrate -path ./schema/clickhouse "your_path" up
```
---

**installation**
1) Clone repository
```
git clone https://github.com/Vlvdik/hezzlService
```
2) Update go.mod
```
go mod tidy
```
3. Setup yours config.yml file (values are given as an example)
~~~
server:
  host: "localhost"
  port: ":8080"
  request_timeout: 30
db:
  port: "5432"
  name: "postgres"
  user: "postgres"
  pwd: ""
cache:
  host: "localhost"
  port: ":6379"
  pwd: ""
  secret: ""
  write_duration: 60
broker:
  url: "nats://0.0.0.0:49652"
  user: ""
  pwd: ""
  subject: ""
  max_pending: 256
clickhouse:
  host: ""
  port: ""
  db: ""
  user: ""
  pwd: ""
  table: ""
  timeout: 30
~~~
4. run application
```
go run cmd/app/main.go
```

### Unfortunately, I was not able to configure docker-compose by the deadline, so all components (Postgres, Redis, Nats, Clickhouse) must be run separately before running the main application

---

#### What I would like to improve:
1) Add swagger
2) Add docker-compose configuration
3) Create a docker container for core-api
4) Add tests
