# Docker i Go pokretanje aplikacije od nule

Ovo su koraci za scenario kao da je projekat tek preuzet sa GitHub-a. Setup podize PostgreSQL bazu u Dockeru, priprema Go dependency-je i pokrece aplikaciju.

## Preduslovi

Instaliraj:

- Docker Desktop ili Docker Engine sa Docker Compose podrskom
- Go verziju koja odgovara projektu, prema `go.mod`: `1.26.4`
- Git

Provera instalacije:

```bash
docker --version
docker compose version
go version
git --version
```

## Kloniranje projekta

```bash
git clone <URL_REPOZITORIJUMA>
cd toctou-demo
```

Sve naredne komande pokrecu se iz korena projekta, gde se nalaze `go.mod`, `main.go`, `docker-compose.yml` i `app/Dockerfile`.

## Priprema Go okruzenja

Go dependency-ji su definisani u `go.mod` i `go.sum`. Preuzimaju se komandom:

```bash
go mod download
```

Ako zelis da proveris da li se aplikacija lokalno kompajlira:

```bash
go build .
```

Ovo proverava Go okruzenje na host masini. Za Docker pokretanje aplikacije nije neophodno da imas lokalno preuzete dependency-je, jer ih Dockerfile preuzima unutar build kontejnera, ali je korisno za razvoj i pokretanje `go run .`.

## Pokretanje preko Docker Compose-a

Najjednostavniji nacin je da Docker Compose podigne i bazu i aplikaciju:

```bash
docker compose up --build
```

Ova komanda radi sledece:

- pokrece PostgreSQL kontejner `toctou_demo_db`
- kreira bazu `toctou`
- izvrsava inicijalni SQL iz `app/db/schema.sql`
- build-uje Go aplikaciju preko `app/Dockerfile`
- pokrece aplikacioni kontejner `toctou_demo_app`
- izlozi aplikaciju na `http://localhost:8080`
- izlozi PostgreSQL na host portu `5433`

Compose koristi `postgres:16-alpine`, da bi verzija baze bila stabilna i da se ne bi slucajno promenila pri sledecem pokretanju.

Kada se u logu aplikacije pojavi poruka da je server pokrenut, otvori:

```text
http://localhost:8080/login
```

## Login podaci

Admin nalog:

```text
email: lobradovic@mail.com
password: password1
```

Korisnicki nalozi:

```text
email: lmihailovic@mail.com
password: password2

email: vlazarevic@mail.com
password: password3
```

## Korisne Docker komande

Pokretanje u pozadini:

```bash
docker compose up --build -d
```

Prikaz logova:

```bash
docker compose logs -f
```

Zaustavljanje kontejnera bez brisanja baze:

```bash
docker compose down
```

Potpuno resetovanje baze i ponovno izvrsavanje `app/db/schema.sql`:

```bash
docker compose down -v
docker compose up --build
```

Ulazak u PostgreSQL shell:

```bash
docker compose exec db psql -U postgres -d toctou
```

Primer provere podataka u bazi:

```sql
SELECT id, email, balance FROM users ORDER BY id;
SELECT id, sender_id, recipient_id, amount, timestamp FROM transfers ORDER BY id;
```

## Lokalno pokretanje Go aplikacije uz Docker bazu

Ako hoces da bazu drzis u Dockeru, a aplikaciju pokrenes lokalno kroz Go, podigni samo bazu:

```bash
docker compose up -d db
```

Zatim pokreni aplikaciju lokalno:

```bash
DATABASE_URL='host=localhost port=5433 user=postgres password=admin dbname=toctou sslmode=disable' go run .
```

Aplikacija ce biti dostupna na:

```text
http://localhost:8080/login
```

Napomena: ako `DATABASE_URL` nije postavljen, aplikacija vec ima podrazumevani connection string za lokalni Postgres na `localhost:5433`, tako da moze da radi i samo:

```bash
go run .
```

## Sta radi Dockerfile

Dockerfile za aplikaciju nalazi se u `app/Dockerfile` i koristi multi-stage build.

Prva faza koristi Go image i:

- postavlja radni direktorijum na `/src`
- kopira `go.mod` i `go.sum`
- preuzima Go dependency-je komandom `go mod download`
- kopira `main.go`, `handlers/`, `middleware/`, `templates/` i `styles/`
- kompajlira aplikaciju u Linux binarni fajl `/bin/toctou-demo`

Druga faza koristi manji Alpine image i:

- postavlja radni direktorijum na `/app`
- kopira gotov binarni fajl
- kopira `templates/` i `styles/`, jer ih aplikacija cita u runtime-u
- otvara port `8080`
- pokrece aplikaciju komandom `./toctou-demo`
