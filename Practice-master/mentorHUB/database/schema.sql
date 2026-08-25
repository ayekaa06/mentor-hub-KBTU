-- MentorHub KBTU relational schema (SQLite 3)
-- The application enables PRAGMA foreign_keys = ON for every connection.

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  full_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE COLLATE NOCASE,
  password_hash TEXT NOT NULL,
  role TEXT NOT NULL CHECK (role IN ('student', 'mentor', 'admin')),
  avatar_color TEXT NOT NULL DEFAULT '#2F6BFF',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS student_profiles (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  program TEXT NOT NULL,
  study_year INTEGER NOT NULL CHECK (study_year BETWEEN 1 AND 6),
  interests_json TEXT NOT NULL DEFAULT '[]',
  mentor_id INTEGER REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS mentor_profiles (
  user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
  program TEXT NOT NULL,
  graduation_year INTEGER,
  title TEXT NOT NULL,
  bio TEXT NOT NULL,
  skills_json TEXT NOT NULL DEFAULT '[]',
  languages_json TEXT NOT NULL DEFAULT '["English"]',
  rating REAL NOT NULL DEFAULT 5 CHECK (rating BETWEEN 0 AND 5),
  sessions_count INTEGER NOT NULL DEFAULT 0 CHECK (sessions_count >= 0),
  capacity INTEGER NOT NULL DEFAULT 5 CHECK (capacity > 0),
  availability TEXT NOT NULL DEFAULT 'Weekdays after 16:00'
);

CREATE TABLE IF NOT EXISTS mentoring_applications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  mentor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message TEXT NOT NULL,
  status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'declined')),
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE(student_id, mentor_id)
);

CREATE TABLE IF NOT EXISTS meetings (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  mentor_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  start_time TEXT NOT NULL,
  duration_minutes INTEGER NOT NULL DEFAULT 30 CHECK (duration_minutes BETWEEN 15 AND 180),
  mode TEXT NOT NULL CHECK (mode IN ('online', 'campus')),
  status TEXT NOT NULL DEFAULT 'scheduled' CHECK (status IN ('scheduled', 'completed', 'cancelled')),
  notes TEXT NOT NULL DEFAULT '',
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS goals (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  student_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  category TEXT NOT NULL,
  due_date TEXT NOT NULL,
  progress INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  completed INTEGER NOT NULL DEFAULT 0 CHECK (completed IN (0, 1)),
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS reviews (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  meeting_id INTEGER NOT NULL UNIQUE REFERENCES meetings(id) ON DELETE CASCADE,
  author_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
  comment TEXT NOT NULL,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS notifications (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  read_at TEXT,
  created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);
CREATE INDEX IF NOT EXISTS idx_meetings_student_time ON meetings(student_id, start_time);
CREATE INDEX IF NOT EXISTS idx_meetings_mentor_time ON meetings(mentor_id, start_time);
CREATE INDEX IF NOT EXISTS idx_goals_student ON goals(student_id, completed);
CREATE INDEX IF NOT EXISTS idx_applications_status ON mentoring_applications(status);
