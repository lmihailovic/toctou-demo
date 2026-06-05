# Docker pokretanje aplikacije

Dockerfile za aplikaciju se nalazi u `app/Dockerfile`. On koristi multi-stage build, sto znaci da se aplikacija prvo kompajlira u posebnom Go image-u, a zatim se u finalni image kopira samo gotov binarni fajl i fajlovi koji su potrebni za rad aplikacije.

## Sta radi Dockerfile

Prva faza koristi image `golang:1.26-alpine` i sluzi za build aplikacije:

- postavlja radni direktorijum na `/src`
- kopira `go.mod` i `go.sum`
- preuzima Go dependency-je komandom `go mod download`
- kopira izvorni kod aplikacije: `main.go`, `handlers/`, `middleware/`
- kopira staticke i template fajlove: `templates/` i `styles/`
- kompajlira aplikaciju u Linux binarni fajl `/bin/toctou-demo`

Druga faza koristi manji image `alpine:3.22` i sluzi za pokretanje aplikacije:

- postavlja radni direktorijum na `/app`
- kopira kompajliranu aplikaciju iz prve faze
- kopira `templates/` i `styles/`, jer ih aplikacija cita u runtime-u
- otvara port `8080`
- pokrece aplikaciju komandom `./toctou-demo`

## Build Docker image-a

Komanda se pokrece iz korena projekta:

```bash
docker build -f app/Dockerfile -t toctou-demo-app .
```

Opcije u komandi:

- `-f app/Dockerfile` govori Docker-u gde se nalazi Dockerfile
- `-t toctou-demo-app` daje ime image-u
- `.` znaci da je build context koren projekta

Build context mora biti koren projekta zato sto Dockerfile kopira fajlove kao sto su `go.mod`, `main.go`, `handlers/`, `templates/` i `styles/`.

## Pokretanje baze

Aplikacija koristi PostgreSQL bazu. Ako Postgres kontejner jos nije napravljen, moze da se pokrene ovako:

```bash
docker run --name toctou_demo \
  -e POSTGRES_PASSWORD=admin \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_DB=toctou \
  -p 5433:5432 \
  -d postgres
```

Ako kontejner vec postoji, ali je zaustavljen, pokrece se ovako:

```bash
docker start toctou_demo
```

## Pokretanje aplikacije

Kada je baza pokrenuta, aplikacija se pokrece komandom:

```bash
docker run --rm -p 8080:8080 \
  -e 'DATABASE_URL=host=host.docker.internal port=5433 user=postgres password=admin dbname=toctou sslmode=disable' \
  toctou-demo-app
```

Opcije u komandi:

- `--rm` brise aplikacioni kontejner kada se zaustavi
- `-p 8080:8080` mapira port aplikacije na host masinu
- `-e DATABASE_URL=...` prosledjuje connection string za PostgreSQL bazu
- `host=host.docker.internal` omogucava aplikacionom kontejneru da pristupi Postgresu koji je mapiran na host port `5433`

Nakon pokretanja, aplikacija je dostupna na:

```text
http://localhost:8080
```

Login stranica je dostupna na:

```text
http://localhost:8080/login
```
