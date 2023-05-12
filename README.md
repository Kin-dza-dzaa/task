### Task for go dev.

### Summary

---

API слушает по `http://localhost:8000`

Хостится swagger-ui `http://localhost:8000/swagger/`

Есть unit тесты, для работы тестов нужен рабочий докер с выключенным TLS на порту 2375: 

```plaintext
make test
make cover
```

### Usage

---

```plaintext
docker compose build
docker compose up -d
```
