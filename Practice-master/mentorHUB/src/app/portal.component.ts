import { CommonModule, isPlatformBrowser } from '@angular/common';
import { Component, PLATFORM_ID, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterLink, RouterLinkActive } from '@angular/router';
import { ApiService, type Role, type Session, type SessionUser, mapGoRole, avatarColor } from './services/api.service';

// ── Локальные типы ────────────────────────────────────────────────────────────

/** Пользователь с ролью mentor, отображается на странице /mentors */
type Mentor = {
  id: string;
  name: string;
  color: string;
  program?: string;
  graduationYear?: number;
  title?: string;
  bio?: string;
  skills: string[];
  languages?: string[];
  rating?: number;
  sessionsCount?: number;
  capacity?: number;
  availability?: string;
};

/** Задача freshman-а (GET /api/v1/freshman/tasks) */
type Task = {
  id: string;
  title: string;
  description?: string;
  status: 'pending' | 'submitted' | 'approved' | 'rejected';
  due_date?: string;
  proof_url?: string;
  comment?: string;
  progress?: number; // вычисляемое: 0 | 50 | 100 (pending | submitted | approved)
  completed?: boolean;
  category?: string;
};

/** Встреча (GET /api/v1/freshman/meetings | /api/v1/mentor/meetings) */
type Meeting = {
  id: string;
  title: string;
  scheduled_at?: string;
  start_time?: string;
  duration_minutes?: number;
  mode?: string;
  status: string;
  mentor_name?: string;
  student_name?: string;
  mentor_color?: string;
  student_color?: string;
  notes?: string;
};

@Component({
  selector: 'app-portal',
  imports: [CommonModule, FormsModule, RouterLink, RouterLinkActive],
  template: `
    <div class="portal-shell">
      <header class="topbar">
        <a routerLink="/" class="brand" aria-label="MentorHub KBTU home">
          <span class="brand-mark"><i></i><i></i><i></i></span>
          <span>MentorHub <em>KBTU</em></span>
        </a>
        <nav class="main-nav" aria-label="Main navigation">
          <a routerLink="/" routerLinkActive="active" [routerLinkActiveOptions]="{ exact: true }"
            >Discover</a
          >
          <a routerLink="/mentors" routerLinkActive="active">Mentors</a>
          @if (session()) {
            <a routerLink="/meetings" routerLinkActive="active">Meetings</a>
            <a routerLink="/roadmap" routerLinkActive="active">{{ session()!.user.role === 'student' ? 'My Tasks' : 'Roadmap' }}</a>
            @if (session()!.user.role === 'admin') {
              <a routerLink="/admin" routerLinkActive="active">Insights</a>
            }
          }
        </nav>
        <div class="header-actions">
          @if (session(); as activeSession) {
            <a
              routerLink="/profile"
              class="avatar-button"
              [style.background]="activeSession.user.avatarColor"
              [attr.aria-label]="'Open ' + activeSession.user.fullName + ' profile'"
              >{{ initials(activeSession.user.fullName) }}</a
            >
            <button class="text-button logout" type="button" (click)="logout()">Sign out</button>
          } @else {
            <a routerLink="/login" class="text-button">Sign in</a>
            <a routerLink="/login" class="button button-dark">Join MentorHub <span>↗</span></a>
          }
        </div>
      </header>

      @if (notice(); as message) {
        <div class="toast" role="status">
          <span>✦</span>{{ message
          }}<button type="button" (click)="notice.set('')" aria-label="Dismiss">×</button>
        </div>
      }

      @if (view === 'home') {
        <main>
          <section class="hero-section page-width">
            <div class="hero-copy">
              <p class="eyebrow"><span></span> KBTU peer mentorship program</p>
              <h1>Start boldly.<br /><span>Grow together.</span></h1>
              <p class="hero-description">
                A more human start to university. MentorHub connects KBTU students with people who
                have already found their way.
              </p>
              <div class="hero-actions">
                <a routerLink="/mentors" class="button button-blue"
                  >Find your mentor <span>→</span></a
                >
                <a href="#how-it-works" class="button button-quiet"
                  ><span class="play-icon">▶</span> See how it works</a
                >
              </div>
              <div class="community-proof">
                <div class="avatar-stack">
                  <b style="background:#d8665e">AK</b><b style="background:#644ec2">DS</b
                  ><b style="background:#2e9c86">SK</b><b style="background:#e7a12a">NB</b>
                </div>
                <p><strong>120+ students</strong> are already building their path with a mentor</p>
              </div>
            </div>
            <div class="hero-art" aria-label="Mentor and student working together">
              <div class="orb orb-one"></div>
              <div class="orb orb-two"></div>
              <div class="grid-dots"></div>
              <div class="floating-card first">
                <span class="mini-avatar violet">E</span>
                <p><b>Elena K.</b><small>Software Engineer</small></p>
                <span class="verified">✓</span>
              </div>
              <div class="collaboration-card">
                <div class="people-illustration">
                  <span class="person p1"></span><span class="person p2"></span
                  ><span class="laptop"></span>
                </div>
                <div class="card-caption">
                  <span class="online-dot"></span> Meaningful guidance, one conversation at a time
                </div>
              </div>
              <div class="floating-card second">
                <span class="mini-avatar coral">A</span>
                <p><b>Amina N.</b><small>First-year student</small></p>
                <span class="sparkle">✦</span>
              </div>
              <div class="meeting-pill">
                <span>◷</span>
                <div><b>Next session</b><small>Tuesday · 16:30</small></div>
              </div>
            </div>
          </section>

          <section class="trust-strip">
            <div class="page-width trust-inner">
              <p>Built for the people who make KBTU feel like home</p>
              <div>
                <b>STUDENTS</b><i></i><b>MENTORS</b><i></i><b>ALUMNI</b><i></i><b>FACULTY</b>
              </div>
            </div>
          </section>

          <section id="how-it-works" class="page-width process-section">
            <div class="section-heading">
              <p class="eyebrow"><span></span> A thoughtful beginning</p>
              <h2>Big ambitions deserve<br />a steady <span>first step.</span></h2>
              <p>
                Whether you need a friendly face, honest advice, or a little direction — your mentor
                is here to help you move forward.
              </p>
            </div>
            <div class="process-grid">
              <article>
                <span class="step-number">01</span>
                <div class="icon-circle">⌕</div>
                <h3>Find your fit</h3>
                <p>Explore mentors by faculty, experience and the topics that matter to you.</p>
              </article>
              <article>
                <span class="step-number">02</span>
                <div class="icon-circle">↗</div>
                <h3>Make the first move</h3>
                <p>Send a quick request. Your mentor will get back to you personally.</p>
              </article>
              <article>
                <span class="step-number">03</span>
                <div class="icon-circle">✦</div>
                <h3>Keep growing</h3>
                <p>Meet, set goals and celebrate progress — at your own pace.</p>
              </article>
            </div>
          </section>

          <section class="mentor-callout page-width">
            <div>
              <p class="eyebrow light"><span></span> Share what you know</p>
              <h2>Someone needs the<br />perspective <em>you</em> have.</h2>
              <p>
                Mentoring is a small commitment with a big ripple effect. Help a new student feel
                seen, supported and capable.
              </p>
              <a routerLink="/login" class="button button-white">Become a mentor <span>→</span></a>
            </div>
            <div class="callout-shape">
              <div class="ripple r1"></div>
              <div class="ripple r2"></div>
              <div class="ripple r3"></div>
              <span>✦</span>
            </div>
          </section>
          <footer class="site-footer page-width">
            <a routerLink="/" class="brand"
              ><span class="brand-mark"><i></i><i></i><i></i></span
              ><span>MentorHub <em>KBTU</em></span></a
            >
            <p>Made with care for the KBTU community.</p>
            <span>© 2026 MentorHub KBTU</span>
          </footer>
        </main>
      } @else if (view === 'login') {
        <main class="auth-layout page-width">
          <section class="auth-intro">
            <p class="eyebrow"><span></span> Welcome to the community</p>
            <h1>Your next chapter<br />starts <em>with people.</em></h1>
            <p>Sign in to find guidance, plan your semester and make every KBTU day count.</p>
            <div class="quote-card">
              <span>"</span>
              <p>
                My mentor made the first month feel possible. I finally knew who to ask and where I
                was headed.
              </p>
              <b>— Amina, Information Systems · Year 1</b>
            </div>
          </section>
          <section class="auth-card">
            <div>
              <span class="auth-icon">✦</span>
              <p class="eyebrow"><span></span> Secure access</p>
              <h2>Welcome back</h2>
              <p class="muted">Enter your KBTU credentials or use a demo account.</p>
            </div>
            <label
              >Email<input
                type="email"
                [(ngModel)]="loginForm.email"
                placeholder="you@kbtu.kz"
                autocomplete="email"
            /></label>
            <label
              >Password<input
                type="password"
                [(ngModel)]="loginForm.password"
                placeholder="Your password"
                autocomplete="current-password"
                (keyup.enter)="login()"
            /></label>
            <button
              type="button"
              class="button button-blue full"
              [disabled]="loading()"
              (click)="login()"
            >
              {{ loading() ? 'Signing you in…' : 'Sign in to MentorHub' }} <span>→</span>
            </button>
            <div class="demo-logins">
              <p>Demo access</p>
              <button type="button" (click)="useDemo('student')">
                <b>Freshman</b><small>freshman@kbtu.kz</small></button
              ><button type="button" (click)="useDemo('mentor')">
                <b>Mentor</b><small>mentor@kbtu.kz</small></button
              ><button type="button" (click)="useDemo('admin')">
                <b>Head</b><small>head@kbtu.kz</small>
              </button>
            </div>
          </section>
        </main>
      } @else if (!session()) {
        <main class="access-layout page-width">
          <div class="empty-state">
            <span class="large-icon">✦</span>
            <p class="eyebrow"><span></span> MentorHub is personal</p>
            <h1>Sign in to continue your journey.</h1>
            <p>Your meetings, goals and mentoring matches are kept in your private workspace.</p>
            <a routerLink="/login" class="button button-blue">Sign in to continue <span>→</span></a>
          </div>
        </main>
      } @else if (view === 'dashboard') {
        <main class="portal-main page-width">
          <section class="portal-hero">
            <div>
              <p class="eyebrow">
                <span></span> {{ greeting() }}, {{ firstName(session()!.user.fullName) }}
              </p>
              <h1>{{ dashboardHeadline() }}</h1>
              <p>{{ dashboardSubline() }}</p>
            </div>
            <div class="hero-date">
              <span>{{ today | date: 'EEE' }}</span
              ><b>{{ today | date: 'd' }}</b
              ><small>{{ today | date: 'MMMM' }}</small>
            </div>
          </section>

          @if (session()!.user.role === 'student') {
            <!-- Dashboard для freshman: задачи вместо "roadmap progress" -->
            <section class="metric-grid">
              <article class="metric-card blue">
                <p>Task progress</p>
                <b>{{ taskProgress() }}<small>%</small></b
                ><span>Keep showing up — you're building momentum.</span>
                <div class="progress"><i [style.width.%]="taskProgress()"></i></div>
              </article>
              <article class="metric-card">
                <p>Pending tasks</p>
                <b>{{ pendingTaskCount() }}</b>
                <span>Tasks waiting for your attention.</span>
                <a routerLink="/roadmap">View all tasks →</a>
              </article>
              <article class="metric-card">
                <p>Upcoming meetings</p>
                <b>{{ upcomingCount() }}</b>
                <span>Small conversations, real progress.</span>
                <a routerLink="/meetings">Open calendar →</a>
              </article>
            </section>
          } @else if (session()!.user.role === 'mentor') {
            <!-- Dashboard для mentor -->
            <section class="metric-grid">
              <article class="metric-card blue">
                <p>Students supported</p>
                <b>{{ dashboard()?.student_count || 0 }}</b>
                <span>Your guidance is already making an impact.</span>
              </article>
              <article class="metric-card">
                <p>Sessions coming up</p>
                <b>{{ upcomingCount() }}</b>
                <span>Make space for a thoughtful conversation.</span>
              </article>
              <article class="metric-card">
                <p>Tasks to review</p>
                <b>{{ dashboard()?.pending_reviews || 0 }}</b>
                <span>Submitted tasks waiting for your approval.</span>
                <a routerLink="/roadmap">Review tasks →</a>
              </article>
            </section>
          } @else {
            <!-- Dashboard для head/advisor -->
            <section class="metric-grid">
              <article class="metric-card blue">
                <p>Active freshmen</p>
                <b>{{ dashboard()?.total_freshmen || 0 }}</b>
                <span>Students currently in the program.</span>
              </article>
              <article class="metric-card">
                <p>Mentors</p>
                <b>{{ dashboard()?.total_mentors || 0 }}</b>
                <span>Experienced people ready to help.</span>
              </article>
              <article class="metric-card">
                <p>Overdue tasks</p>
                <b>{{ dashboard()?.overdue_tasks || 0 }}</b>
                <span>Tasks past their deadline.</span>
                <a routerLink="/admin">Open insights →</a>
              </article>
            </section>
          }

          <section class="content-grid">
            <div class="surface agenda">
              <div class="surface-heading">
                <div>
                  <p class="label">YOUR RHYTHM</p>
                  <h2>Coming up</h2>
                </div>
                <a routerLink="/meetings">Full calendar →</a>
              </div>
              @if (meetings().length) {
                @for (meeting of meetings().slice(0, 4); track meeting.id) {
                  <article class="agenda-row">
                    <div class="agenda-date">
                      <b>{{ (meeting.scheduled_at || meeting.start_time) | date: 'd' }}</b>
                      <span>{{ (meeting.scheduled_at || meeting.start_time) | date: 'MMM' }}</span>
                    </div>
                    <div>
                      <b>{{ meeting.title }}</b>
                      <p>
                        {{ meeting.mentor_name || meeting.student_name }} ·
                        {{ (meeting.scheduled_at || meeting.start_time) | date: 'HH:mm' }} · {{ meeting.mode || 'online' }}
                      </p>
                    </div>
                    <span class="status" [class.done]="meeting.status === 'completed'">{{
                      meeting.status
                    }}</span>
                  </article>
                }
              } @else {
                <div class="blank-row">
                  No sessions yet. The right conversation can start today.
                </div>
              }
            </div>
            <aside class="surface nudge">
              <p class="label">A GENTLE NUDGE</p>
              <span class="nudge-art">✦</span>
              <h2>Progress loves consistency.</h2>
              <p>Set one small, meaningful goal for this week. Your future self will notice.</p>
              <a routerLink="/roadmap" class="button button-dark">{{ session()!.user.role === 'student' ? 'Open my tasks' : 'Open roadmap' }} <span>→</span></a>
            </aside>
          </section>
        </main>
      } @else if (view === 'mentors') {
        <main class="portal-main page-width">
          <section class="directory-heading">
            <p class="eyebrow"><span></span> People who have been there</p>
            <h1>Find your person.</h1>
            <p>
              Browse experienced KBTU students and alumni who are ready to listen, share and help
              you take the next step.
            </p>
          </section>
          <section class="directory-tools surface">
            <label class="search-field"
              ><span>⌕</span
              ><input [(ngModel)]="mentorSearch" placeholder="Search by name, skill or topic"
            /></label>
            <div class="filter-chips">
              <button
                type="button"
                [class.selected]="mentorFocus === 'All'"
                (click)="mentorFocus = 'All'"
              >
                All topics</button
              ><button
                type="button"
                [class.selected]="mentorFocus === 'Career'"
                (click)="mentorFocus = 'Career'"
              >
                Career</button
              ><button
                type="button"
                [class.selected]="mentorFocus === 'Academic'"
                (click)="mentorFocus = 'Academic'"
              >
                Academic</button
              ><button
                type="button"
                [class.selected]="mentorFocus === 'Campus'"
                (click)="mentorFocus = 'Campus'"
              >
                Campus life
              </button>
            </div>
          </section>
          <section class="mentor-grid">
            @for (mentor of visibleMentors(); track mentor.id) {
              <article class="mentor-card">
                <div class="mentor-card-top">
                  <span class="profile-orb" [style.background]="mentor.color">{{
                    initials(mentor.name)
                  }}</span>
                  <div><span class="availability-dot"></span> {{ mentor.availability || 'Available' }}</div>
                  <button
                    type="button"
                    class="icon-button"
                    (click)="notice.set(mentor.name + ' saved to your favourites.')"
                    aria-label="Save mentor"
                  >
                    ♡
                  </button>
                </div>
                <p class="program">{{ mentor.program || 'KBTU' }}{{ mentor.graduationYear ? ' · Class of ' + mentor.graduationYear : '' }}</p>
                <h2>{{ mentor.name }}</h2>
                <p class="mentor-title">{{ mentor.title || 'Mentor' }}</p>
                <p class="mentor-bio">{{ mentor.bio || 'Ready to help you navigate your first year.' }}</p>
                <div class="tag-row">
                  @for (skill of (mentor.skills || []).slice(0, 4); track skill) {
                    <span>{{ skill }}</span>
                  }
                </div>
                <div class="mentor-card-footer">
                  <p>
                    <b>★ {{ mentor.rating || '5.0' }}</b> <span>· {{ mentor.sessionsCount || 0 }} sessions</span>
                  </p>
                  <button
                    type="button"
                    class="button button-dark small"
                    (click)="requestMentor(mentor)"
                  >
                    {{ session()?.user?.role === 'student' ? 'Connect' : 'View profile' }}
                    <span>→</span>
                  </button>
                </div>
              </article>
            } @empty {
              <div class="empty-search">
                No mentors match that search yet. Try a different topic.
              </div>
            }
          </section>
        </main>
      } @else if (view === 'meetings') {
        <main class="portal-main page-width">
          <section class="workspace-heading">
            <div>
              <p class="eyebrow"><span></span> Make time for what matters</p>
              <h1>Your conversations.</h1>
              <p>Everything you need for a thoughtful mentor session, in one calm place.</p>
            </div>
            @if (session()!.user.role === 'mentor') {
              <button type="button" class="button button-blue" (click)="bookingOpen.set(true)">
                Schedule a meeting <span>+</span>
              </button>
            }
          </section>

          @if (bookingOpen() && session()!.user.role === 'mentor') {
            <!-- Ментор создаёт встречу через POST /api/v1/mentor/meetings -->
            <section class="booking-panel surface">
              <div class="surface-heading">
                <div>
                  <p class="label">NEW CONVERSATION</p>
                  <h2>Schedule with intention</h2>
                </div>
                <button type="button" class="icon-button" (click)="bookingOpen.set(false)">
                  ×
                </button>
              </div>
              <div class="booking-fields">
                <label
                  >Topic<input
                    [(ngModel)]="booking.title"
                    placeholder="What would you like to explore?" /></label
                ><label
                  >Date and time<input
                    type="datetime-local"
                    [(ngModel)]="booking.startTime"
                    [min]="minDateTime()" /></label
                ><label
                  >Meeting format<select [(ngModel)]="booking.mode">
                    <option value="online">Online</option>
                    <option value="campus">On campus</option>
                  </select></label
                >
              </div>
              <div class="booking-footer">
                <p>Choose a future time. Your group will see the meeting in their workspace.</p>
                <button
                  type="button"
                  class="button button-blue"
                  [disabled]="loading()"
                  (click)="createMeeting()"
                >
                  Confirm meeting <span>→</span>
                </button>
              </div>
            </section>
          }

          <section class="meeting-layout">
            <div class="surface agenda">
              <div class="surface-heading">
                <div>
                  <p class="label">CALENDAR</p>
                  <h2>All sessions</h2>
                </div>
                <span class="count-badge">{{ meetings().length }}</span>
              </div>
              @for (meeting of meetings(); track meeting.id) {
                <article class="meeting-row">
                  <div class="date-tile">
                    <b>{{ (meeting.scheduled_at || meeting.start_time) | date: 'd' }}</b>
                    <span>{{ (meeting.scheduled_at || meeting.start_time) | date: 'MMM' }}</span>
                  </div>
                  <div class="meeting-detail">
                    <span class="status" [class.done]="meeting.status === 'completed'">{{
                      meeting.status
                    }}</span>
                    <h3>{{ meeting.title }}</h3>
                    <p>
                      {{ (meeting.scheduled_at || meeting.start_time) | date: 'EEEE, d MMMM · HH:mm' }} ·
                      {{ meeting.duration_minutes || 30 }} min ·
                      {{ meeting.mode === 'campus' ? 'On campus' : 'Online' }}
                    </p>
                    <p class="partner-line">
                      <span
                        class="small-avatar"
                        [style.background]="
                          meeting.mentor_color || meeting.student_color || '#385cba'
                        "
                        >{{ initials(meeting.mentor_name || meeting.student_name || '') }}</span
                      >
                      {{ session()!.user.role === 'student' ? 'with ' : 'Student: '
                      }}{{ meeting.mentor_name || meeting.student_name }}
                    </p>
                  </div>
                  @if (meeting.status === 'scheduled') {
                    @if (session()!.user.role === 'mentor') {
                      <button
                        type="button"
                        class="button button-dark small"
                        (click)="completeMeeting(meeting.id)"
                      >
                        Complete ✓
                      </button>
                    }
                  } @else {
                    <span class="completed-mark">{{
                      meeting.status === 'completed' ? '✓' : '—'
                    }}</span>
                  }
                </article>
              } @empty {
                <div class="blank-row">
                  No conversations scheduled yet. Choose a mentor and make the first move.
                </div>
              }
            </div>
            <aside class="surface prep-card">
              <p class="label">MAKE IT COUNT</p>
              <h2>Before your meeting</h2>
              <ul>
                <li>Write down the one thing you want clarity on.</li>
                <li>Bring a small example or a recent challenge.</li>
                <li>Leave with one next action, not ten.</li>
              </ul>
              <p class="tip"><b>Tip</b> The best sessions begin with a real question.</p>
            </aside>
          </section>
        </main>
      } @else if (view === 'roadmap') {
        <main class="portal-main page-width">
          @if (session()!.user.role === 'student') {
            <!-- Freshman: задачи из /api/v1/freshman/tasks -->
            <section class="workspace-heading">
              <div>
                <p class="eyebrow"><span></span> Your path, your pace</p>
                <h1>Build your momentum.</h1>
                <p>Complete your tasks and track progress with your mentor.</p>
              </div>
            </section>

            <section class="roadmap-overview">
              <div class="progress-ring">
                <svg viewBox="0 0 42 42">
                  <circle class="ring-bg" cx="21" cy="21" r="15.9"></circle>
                  <circle
                    class="ring-value"
                    cx="21"
                    cy="21"
                    r="15.9"
                    [attr.stroke-dasharray]="taskProgress() + ' ' + (100 - taskProgress())"
                  ></circle></svg
                ><b>{{ taskProgress() }}%</b>
              </div>
              <div>
                <p class="label">YOUR TASKS</p>
                <h2>Small steps. Real direction.</h2>
                <p>
                  {{ completedTaskCount() }} of {{ tasks().length }} tasks completed. Keep the promise
                  you made to yourself.
                </p>
              </div>
            </section>

            <section class="goal-list">
              @for (task of tasks(); track task.id) {
                <article class="goal-card" [class.complete]="task.status === 'approved'">
                  <div class="check-button" [class]="taskStatusClass(task.status)">
                    {{ task.status === 'approved' ? '✓' : task.status === 'submitted' ? '⏳' : task.status === 'rejected' ? '✗' : '' }}
                  </div>
                  <div class="goal-main">
                    <span>{{ task.category || 'Task' }}</span>
                    <h2>{{ task.title }}</h2>
                    @if (task.description) {
                      <p style="color:#6c7890;font-size:12px;margin:3px 0 0">{{ task.description }}</p>
                    }
                    @if (task.due_date) {
                      <p>Due: {{ task.due_date | date: 'd MMMM' }}</p>
                    }
                    @if (task.comment) {
                      <p style="color:#c0392b;font-size:11px;margin:4px 0 0">Feedback: {{ task.comment }}</p>
                    }
                  </div>
                  <div class="goal-progress">
                    <b [style.color]="taskStatusColor(task.status)">{{ task.status }}</b>
                    @if (task.status === 'pending' || task.status === 'rejected') {
                      <button
                        type="button"
                        class="button button-blue small"
                        [disabled]="loading()"
                        (click)="submitTask(task)"
                        style="margin-top:8px;"
                      >
                        Submit <span>→</span>
                      </button>
                    }
                  </div>
                </article>
              } @empty {
                <div class="empty-search">
                  No tasks assigned yet. Your mentor will add them soon.
                </div>
              }
            </section>

          } @else if (session()!.user.role === 'mentor') {
            <!-- Mentor: задачи группы из /api/v1/mentor/tasks -->
            <section class="workspace-heading">
              <div>
                <p class="eyebrow"><span></span> Your group's progress</p>
                <h1>Review tasks.</h1>
                <p>Check submitted tasks from your freshmen and give them feedback.</p>
              </div>
              <div class="filter-chips">
                <button type="button" [class.selected]="taskFilter === ''" (click)="taskFilter = ''; loadMentorTasks()">All</button>
                <button type="button" [class.selected]="taskFilter === 'submitted'" (click)="taskFilter = 'submitted'; loadMentorTasks()">Submitted</button>
                <button type="button" [class.selected]="taskFilter === 'pending'" (click)="taskFilter = 'pending'; loadMentorTasks()">Pending</button>
              </div>
            </section>

            <section class="goal-list">
              @for (task of tasks(); track task.id) {
                <article class="goal-card" [class.complete]="task.status === 'approved'">
                  <div class="check-button" [class]="taskStatusClass(task.status)">
                    {{ task.status === 'approved' ? '✓' : task.status === 'submitted' ? '⏳' : '' }}
                  </div>
                  <div class="goal-main" style="flex:1">
                    <span>{{ task.category || 'Task' }}</span>
                    <h2>{{ task.title }}</h2>
                    @if (task.description) {
                      <p style="color:#6c7890;font-size:12px;margin:3px 0 0">{{ task.description }}</p>
                    }
                    @if (task.proof_url) {
                      <p style="font-size:11px;margin:4px 0 0"><a [href]="task.proof_url" target="_blank" style="color:#3474ed">View submission →</a></p>
                    }
                  </div>
                  @if (task.status === 'submitted') {
                    <div style="display:flex;gap:8px;align-items:center;">
                      <button
                        type="button"
                        class="text-button cancel"
                        (click)="rejectTask(task.id)"
                      >Reject</button>
                      <button
                        type="button"
                        class="button button-dark small"
                        (click)="approveTask(task.id)"
                      >Approve ✓</button>
                    </div>
                  } @else {
                    <span class="status" [class.done]="task.status === 'approved'">{{ task.status }}</span>
                  }
                </article>
              } @empty {
                <div class="empty-search">
                  No tasks match the filter.
                </div>
              }
            </section>

          } @else {
            <!-- Head/Advisor: общий обзор -->
            <section class="access-layout">
              <div class="empty-state">
                <span class="large-icon">⌁</span>
                <p class="eyebrow"><span></span> Program overview</p>
                <h1>Track the community's progress.</h1>
                <p>
                  Monitor task completion rates and mentor activity from the Insights dashboard.
                </p>
                <a routerLink="/admin" class="button button-blue"
                  >Open insights <span>→</span></a
                >
              </div>
            </section>
          }
        </main>
      } @else if (view === 'admin') {
        <main class="portal-main page-width">
          @if (session()!.user.role === 'admin') {
            <!-- Head: аналитика через /api/v1/head/dashboard -->
            <section class="workspace-heading">
              <div>
                <p class="eyebrow"><span></span> Program pulse</p>
                <h1>Community, at a glance.</h1>
                <p>
                  Keep the mentorship program welcoming, responsive and growing in the right
                  direction.
                </p>
              </div>
            </section>
            <section class="metric-grid admin-metrics">
              <article class="metric-card blue">
                <p>Active freshmen</p>
                <b>{{ adminMetrics()?.total_freshmen || 0 }}</b>
                <span>Registered in MentorHub</span>
              </article>
              <article class="metric-card">
                <p>Mentors</p>
                <b>{{ adminMetrics()?.total_mentors || 0 }}</b>
                <span>Ready to support new students</span>
              </article>
              <article class="metric-card">
                <p>Overdue tasks</p>
                <b>{{ adminMetrics()?.overdue_tasks || 0 }}</b>
                <span>Tasks past deadline</span>
              </article>
              <article class="metric-card">
                <p>Average progress</p>
                <b>{{ adminMetrics()?.avg_progress_pct || 0 }}<small>%</small></b>
                <span>Community task average</span>
              </article>
            </section>

            <!-- Список пользователей (mentors) из /api/v1/head/users?role=mentor -->
            <section class="surface applications">
              <div class="surface-heading">
                <div>
                  <p class="label">USER MANAGEMENT</p>
                  <h2>All mentors</h2>
                </div>
                <span class="count-badge">{{ adminUsers().length }}</span>
              </div>
              @for (user of adminUsers(); track user.id) {
                <article class="application-row">
                  <span class="small-avatar" [style.background]="userColor(user.id)"
                    >{{ initials((user.first_name || '') + ' ' + (user.last_name || '')) }}</span
                  >
                  <div>
                    <h3>{{ user.first_name }} {{ user.last_name }}</h3>
                    <p>{{ user.email }}</p>
                    <small><b>{{ user.role }}</b></small>
                  </div>
                  <span class="status done">active</span>
                </article>
              } @empty {
                <div class="blank-row">No users loaded yet.</div>
              }
            </section>

          } @else {
            <section class="access-layout">
              <div class="empty-state">
                <span class="large-icon">⌁</span>
                <h1>This space is for program administrators.</h1>
                <p>
                  Administrators can review mentor statistics and track the health of the
                  community.
                </p>
                <a routerLink="/dashboard" class="button button-blue"
                  >Back to my workspace <span>→</span></a
                >
              </div>
            </section>
          }
        </main>
      } @else if (view === 'profile') {
        <main class="portal-main page-width">
          <section class="profile-layout">
            <aside class="profile-aside">
              <div class="profile-orb large" [style.background]="session()!.user.avatarColor">
                {{ initials(session()!.user.fullName) }}
              </div>
              <h1>{{ session()!.user.fullName }}</h1>
              <p>{{ roleLabel(session()!.user.role) }} · KBTU</p>
              <span class="status done">Active member</span>
              <div class="profile-facts">
                <p><span>✦</span> KBTU verified</p>
                <p><span>◷</span> Member since 2026</p>
                <p><span>⌁</span> {{ session()!.user.email }}</p>
              </div>
            </aside>
            <section class="surface profile-editor">
              <div class="surface-heading">
                <div>
                  <p class="label">YOUR DETAILS</p>
                  <h2>Profile settings</h2>
                </div>
                <span class="safe-badge">✓ Private</span>
              </div>
              <label>First name<input [(ngModel)]="profileFirstName" maxlength="100" /></label>
              <label>Last name<input [(ngModel)]="profileLastName" maxlength="100" /></label>
              <label>Email address<input [value]="session()!.user.email" disabled /></label>
              <button
                type="button"
                class="button button-blue"
                [disabled]="loading()"
                (click)="saveProfile()"
              >
                Save changes <span>→</span>
              </button>
            </section>
          </section>
        </main>
      }
    </div>
  `,
})
export class PortalComponent {
  private readonly api = inject(ApiService);
  private readonly router = inject(Router);
  private readonly route = inject(ActivatedRoute);
  private readonly platformId = inject(PLATFORM_ID);

  readonly today = new Date();

  // ── State signals ──────────────────────────────────────────────────────────
  readonly session = signal<Session | null>(null);
  readonly mentors = signal<Mentor[]>([]);
  readonly meetings = signal<Meeting[]>([]);
  readonly tasks = signal<Task[]>([]);
  readonly dashboard = signal<any>(null);
  readonly adminMetrics = signal<any>(null);
  readonly adminUsers = signal<any[]>([]);
  readonly notice = signal('');
  readonly loading = signal(false);
  readonly bookingOpen = signal(false);

  view = 'home';
  mentorSearch = '';
  mentorFocus = 'All';
  profileFirstName = '';
  profileLastName = '';
  taskFilter = '';
  loginForm = { email: 'freshman@kbtu.kz', password: 'Freshman2026!' };
  booking = { title: '', startTime: '', mode: 'online' };

  // ── Computed ───────────────────────────────────────────────────────────────
  readonly completedTaskCount = computed(
    () => this.tasks().filter((t) => t.status === 'approved').length,
  );

  readonly taskProgress = computed(() => {
    const all = this.tasks();
    if (!all.length) return 0;
    const done = all.filter((t) => t.status === 'approved').length;
    return Math.round((done / all.length) * 100);
  });

  readonly pendingTaskCount = computed(
    () => this.tasks().filter((t) => t.status === 'pending' || t.status === 'rejected').length,
  );

  constructor() {
    this.route.data.subscribe((data) => {
      this.view = String(data['view'] || 'home');
      if (isPlatformBrowser(this.platformId)) this.load();
    });
  }

  // ── Session ────────────────────────────────────────────────────────────────

  private readSession(): Session | null {
    try {
      return JSON.parse(localStorage.getItem('mentorhub_session') || 'null') as Session | null;
    } catch {
      return null;
    }
  }

  // ── Load ───────────────────────────────────────────────────────────────────

  private load(): void {
    this.session.set(this.readSession());
    // Менторов загружаем всегда (публичная страница)
    this.loadMentors();

    const active = this.session();
    if (!active) return;

    const [first, ...rest] = active.user.fullName.split(' ');
    this.profileFirstName = first || '';
    this.profileLastName = rest.join(' ');

    const role = active.user.role;
    const basePath = this.api.roleBasePath(role);

    // Dashboard
    this.api.get<any>(`${basePath}/dashboard`).subscribe({
      next: (result) => this.dashboard.set(result?.data ?? result),
      error: () => undefined,
    });

    // Meetings (для dashboard и meetings views)
    if (['dashboard', 'meetings'].includes(this.view)) {
      this.loadMeetings();
    }

    // Tasks (freshman → /freshman/tasks, mentor → /mentor/tasks)
    if (['dashboard', 'roadmap'].includes(this.view)) {
      if (role === 'student') this.loadFreshmanTasks();
      if (role === 'mentor') this.loadMentorTasks();
    }

    // Admin metrics (head)
    if (this.view === 'admin' && role === 'admin') {
      this.loadAdminData();
    }
  }

  private loadMentors(): void {
    // Используем список пользователей с ролью mentor
    // Используем новый защищенный эндпоинт /mentors, доступный всем авторизованным пользователям
    this.api.get<any>('mentors?role=mentor&per_page=50').subscribe({
      next: (result) => {
        const users: any[] = result?.data ?? result ?? [];
        const mentors: Mentor[] = users.map((u: any) => ({
          id: u.id,
          name: `${u.first_name || ''} ${u.last_name || ''}`.trim(),
          color: avatarColor(u.id),
          program: u.program || undefined,
          graduationYear: u.graduation_year || undefined,
          title: u.title || undefined,
          bio: u.bio || undefined,
          skills: u.skills || [],
          languages: u.languages || [],
          rating: u.rating || undefined,
          sessionsCount: u.sessions_count || 0,
          capacity: u.capacity || undefined,
          availability: u.availability || 'Available',
        }));
        this.mentors.set(mentors);
      },
      error: () => this.mentors.set([]),
    });
  }

  private loadMeetings(): void {
    const role = this.session()?.user.role;
    const basePath = this.api.roleBasePath(role ?? 'student');
    this.api.get<any>(`${basePath}/meetings`).subscribe({
      next: (result) => {
        const data = result?.data ?? result;
        const meetingsList = Array.isArray(data) ? data : data?.meetings ?? [];
        const mapped = meetingsList.map((m: any) => {
          let mentorName = '';
          if (m.mentor) {
            mentorName = `${m.mentor.first_name || ''} ${m.mentor.last_name || ''}`.trim();
          }
          return {
            id: m.id,
            title: m.title,
            scheduled_at: m.scheduled_at,
            duration_minutes: m.duration_minutes || 45,
            mode: m.mode || 'online',
            status: m.held ? 'completed' : 'scheduled',
            mentor_name: mentorName,
            mentor_color: m.mentor ? avatarColor(m.mentor.id) : '#385cba',
            notes: m.notes || '',
          };
        });
        this.meetings.set(mapped);
      },
      error: () => undefined,
    });
  }

  private loadFreshmanTasks(): void {
    this.api.get<any>('freshman/tasks').subscribe({
      next: (result) => {
        const data = result?.data ?? result;
        const tasks: Task[] = (Array.isArray(data) ? data : data?.tasks ?? []).map((t: any) => ({
          id: t.id,
          title: t.template?.title || 'MentorHub Task',
          description: t.template?.description || '',
          status: t.status,
          due_date: t.due_date,
          proof_url: t.proof_url || undefined,
          comment: t.comment || undefined,
          category: t.category || 'Academic',
        }));
        this.tasks.set(tasks);
      },
      error: () => undefined,
    });
  }

  loadMentorTasks(): void {
    const filter = this.taskFilter ? `?status=${this.taskFilter}` : '';
    this.api.get<any>(`mentor/tasks${filter}`).subscribe({
      next: (result) => {
        const data = result?.data ?? result;
        const tasks: Task[] = (Array.isArray(data) ? data : data?.tasks ?? []).map((t: any) => ({
          id: t.id,
          title: t.template?.title || 'MentorHub Task',
          description: t.template?.description || '',
          status: t.status,
          due_date: t.due_date,
          proof_url: t.proof_url || undefined,
          comment: t.comment || undefined,
          category: t.category || 'Academic',
        }));
        this.tasks.set(tasks);
      },
      error: () => undefined,
    });
  }

  private loadAdminData(): void {
    this.api.get<any>('head/dashboard').subscribe({
      next: (result) => this.adminMetrics.set(result?.data ?? result),
      error: () => undefined,
    });
    this.api.get<any>('head/users?per_page=50').subscribe({
      next: (result) => {
        const data = result?.data ?? result;
        this.adminUsers.set(Array.isArray(data) ? data : data?.users ?? []);
      },
      error: () => undefined,
    });
  }

  // ── Auth ───────────────────────────────────────────────────────────────────

  login(): void {
    this.loading.set(true);
    this.api.login(this.loginForm.email, this.loginForm.password).subscribe({
      next: (session) => {
        localStorage.setItem('mentorhub_session', JSON.stringify(session));
        this.session.set(session);
        this.loading.set(false);
        this.notice.set(
          `Welcome, ${this.firstName(session.user.fullName)}. Your workspace is ready.`,
        );
        this.router.navigate(['/dashboard']);
      },
      error: (error) => {
        this.loading.set(false);
        this.notice.set(error?.error?.message || error?.error?.error || 'Invalid email or password.');
      },
    });
  }

  useDemo(role: Role): void {
    const accounts: Record<Role, [string, string]> = {
      student: ['freshman@kbtu.kz', 'Freshman2026!'],
      mentor: ['mentor@kbtu.kz', 'Mentor2026!'],
      admin: ['head@kbtu.kz', 'Head2026!'],
    };
    [this.loginForm.email, this.loginForm.password] = accounts[role];
    this.login();
  }

  logout(): void {
    // POST /api/v1/auth/logout (stateless — просто чистим локально)
    this.api.post<any>('auth/logout', {}).subscribe({ error: () => undefined });
    localStorage.removeItem('mentorhub_session');
    this.session.set(null);
    this.notice.set('You have been signed out safely.');
    this.router.navigate(['/']);
  }

  // ── Mentors ────────────────────────────────────────────────────────────────

  requestMentor(mentor: Mentor): void {
    const active = this.session();
    if (!active) {
      this.router.navigate(['/login']);
      return;
    }
    if (active.user.role !== 'student') {
      this.notice.set('Connection requests are available from a student account.');
      return;
    }
    this.notice.set(`Request to connect with ${mentor.name} noted. Contact your head to assign a mentor.`);
  }

  // ── Meetings ───────────────────────────────────────────────────────────────

  createMeeting(): void {
    this.loading.set(true);
    // Mentor creates meeting: POST /api/v1/mentor/meetings
    this.api
      .post<any>('mentor/meetings', {
        title: this.booking.title,
        scheduled_at: this.booking.startTime,
        mode: this.booking.mode,
        duration_minutes: 30,
      })
      .subscribe({
        next: () => {
          this.loading.set(false);
          this.bookingOpen.set(false);
          this.booking.title = '';
          this.notice.set('Your meeting is on the calendar. Great first step!');
          this.loadMeetings();
        },
        error: (error) => {
          this.loading.set(false);
          this.handleError(error);
        },
      });
  }

  completeMeeting(id: string): void {
    // PUT /api/v1/mentor/meetings/{id}/complete
    this.api.put<any>(`mentor/meetings/${id}/complete`, { notes: '' }).subscribe({
      next: () => {
        this.notice.set('Meeting marked as completed.');
        this.loadMeetings();
      },
      error: (error) => this.handleError(error),
    });
  }

  // ── Tasks ──────────────────────────────────────────────────────────────────

  submitTask(task: Task): void {
    const proofUrl = prompt('Enter a valid submission URL (must start with http:// or https://):');
    if (!proofUrl) return;
    this.loading.set(true);
    // PUT /api/v1/freshman/tasks/{id}/submit
    this.api.put<any>(`freshman/tasks/${task.id}/submit`, { proof_url: proofUrl }).subscribe({
      next: () => {
        this.loading.set(false);
        this.notice.set('Task submitted! Your mentor will review it.');
        this.loadFreshmanTasks();
      },
      error: (error) => {
        this.loading.set(false);
        this.handleError(error);
      },
    });
  }

  approveTask(taskId: string): void {
    // PUT /api/v1/mentor/tasks/{id}/approve
    this.api.put<any>(`mentor/tasks/${taskId}/approve`, { comment: '' }).subscribe({
      next: () => {
        this.notice.set('Task approved.');
        this.loadMentorTasks();
      },
      error: (error) => this.handleError(error),
    });
  }

  rejectTask(taskId: string): void {
    const comment = prompt('Reason for rejection (optional):') || '';
    // PUT /api/v1/mentor/tasks/{id}/reject
    this.api.put<any>(`mentor/tasks/${taskId}/reject`, { comment }).subscribe({
      next: () => {
        this.notice.set('Task rejected. Student has been notified.');
        this.loadMentorTasks();
      },
      error: (error) => this.handleError(error),
    });
  }

  // ── Profile ────────────────────────────────────────────────────────────────

  saveProfile(): void {
    this.loading.set(true);
    this.api
      .patch<any>('auth/me', {
        first_name: this.profileFirstName,
        last_name: this.profileLastName,
      })
      .subscribe({
        next: (result) => {
          const me = result?.data ?? result;
          const current = this.session()!;
          const updated: Session = {
            ...current,
            user: {
              ...current.user,
              fullName: `${me.first_name || this.profileFirstName} ${me.last_name || this.profileLastName}`.trim(),
            },
          };
          localStorage.setItem('mentorhub_session', JSON.stringify(updated));
          this.session.set(updated);
          this.loading.set(false);
          this.notice.set('Your profile has been updated.');
        },
        error: (error) => {
          this.loading.set(false);
          this.handleError(error);
        },
      });
  }

  // ── Helpers ────────────────────────────────────────────────────────────────

  visibleMentors(): Mentor[] {
    const search = this.mentorSearch.toLowerCase().trim();
    return this.mentors().filter((mentor) => {
      const allText = [mentor.name, mentor.title, mentor.program, ...(mentor.skills || [])]
        .join(' ')
        .toLowerCase();
      const matchesSearch = !search || allText.includes(search);
      const matchesFocus =
        this.mentorFocus === 'All' ||
        (this.mentorFocus === 'Career' && /career|cv|product/i.test(allText)) ||
        (this.mentorFocus === 'Academic' && /algorithm|web|data/i.test(allText)) ||
        (this.mentorFocus === 'Campus' && /campus|community|time/i.test(allText));
      return matchesSearch && matchesFocus;
    });
  }

  greeting(): string {
    const hour = new Date().getHours();
    return hour < 12 ? 'Good morning' : hour < 18 ? 'Good afternoon' : 'Good evening';
  }

  dashboardHeadline(): string {
    const role = this.session()?.user.role;
    return role === 'mentor'
      ? 'Your experience is someone\'s shortcut.'
      : role === 'admin'
        ? 'A stronger start for every student.'
        : 'You\'re building more than a semester.';
  }

  dashboardSubline(): string {
    const role = this.session()?.user.role;
    return role === 'mentor'
      ? 'Every thoughtful conversation can change a student\'s first year.'
      : role === 'admin'
        ? 'See how the MentorHub community is connecting today.'
        : 'You\'re building a circle of support, a sense of direction and a future you can trust.';
  }

  upcomingCount(): number {
    return this.meetings().filter(
      (m) => m.status === 'scheduled' && new Date(m.scheduled_at || m.start_time || '') > new Date(),
    ).length;
  }

  minDateTime(): string {
    const date = new Date(Date.now() + 5 * 60000);
    date.setSeconds(0, 0);
    return new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 16);
  }

  initials(name: string): string {
    return name
      .split(' ')
      .filter(Boolean)
      .slice(0, 2)
      .map((part) => part[0])
      .join('')
      .toUpperCase();
  }

  firstName(name: string): string {
    return name.split(' ')[0];
  }

  roleLabel(role: Role): string {
    return role === 'student' ? 'Freshman' : role === 'mentor' ? 'Mentor' : 'Head / Administrator';
  }

  taskStatusClass(status: string): string {
    if (status === 'approved') return 'check-button approved';
    if (status === 'submitted') return 'check-button submitted';
    if (status === 'rejected') return 'check-button rejected';
    return 'check-button';
  }

  taskStatusColor(status: string): string {
    if (status === 'approved') return '#388b5e';
    if (status === 'submitted') return '#aa7d24';
    if (status === 'rejected') return '#bf6c66';
    return '#7d899f';
  }

  userColor(id: string): string {
    return avatarColor(id);
  }

  private handleError(error: any, showNotice = true): void {
    if (error?.status === 401) {
      localStorage.removeItem('mentorhub_session');
      this.session.set(null);
    }
    if (showNotice)
      this.notice.set(
        error?.error?.message || error?.error?.error || 'We could not complete that action. Please try again.',
      );
  }
}
