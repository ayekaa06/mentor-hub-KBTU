import {
  AngularNodeAppEngine,
  createNodeRequestHandler,
  isMainModule,
  writeResponseToNodeResponse,
} from '@angular/ssr/node';
import { randomBytes, scryptSync, createHmac, timingSafeEqual } from 'node:crypto';
import { existsSync, mkdirSync, readFileSync } from 'node:fs';
import { join, dirname } from 'node:path';
import express, { type NextFunction, type Request, type Response } from 'express';

type Role = 'student' | 'mentor' | 'admin';
type User = { id: number; fullName: string; email: string; role: Role; avatarColor: string };
type AuthRequest = Request & { user?: User };

type Statement = {
  get(...params: unknown[]): unknown;
  all(...params: unknown[]): unknown[];
  run(...params: unknown[]): { lastInsertRowid: number | bigint; changes: number };
};
type Database = { exec(sql: string): void; prepare(sql: string): Statement; close(): void };

const getBuiltinModule = (process as any).getBuiltinModule;
const sqlite = getBuiltinModule ? getBuiltinModule('node:sqlite') as {
  DatabaseSync: new (path: string) => Database;
} : null;

const dbPath = process.env['MENTORHUB_DB'] || join(process.cwd(), 'data', 'mentorhub.db');
mkdirSync(dirname(dbPath), { recursive: true });

const db = sqlite ? new sqlite.DatabaseSync(dbPath) : {
  exec: (sql: string) => {},
  prepare: (sql: string) => ({
    get: (...params: unknown[]) => undefined as any,
    all: (...params: unknown[]) => [] as any[],
    run: (...params: unknown[]) => ({ lastInsertRowid: 0, changes: 0 })
  }),
  close: () => {}
} as any as Database;

if (sqlite) {
  db.exec('PRAGMA foreign_keys = ON; PRAGMA journal_mode = WAL;');
}

function readSchema(): string {
  const candidates = [
    join(process.cwd(), 'database', 'schema.sql'),
    join(import.meta.dirname, '..', '..', '..', 'database', 'schema.sql'),
  ];
  const schemaPath = candidates.find(existsSync);
  if (!schemaPath) throw new Error('database/schema.sql could not be found. Start the server from the project root.');
  return readFileSync(schemaPath, 'utf8');
}

function hashPassword(password: string, salt = randomBytes(16).toString('hex')): string {
  const hash = scryptSync(password, salt, 64).toString('hex');
  return `scrypt$${salt}$${hash}`;
}

function verifyPassword(password: string, stored: string): boolean {
  const [, salt, storedHash] = stored.split('$');
  if (!salt || !storedHash) return false;
  const calculated = hashPassword(password, salt).split('$')[2];
  const left = Buffer.from(calculated, 'hex');
  const right = Buffer.from(storedHash, 'hex');
  return left.length === right.length && timingSafeEqual(left, right);
}

function row<T>(sql: string, ...params: unknown[]): T | undefined {
  return db.prepare(sql).get(...params) as T | undefined;
}

function rows<T>(sql: string, ...params: unknown[]): T[] {
  return db.prepare(sql).all(...params) as T[];
}

function parseList(value: string | null | undefined): string[] {
  try {
    return JSON.parse(value || '[]') as string[];
  } catch {
    return [];
  }
}

function seedDatabase(): void {
  db.exec(readSchema());
  const seedSalt = 'mentorhub-kbtu-demo-salt';
  const demoUsers: Array<[string, string, Role, string, string]> = [
    ['Amina Nurgaliyeva', 'student@kbtu.kz', 'student', '#3578F6', 'Student2026!'],
    ['Elena Kadyrova', 'mentor@kbtu.kz', 'mentor', '#8B5CF6', 'Mentor2026!'],
    ['Daniyar Serik', 'daniyar@kbtu.kz', 'mentor', '#EC4899', 'Mentor2026!'],
    ['Sofia Karimova', 'sofia@kbtu.kz', 'mentor', '#00AFA0', 'Mentor2026!'],
    ['Nurlan Bekov', 'admin@kbtu.kz', 'admin', '#F59E0B', 'Admin2026!'],
  ];
  const insertUser = db.prepare(
    'INSERT OR IGNORE INTO users (full_name, email, role, avatar_color, password_hash) VALUES (?, ?, ?, ?, ?)',
  );
  for (const [name, email, role, color, password] of demoUsers) {
    insertUser.run(name, email, role, color, hashPassword(password, seedSalt));
  }

  const idFor = (email: string) => row<{ id: number }>('SELECT id FROM users WHERE email = ?', email)!.id;
  const amina = idFor('student@kbtu.kz');
  const elena = idFor('mentor@kbtu.kz');
  const daniyar = idFor('daniyar@kbtu.kz');
  const sofia = idFor('sofia@kbtu.kz');

  db.prepare(
    'INSERT OR IGNORE INTO student_profiles (user_id, program, study_year, interests_json, mentor_id) VALUES (?, ?, ?, ?, ?)',
  ).run(amina, 'Information Systems', 1, JSON.stringify(['Product design', 'Web development', 'Student life']), elena);

  const addMentor = db.prepare(
    `INSERT OR IGNORE INTO mentor_profiles
      (user_id, program, graduation_year, title, bio, skills_json, languages_json, rating, sessions_count, capacity, availability)
      VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
  );
  addMentor.run(elena, 'Computer Science', 2024, 'Software Engineer · Kaspi.kz', 'Helps new KBTU students turn uncertainty into a practical plan for their first semester.', JSON.stringify(['Web development', 'Career planning', 'Algorithms']), JSON.stringify(['English', 'Russian', 'Kazakh']), 4.9, 28, 6, 'Tue–Thu after 16:00');
  addMentor.run(daniyar, 'Information Systems', 2023, 'Product Analyst · Freedom Bank', 'Enjoys helping students discover product thinking, data analysis and meaningful university projects.', JSON.stringify(['Product management', 'Data analytics', 'CV review']), JSON.stringify(['English', 'Russian']), 4.8, 19, 5, 'Mon, Wed, Fri after 17:00');
  addMentor.run(sofia, 'Management', 2025, 'KBTU Student Ambassador', 'Your guide to communities, academic routines and finding your place at KBTU.', JSON.stringify(['Campus life', 'Public speaking', 'Time management']), JSON.stringify(['English', 'Kazakh']), 5, 32, 8, 'Weekdays 14:00–18:00');

  const meetingCount = row<{ count: number }>('SELECT COUNT(*) AS count FROM meetings')!.count;
  if (meetingCount === 0) {
    const addMeeting = db.prepare(
      `INSERT INTO meetings (student_id, mentor_id, title, start_time, duration_minutes, mode, status, notes)
       VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
    );
    addMeeting.run(amina, elena, 'First-semester game plan', new Date(Date.now() + 2 * 86400000).toISOString(), 45, 'online', 'scheduled', 'Bring your current course list and two questions.');
    addMeeting.run(amina, elena, 'Getting to know KBTU', new Date(Date.now() - 4 * 86400000).toISOString(), 30, 'campus', 'completed', 'Discussed clubs, LMS and building navigation.');
  }
  const goalCount = row<{ count: number }>('SELECT COUNT(*) AS count FROM goals')!.count;
  if (goalCount === 0) {
    const addGoal = db.prepare('INSERT INTO goals (student_id, title, category, due_date, progress, completed) VALUES (?, ?, ?, ?, ?, ?)');
    addGoal.run(amina, 'Build a reliable study routine', 'Academic', new Date(Date.now() + 14 * 86400000).toISOString(), 70, 0);
    addGoal.run(amina, 'Join one student community', 'Campus life', new Date(Date.now() + 21 * 86400000).toISOString(), 40, 0);
    addGoal.run(amina, 'Prepare a first CV draft', 'Career', new Date(Date.now() + 30 * 86400000).toISOString(), 15, 0);
  }
  const appCount = row<{ count: number }>('SELECT COUNT(*) AS count FROM mentoring_applications')!.count;
  if (appCount === 0) {
    db.prepare('INSERT INTO mentoring_applications (student_id, mentor_id, message, status) VALUES (?, ?, ?, ?)').run(amina, daniyar, 'I would love guidance on product analytics and choosing electives.', 'pending');
  }
}

if (sqlite) {
  seedDatabase();
}

const JWT_SECRET = process.env['JWT_SECRET'] || 'mentorhub-development-secret-change-before-deployment';
function signToken(user: User): string {
  const payload = Buffer.from(JSON.stringify({ sub: user.id, role: user.role, exp: Math.floor(Date.now() / 1000) + 8 * 60 * 60 })).toString('base64url');
  const signature = createHmac('sha256', JWT_SECRET).update(payload).digest('base64url');
  return `${payload}.${signature}`;
}

function verifyToken(token: string): { sub: number; role: Role } | null {
  const [payload, signature] = token.split('.');
  if (!payload || !signature) return null;
  const expected = createHmac('sha256', JWT_SECRET).update(payload).digest('base64url');
  const left = Buffer.from(signature);
  const right = Buffer.from(expected);
  if (left.length !== right.length || !timingSafeEqual(left, right)) return null;
  try {
    const decoded = JSON.parse(Buffer.from(payload, 'base64url').toString('utf8')) as { sub: number; role: Role; exp: number };
    return decoded.exp > Math.floor(Date.now() / 1000) ? { sub: decoded.sub, role: decoded.role } : null;
  } catch {
    return null;
  }
}

function toUser(record: any): User {
  return { id: Number(record.id), fullName: record.full_name, email: record.email, role: record.role as Role, avatarColor: record.avatar_color };
}

function authenticate(req: AuthRequest, res: Response, next: NextFunction): void {
  const token = req.header('Authorization')?.replace(/^Bearer\s+/i, '');
  const tokenData = token ? verifyToken(token) : null;
  const record = tokenData ? row<any>('SELECT id, full_name, email, role, avatar_color FROM users WHERE id = ?', tokenData.sub) : undefined;
  if (!record) {
    res.status(401).json({ message: 'Authentication is required.' });
    return;
  }
  req.user = toUser(record);
  next();
}

function only(...allowedRoles: Role[]) {
  return (req: AuthRequest, res: Response, next: NextFunction): void => {
    if (!req.user || !allowedRoles.includes(req.user.role)) {
      res.status(403).json({ message: 'You do not have permission to perform this action.' });
      return;
    }
    next();
  };
}

function cleanText(value: unknown, label: string, maxLength = 300): string | null {
  if (typeof value !== 'string') return null;
  const cleaned = value.trim().replace(/\s+/g, ' ');
  return cleaned.length > 0 && cleaned.length <= maxLength ? cleaned : null;
}

function publicMentor(record: any) {
  return {
    id: Number(record.id),
    name: record.full_name,
    color: record.avatar_color,
    program: record.program,
    graduationYear: record.graduation_year,
    title: record.title,
    bio: record.bio,
    skills: parseList(record.skills_json),
    languages: parseList(record.languages_json),
    rating: Number(record.rating),
    sessionsCount: Number(record.sessions_count),
    capacity: Number(record.capacity),
    availability: record.availability,
  };
}

const browserDistFolder = join(import.meta.dirname, '../browser');
export const app = express();
app.disable('x-powered-by');
app.use((_, res, next) => {
  res.setHeader('X-Content-Type-Options', 'nosniff');
  res.setHeader('X-Frame-Options', 'DENY');
  res.setHeader('Referrer-Policy', 'strict-origin-when-cross-origin');
  next();
});
app.use(express.json({ limit: '50kb' }));

app.get('/api/health', (_, res) => res.json({ status: 'ok', service: 'mentorhub-api' }));

app.post('/api/auth/login', (req, res) => {
  const email = cleanText(req.body?.email, 'email', 120)?.toLowerCase();
  const password = typeof req.body?.password === 'string' ? req.body.password : '';
  const record = email ? row<any>('SELECT * FROM users WHERE email = ?', email) : undefined;
  if (!record || !verifyPassword(password, record.password_hash)) {
    res.status(401).json({ message: 'Incorrect email or password.' });
    return;
  }
  const user = toUser(record);
  res.json({ token: signToken(user), user });
});

app.get('/api/auth/me', authenticate, (req: AuthRequest, res) => res.json({ user: req.user }));

app.get('/api/mentors', (_req, res) => {
  const mentors = rows<any>(
    `SELECT u.id, u.full_name, u.avatar_color, mp.*
       FROM users u JOIN mentor_profiles mp ON mp.user_id = u.id
       ORDER BY mp.rating DESC, mp.sessions_count DESC`,
  ).map(publicMentor);
  res.json({ mentors });
});

app.get('/api/profile', authenticate, (req: AuthRequest, res) => {
  const user = req.user!;
  const profile = user.role === 'student'
    ? row<any>(`SELECT sp.*, mentor.full_name AS mentor_name FROM student_profiles sp
                 LEFT JOIN users mentor ON mentor.id = sp.mentor_id WHERE sp.user_id = ?`, user.id)
    : user.role === 'mentor'
      ? row<any>('SELECT * FROM mentor_profiles WHERE user_id = ?', user.id)
      : null;
  res.json({ user, profile: profile ? { ...profile, interests: parseList(profile.interests_json), skills: parseList(profile.skills_json), languages: parseList(profile.languages_json) } : null });
});

app.patch('/api/profile', authenticate, (req: AuthRequest, res) => {
  const name = cleanText(req.body?.fullName, 'full name', 80);
  if (!name) {
    res.status(400).json({ message: 'Please provide a name between 1 and 80 characters.' });
    return;
  }
  db.prepare("UPDATE users SET full_name = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?").run(name, req.user!.id);
  const record = row<any>('SELECT id, full_name, email, role, avatar_color FROM users WHERE id = ?', req.user!.id)!;
  res.json({ user: toUser(record) });
});

app.get('/api/dashboard', authenticate, (req: AuthRequest, res) => {
  const user = req.user!;
  const meetings = user.role === 'mentor'
    ? rows<any>(`SELECT m.*, s.full_name AS student_name, s.avatar_color AS student_color FROM meetings m JOIN users s ON s.id=m.student_id WHERE m.mentor_id=? ORDER BY m.start_time ASC LIMIT 4`, user.id)
    : user.role === 'student'
      ? rows<any>(`SELECT m.*, mentor.full_name AS mentor_name, mentor.avatar_color AS mentor_color FROM meetings m JOIN users mentor ON mentor.id=m.mentor_id WHERE m.student_id=? ORDER BY m.start_time ASC LIMIT 4`, user.id)
      : rows<any>(`SELECT m.*, s.full_name AS student_name, mentor.full_name AS mentor_name FROM meetings m JOIN users s ON s.id=m.student_id JOIN users mentor ON mentor.id=m.mentor_id ORDER BY m.start_time ASC LIMIT 4`);
  const metrics = user.role === 'student'
    ? { completedGoals: row<{ count: number }>('SELECT COUNT(*) AS count FROM goals WHERE student_id=? AND completed=1', user.id)!.count, totalGoals: row<{ count: number }>('SELECT COUNT(*) AS count FROM goals WHERE student_id=?', user.id)!.count, mentorName: row<{ name: string }>('SELECT u.full_name AS name FROM student_profiles sp JOIN users u ON u.id=sp.mentor_id WHERE sp.user_id=?', user.id)?.name || 'Choose a mentor' }
    : user.role === 'mentor'
      ? { activeStudents: row<{ count: number }>('SELECT COUNT(DISTINCT student_id) AS count FROM meetings WHERE mentor_id=?', user.id)!.count, upcoming: row<{ count: number }>("SELECT COUNT(*) AS count FROM meetings WHERE mentor_id=? AND status='scheduled'", user.id)!.count }
      : { students: row<{ count: number }>("SELECT COUNT(*) AS count FROM users WHERE role='student'")!.count, mentors: row<{ count: number }>("SELECT COUNT(*) AS count FROM users WHERE role='mentor'")!.count, meetings: row<{ count: number }>('SELECT COUNT(*) AS count FROM meetings')!.count };
  res.json({ meetings, metrics });
});

app.get('/api/meetings', authenticate, (req: AuthRequest, res) => {
  const user = req.user!;
  const meetingSql = user.role === 'mentor'
    ? `SELECT m.*, s.full_name AS student_name, s.avatar_color AS student_color, mentor.full_name AS mentor_name FROM meetings m JOIN users s ON s.id=m.student_id JOIN users mentor ON mentor.id=m.mentor_id WHERE m.mentor_id=? ORDER BY m.start_time ASC`
    : user.role === 'student'
      ? `SELECT m.*, mentor.full_name AS mentor_name, mentor.avatar_color AS mentor_color, s.full_name AS student_name FROM meetings m JOIN users mentor ON mentor.id=m.mentor_id JOIN users s ON s.id=m.student_id WHERE m.student_id=? ORDER BY m.start_time ASC`
      : `SELECT m.*, mentor.full_name AS mentor_name, s.full_name AS student_name FROM meetings m JOIN users mentor ON mentor.id=m.mentor_id JOIN users s ON s.id=m.student_id ORDER BY m.start_time ASC`;
  res.json({ meetings: rows<any>(meetingSql, ...(user.role === 'admin' ? [] : [user.id])) });
});

app.post('/api/meetings', authenticate, only('student'), (req: AuthRequest, res) => {
  const title = cleanText(req.body?.title, 'title', 100);
  const mentorId = Number(req.body?.mentorId);
  const start = new Date(req.body?.startTime);
  const duration = Number(req.body?.durationMinutes || 30);
  const mode = req.body?.mode === 'campus' ? 'campus' : req.body?.mode === 'online' ? 'online' : null;
  if (!title || !Number.isInteger(mentorId) || Number.isNaN(start.valueOf()) || start <= new Date() || !Number.isInteger(duration) || duration < 15 || duration > 180 || !mode) {
    res.status(400).json({ message: 'Please provide a valid future date, mentor, meeting title and duration.' });
    return;
  }
  if (!row('SELECT user_id FROM mentor_profiles WHERE user_id=?', mentorId)) {
    res.status(404).json({ message: 'Mentor not found.' });
    return;
  }
  const result = db.prepare('INSERT INTO meetings (student_id, mentor_id, title, start_time, duration_minutes, mode, status, notes) VALUES (?, ?, ?, ?, ?, ?, ?, ?)').run(req.user!.id, mentorId, title, start.toISOString(), duration, mode, 'scheduled', '');
  const meeting = row<any>('SELECT m.*, mentor.full_name AS mentor_name FROM meetings m JOIN users mentor ON mentor.id=m.mentor_id WHERE m.id=?', Number(result.lastInsertRowid));
  res.status(201).json({ meeting });
});

app.patch('/api/meetings/:id', authenticate, (req: AuthRequest, res) => {
  const id = Number(req.params['id']);
  const meeting = row<any>('SELECT * FROM meetings WHERE id=?', id);
  if (!meeting) {
    res.status(404).json({ message: 'Meeting not found.' });
    return;
  }
  if (req.user!.role !== 'admin' && meeting.student_id !== req.user!.id && meeting.mentor_id !== req.user!.id) {
    res.status(403).json({ message: 'You cannot update this meeting.' });
    return;
  }
  const status = req.body?.status;
  if (!['scheduled', 'completed', 'cancelled'].includes(status)) {
    res.status(400).json({ message: 'A valid meeting status is required.' });
    return;
  }
  db.prepare('UPDATE meetings SET status=? WHERE id=?').run(status, id);
  res.json({ meeting: row<any>('SELECT * FROM meetings WHERE id=?', id) });
});

app.get('/api/goals', authenticate, only('student'), (req: AuthRequest, res) => {
  res.json({ goals: rows<any>('SELECT * FROM goals WHERE student_id=? ORDER BY completed ASC, due_date ASC', req.user!.id) });
});

app.post('/api/goals', authenticate, only('student'), (req: AuthRequest, res) => {
  const title = cleanText(req.body?.title, 'title', 140);
  const category = cleanText(req.body?.category, 'category', 50);
  const dueDate = new Date(req.body?.dueDate);
  if (!title || !category || Number.isNaN(dueDate.valueOf())) {
    res.status(400).json({ message: 'A goal needs a title, category and due date.' });
    return;
  }
  const result = db.prepare('INSERT INTO goals (student_id, title, category, due_date) VALUES (?, ?, ?, ?)').run(req.user!.id, title, category, dueDate.toISOString());
  res.status(201).json({ goal: row<any>('SELECT * FROM goals WHERE id=?', Number(result.lastInsertRowid)) });
});

app.patch('/api/goals/:id', authenticate, only('student'), (req: AuthRequest, res) => {
  const id = Number(req.params['id']);
  const progress = Number(req.body?.progress);
  if (!Number.isInteger(progress) || progress < 0 || progress > 100) {
    res.status(400).json({ message: 'Progress must be a whole number from 0 to 100.' });
    return;
  }
  const result = db.prepare('UPDATE goals SET progress=?, completed=? WHERE id=? AND student_id=?').run(progress, progress === 100 ? 1 : 0, id, req.user!.id);
  if (!result.changes) {
    res.status(404).json({ message: 'Goal not found.' });
    return;
  }
  res.json({ goal: row<any>('SELECT * FROM goals WHERE id=?', id) });
});

app.get('/api/applications', authenticate, only('admin'), (_req, res) => {
  const applications = rows<any>(`SELECT a.*, s.full_name AS student_name, mentor.full_name AS mentor_name
    FROM mentoring_applications a JOIN users s ON s.id=a.student_id JOIN users mentor ON mentor.id=a.mentor_id
    ORDER BY CASE a.status WHEN 'pending' THEN 0 ELSE 1 END, a.created_at DESC`);
  res.json({ applications });
});

app.post('/api/applications', authenticate, only('student'), (req: AuthRequest, res) => {
  const mentorId = Number(req.body?.mentorId);
  const message = cleanText(req.body?.message, 'message', 400) || 'I would like to request mentorship.';
  if (!Number.isInteger(mentorId) || !row('SELECT user_id FROM mentor_profiles WHERE user_id=?', mentorId)) {
    res.status(400).json({ message: 'Please select a valid mentor.' });
    return;
  }
  try {
    const result = db.prepare('INSERT INTO mentoring_applications (student_id, mentor_id, message) VALUES (?, ?, ?)').run(req.user!.id, mentorId, message);
    res.status(201).json({ application: row<any>('SELECT * FROM mentoring_applications WHERE id=?', Number(result.lastInsertRowid)) });
  } catch {
    res.status(409).json({ message: 'You have already sent a request to this mentor.' });
  }
});

app.patch('/api/applications/:id', authenticate, only('admin'), (req: AuthRequest, res) => {
  const id = Number(req.params['id']);
  const status = req.body?.status;
  if (!['approved', 'declined'].includes(status)) {
    res.status(400).json({ message: 'Choose approved or declined.' });
    return;
  }
  const application = row<any>('SELECT * FROM mentoring_applications WHERE id=?', id);
  if (!application) {
    res.status(404).json({ message: 'Application not found.' });
    return;
  }
  db.prepare('UPDATE mentoring_applications SET status=? WHERE id=?').run(status, id);
  if (status === 'approved') {
    db.prepare('UPDATE student_profiles SET mentor_id=? WHERE user_id=?').run(application.mentor_id, application.student_id);
  }
  res.json({ application: row<any>('SELECT * FROM mentoring_applications WHERE id=?', id) });
});

app.get('/api/admin/metrics', authenticate, only('admin'), (_req, res) => {
  res.json({
    metrics: {
      students: row<{ count: number }>("SELECT COUNT(*) AS count FROM users WHERE role='student'")!.count,
      mentors: row<{ count: number }>("SELECT COUNT(*) AS count FROM users WHERE role='mentor'")!.count,
      meetings: row<{ count: number }>('SELECT COUNT(*) AS count FROM meetings')!.count,
      pendingApplications: row<{ count: number }>("SELECT COUNT(*) AS count FROM mentoring_applications WHERE status='pending'")!.count,
      completionRate: row<{ rate: number }>("SELECT COALESCE(ROUND(100.0 * SUM(CASE WHEN status='completed' THEN 1 ELSE 0 END) / NULLIF(COUNT(*), 0)), 0) AS rate FROM meetings")!.rate,
    },
  });
});

app.use(express.static(browserDistFolder, { maxAge: '1y', index: false, redirect: false }));
const angularApp = new AngularNodeAppEngine();
app.use((req, res, next) => {
  angularApp.handle(req).then((response) => response ? writeResponseToNodeResponse(response, res) : next()).catch(next);
});
app.use((error: unknown, _req: Request, res: Response, _next: NextFunction) => {
  console.error(error);
  res.status(500).json({ message: 'Something went wrong on the server.' });
});

if (isMainModule(import.meta.url) || process.env['pm_id']) {
  const port = Number(process.env['PORT'] || 4000);
  app.listen(port, () => console.log(`MentorHub is running at http://localhost:${port}`));
}

export const reqHandler = createNodeRequestHandler(app);

/** Allows the integration test process to release the SQLite file before cleanup. */
export function closeDatabaseForTests(): void { db.close(); }
