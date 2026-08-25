# MentorHub KBTU

MentorHub is a full-stack mentorship platform for the KBTU community. It pairs students with experienced students and alumni, supports meeting planning and goals, and gives program staff a clear matching queue.

## Run the complete application

```bash
npm install
npm run start:full
```

Open [http://localhost:4000](http://localhost:4000). This command builds the Angular SSR application and starts the Express API and website together. The SQLite database is created automatically at `data/mentorhub.db` using `database/schema.sql`.

For frontend-only development with hot reload, use `npm start`. Use `npm run start:full` when testing login and API-backed functionality.

## Demo accounts

| Role | Email | Password |
| --- | --- | --- |
| Student | `student@kbtu.kz` | `Student2026!` |
| Mentor | `mentor@kbtu.kz` | `Mentor2026!` |
| Administrator | `admin@kbtu.kz` | `Admin2026!` |

## API and security

The API lives in `src/server.ts`. It uses parameterized SQL queries, scrypt password hashes, JWT-style HMAC signed access tokens, role checks for every private mutation, foreign keys and request validation. Set a strong `JWT_SECRET` environment variable before deploying.

Main endpoints include `/api/auth/login`, `/api/mentors`, `/api/meetings`, `/api/goals`, `/api/applications`, and `/api/admin/metrics`. `GET /api/health` is available for deployment checks.

## Tests and build

```bash
npm run build
npm test -- --watch=false
npm run test:api
```

The database model is documented directly in `database/schema.sql`: users, student and mentor profiles, applications, meetings, goals, reviews and notifications, with integrity constraints and indexes.
