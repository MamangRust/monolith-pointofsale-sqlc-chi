# Migrasi Apigateway: Echo → Chi

> **Lingkup:** `monolith-pointofsale-grpc/service/apigateway` (+ harness `tests/` dan `pkg/upload_image`)
> **Router:** `github.com/labstack/echo/v4` → `github.com/go-chi/chi/v5`
> **Status:** ✅ Selesai — build ✅ · vet ✅ · unit test ✅
> **Pola:** sama dengan migrasi apigateway `monolith-payment-gateway-grpc` dan `monolith-ecommerce-grpc` (lihat `MIGRATION_ECHO_TO_CHI.md` masing-masing).

---

## 1. Ringkasan

| Item | Jumlah |
|---|---|
| File yang import echo (sebelum) | 16 |
| Handler domain yang dimigrasi | 11 (`auth`, `cashier`, `category`, `merchant`, `merchant_document`, `order`, `order_item`, `product`, `role`, `transaction`, `user`) — struktur flat, satu paket |
| Rute terdaftar (terkonversi) | 169 (99 GET · 60 POST · 10 DELETE) |
| Middleware global | 8 |
| Middleware per-route | — (semua rute dibungkus `apiHandler.Handle`, tanpa chain) |
| Test harness `tests/` yang dikonversi | 13 file |

Pos ini adalah yang paling sederhana dari ketiganya: handler flat tanpa sub-paket, registrasi rute langsung `router.Group(...)` + `apiHandler.Handle(...)`.

## 2. Karakteristik Khusus POS

1. **Error semantics gaya payment-gateway** — `HTTPErrorHandler` generik (`{"status":"error","message":...,"code":...}`), `shared.HandleApiError` tanpa penanganan `*echo.HTTPError`. Adapter `httpx.Handler`/`apierror.HandleApiError` yang dipakai adalah versi payment-gateway.
2. **Context key `userID`** (bukan `userId`/`user_id`) — `JWTAuth()` custom menyimpan claim `sub` apa adanya ke context key `"userID"`, sama seperti `SuccessHandler` echo-jwt lama.
3. **`TraceContextMiddleware`** — sekadar ekstraksi W3C trace context, di-rewrite ke chi tanpa perlu `wrapResponseWriter` (tidak membaca status).
4. **Upload gambar** — `pkg/upload_image.ProcessImageUpload(c echo.Context, file)` → `(w http.ResponseWriter, file)` dengan helper `writeJSONError`; dipakai handler `product` (multipart `c.FormValue`/`c.FormFile` → `r.FormValue`/`r.FormFile`, arity 3).
5. **Helper `parseQueryInt`/`parseQueryStringRequired`/`parseQueryIntWithValidation`** — signature `echo.Context` → `*http.Request`.
6. **RateLimiter & TraceContext** — tidak terpasang di chain global (dead code di bootstrap), tetap dikonversi agar konsisten.

## 3. Dependensi

**Ditambah:** `go-chi/chi/v5`, `go-chi/cors`, `swaggo/http-swagger`.
**Tidak lagi di-import langsung:** `labstack/echo/v4`, `echo-jwt/v4`, `echo-swagger` (masih *indirect* via `shared/errors` echo-typed yang dipakai modul `tests/`).

## 4. Struktur Baru

```
service/apigateway/
├── apierror/          # NEW — ApiHandler versi net/http
├── httpx/             # NEW — helper JSON/Bind/RealIP + request-scoped values
├── apps/client.go     # bootstrap chi (Recoverer → RequestID → Logger → Trace → CORS → Compress → Secure → JWT)
├── handler/           # 13 file flat → chi.Router
└── middlewares/       # auth (JWT), ratemiddleware, trace
pkg/upload_image/      # ProcessImageUpload(w http.ResponseWriter, file)
```

## 5. Verifikasi

```bash
cd monolith-pointofsale-grpc/service/apigateway
go build ./...    # ✅
go vet ./...      # ✅
cd ../../shared && go build ./...   # ✅ tidak berubah
cd ../tests && go vet ./...         # ✅ harness terkonversi
```

Harness `tests/` (13 file): `NewHandlerX(s.echo, …)` → `NewHandlerX(s.echo /* chi.Router */, …)`, bypass middleware `c.Set("userId", …)` → `context.WithValue`, `MockImageUpload` ikut signature baru, `apierror.NewApiHandler` menggantikan `errors.NewApiHandler` shared.

## 6. Catatan Perilaku

1. **Respons error identik** — 404/405 JSON eksplisit `{"status":"error",...}` via `r.NotFound`/`r.MethodNotAllowed` (bentuk sama dengan error handler echo lama), error `AppError` dari rute wrapped tetap JSON `ErrorResponse` dengan trace_id dari span OTel.
2. **Echo masih indirect dependency** via `shared/errors`.

## 7. Follow-up (Opsional)

- [ ] Smoke test `hurl/` atau full stack docker.
- [ ] Migrasi `tests/` + `shared/errors` penuh ke net/http → echo hilang total dari ketiga repo.
- [ ] Ketiga apigateway kini konsisten memakai chi; pertimbangkan menyamakan lokasi package fondasi (`httpx`/`apierror`) dan versi chi lintas repo.
