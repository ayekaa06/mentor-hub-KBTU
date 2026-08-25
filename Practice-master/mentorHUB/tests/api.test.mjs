import assert from 'node:assert/strict';
import { existsSync, rmSync } from 'node:fs';
import { join } from 'node:path';
import test from 'node:test';

const databasePath = join(process.cwd(), 'data', 'mentorhub-api-test.db');
if (existsSync(databasePath)) rmSync(databasePath, { force: true });
process.env.MENTORHUB_DB = databasePath;
process.env.JWT_SECRET = 'test-only-secret';
const { app, closeDatabaseForTests } = await import('../dist/mentorHUB/server/server.mjs');
const server = app.listen(0);
await new Promise((resolve) => server.once('listening', resolve));
const base = `http://127.0.0.1:${server.address().port}`;

async function request(path, options = {}) {
  const response = await fetch(`${base}${path}`, options);
  return { response, body: await response.json() };
}

test('health endpoint and invalid credentials are handled safely', async () => {
  const health = await request('/api/health');
  assert.equal(health.response.status, 200);
  assert.equal(health.body.status, 'ok');
  const badLogin = await request('/api/auth/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'student@kbtu.kz', password: 'wrong' }) });
  assert.equal(badLogin.response.status, 401);
});

test('student can sign in, list mentors, schedule a meeting, and update a goal', async () => {
  const login = await request('/api/auth/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ email: 'student@kbtu.kz', password: 'Student2026!' }) });
  assert.equal(login.response.status, 200);
  assert.equal(login.body.user.role, 'student');
  const headers = { 'Content-Type': 'application/json', Authorization: `Bearer ${login.body.token}` };
  const mentors = await request('/api/mentors', { headers });
  assert.equal(mentors.response.status, 200);
  assert.ok(mentors.body.mentors.length >= 3);
  const meeting = await request('/api/meetings', { method: 'POST', headers, body: JSON.stringify({ mentorId: mentors.body.mentors[0].id, title: 'API integration check', startTime: new Date(Date.now() + 86400000).toISOString(), mode: 'online', durationMinutes: 30 }) });
  assert.equal(meeting.response.status, 201);
  const goals = await request('/api/goals', { headers });
  assert.equal(goals.response.status, 200);
  const update = await request(`/api/goals/${goals.body.goals[0].id}`, { method: 'PATCH', headers, body: JSON.stringify({ progress: 100 }) });
  assert.equal(update.response.status, 200);
  assert.equal(update.body.goal.completed, 1);
});

test.after(async () => {
  await new Promise((resolve) => server.close(resolve));
  closeDatabaseForTests();
  for (const suffix of ['', '-wal', '-shm']) { if (existsSync(`${databasePath}${suffix}`)) rmSync(`${databasePath}${suffix}`, { force: true }); }
});
