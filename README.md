# MentorHub KBTU — Security Module

Часть проекта, отвечающая за информационную безопасность: аутентификация, авторизация, защита API, защита персональных данных.

## Установка

```bash
npm install jsonwebtoken bcrypt joi xss-clean helmet express-rate-limit cors \
  express-mongo-sanitize passport passport-google-oauth20 winston speakeasy qrcode
```

Скопируй `.env.example` в `.env` и заполни значениями:

```bash
cp .env.example .env
```

## Структура

| Файл | Что делает |
|---|---|
| `utils/password.js` | Хеширование паролей (bcrypt) |
| `middleware/auth.js` | Генерация и проверка JWT |
| `middleware/oauth.js` | Вход через Google OAuth (@kbtu.kz) |
| `middleware/rbac.js` | Контроль доступа по ролям |
| `middleware/security.js` | helmet, rate limiting, XSS/инъекции, CORS |
| `middleware/https.js` | Принудительный HTTPS + конфиг nginx |
| `middleware/passwordReset.js` | Сброс пароля через одноразовую ссылку |
| `middleware/twoFactorAuth.js` | 2FA для админов (TOTP / Google Authenticator) |
| `validation/schemas.js` | Валидация входных данных (Joi) |
| `logging/securityLogger.js` | Логирование подозрительной активности |
| `demo/attack-simulation.js` | Демо: доказательство защиты от SQLi/XSS |

## Как подключить (для Backend Developer)

```javascript
const express = require('express');
const app = express();
require('dotenv').config();

const { applySecurity } = require('./middleware/security');
const { enforceHttps } = require('./middleware/https');
const { verifyToken } = require('./middleware/auth');
const { requireRole } = require('./middleware/rbac');

app.use(enforceHttps);
applySecurity(app);
app.use(express.json({ limit: '10kb' }));

// Пример защищённого роута
app.get('/api/admin/users', verifyToken, requireRole('admin'), (req, res) => {
  res.json({ message: 'Только для админов' });
});

app.listen(3000);
```

## Проверка, что всё работает

```bash
node demo/attack-simulation.js
```

Должно вывести 6 успешных проверок (SQLi, XSS, поддельная почта, попытка регистрации админом — всё отклонено).

## Роли в системе

- `student` — обычный студент
- `mentor` — ментор
- `admin` — администратор (единственная роль, которую нельзя выбрать при регистрации — назначается только вручную в БД)
