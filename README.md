# product-service

Product service adalah rest api basic untuk handle data product, dibangun menggunakan **GO**, **PostgreSQL**,**Redis**,dengan **Hexagonal Arsitektur**.

## Arsitektur
### Hexagonal Arsitektur

Project ini menggunakan hexagonal untuk memisahkan bisnis logic dari infrastruktur

### Alasan mengunakan Hexagonal
- Saya ingin memisahkan bisnis logic utama aplikasi dengan infrastruktur baik dalam pemilihan httphandler,cache  dan database yang di pakai.
- Lebih Clear dan mudah untuk di test ,karna bisnis logic uda dipisah
- Flexibility,Arsitektur memudahkan untuk ganti teknology karna cukup ganti implement tanpa harus ngerubah bisnis logic
- Menurut saya struktur lebih rapi , gampang untuk dimaintenance dan dikembangkan


## Tech Stack
- Bahasa: Go
- Database: PostgreSQL
- Cache: Redis
- Router: Gorilla Mux
- Authentication: HTTP Basic Auth
- Testing: `testing`, `sqlmock`, `mockery`
- Containerization: Docker & Docker Compose


## For Testing
```bash
docker compose up --build
```
### Service yang Up:
- product-service(http://localhost:8181) 
- Redis
- PostgreSQL
- migration database

Unit Test di jalakan waktu docker build

### Example API

- Add Product
curl -u admin:admin123 \
  -X POST http://localhost:8181/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "product 1",
    "price": 2000,
    "description": "desc product 1",
    "quantity": 10
  }'

- default product list
curl -u admin:admin123 http://localhost:8181/products

- short price products
curl -u admin:admin123 "http://localhost:8181/products?sort=price_asc"

- sort Name A-Z 
curl -u admin:admin123 \
  "http://localhost:8181/products?sort=name_asc"








