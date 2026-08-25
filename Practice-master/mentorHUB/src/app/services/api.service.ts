import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Observable, switchMap } from 'rxjs';
import { map } from 'rxjs/operators';

// ── Роли ──────────────────────────────────────────────────────────────────────
// Go backend использует: 'head' | 'advisor' | 'mentor' | 'freshman'
// Angular UI использует: 'student' | 'mentor' | 'admin'
// Маппинг: freshman → student, head → admin, mentor → mentor, advisor → admin
export type GoRole = 'head' | 'advisor' | 'mentor' | 'freshman';
export type Role = 'student' | 'mentor' | 'admin';

export interface SessionUser {
  id: string;          // UUID (Go использует UUID, не integer)
  fullName: string;
  email: string;
  role: Role;
  avatarColor: string;
}
export interface Session { token: string; refreshToken: string; user: SessionUser; }

// ── Go API response types ──────────────────────────────────────────────────────
interface GoTokenResponse {
  access_token: string;
  refresh_token: string;
  token_type: string;
  expires_in: number;
}

interface GoMeResponse {
  id: string;
  email: string;
  role: GoRole;
  first_name: string;
  last_name: string;
  avatar_url?: string;
  is_active: boolean;
}

// ── Helpers ──────────────────────────────────────────────────────────────────

/** Маппинг Go-роли на UI-роль */
export function mapGoRole(role: GoRole): Role {
  if (role === 'freshman') return 'student';
  if (role === 'head' || role === 'advisor') return 'admin';
  return 'mentor';
}

/** Генерация цвета аватара по ID пользователя */
export function avatarColor(id: string): string {
  const colors = ['#3578F6', '#8B5CF6', '#EC4899', '#00AFA0', '#F59E0B', '#EF4444', '#10B981', '#6366F1'];
  let hash = 0;
  for (let i = 0; i < id.length; i++) hash = id.charCodeAt(i) + ((hash << 5) - hash);
  return colors[Math.abs(hash) % colors.length];
}

@Injectable({ providedIn: 'root' })
export class ApiService {
  private readonly http = inject(HttpClient);
  // Все API запросы идут на Go backend через api/v1 (относительно base href)
  private readonly baseUrl = 'api/v1';

  // ── HTTP methods ────────────────────────────────────────────────────────────

  get<T>(path: string): Observable<T> {
    return this.http.get<T>(`${this.baseUrl}/${path}`, { headers: this.authHeaders() });
  }

  post<T>(path: string, body: unknown): Observable<T> {
    return this.http.post<T>(`${this.baseUrl}/${path}`, body, { headers: this.authHeaders() });
  }

  patch<T>(path: string, body: unknown): Observable<T> {
    return this.http.patch<T>(`${this.baseUrl}/${path}`, body, { headers: this.authHeaders() });
  }

  put<T>(path: string, body: unknown = {}): Observable<T> {
    return this.http.put<T>(`${this.baseUrl}/${path}`, body, { headers: this.authHeaders() });
  }

  delete<T>(path: string): Observable<T> {
    return this.http.delete<T>(`${this.baseUrl}/${path}`, { headers: this.authHeaders() });
  }

  // ── Auth ────────────────────────────────────────────────────────────────────

  /**
   * Логин через Go API:
   * 1. POST /api/v1/auth/login → { access_token, refresh_token }
   * 2. GET  /api/v1/auth/me   → { id, email, role, first_name, last_name }
   * 3. Собрать Session совместимый с Angular UI
   */
  login(email: string, password: string): Observable<Session> {
    return this.http
      .post<{ success: boolean; data: GoTokenResponse }>(`${this.baseUrl}/auth/login`, { email, password })
      .pipe(
        switchMap((resp) => {
          console.log('Login API Response:', resp);
          const tokens = resp.data;
          console.log('Extracted tokens:', tokens);
          const headers = new HttpHeaders({ Authorization: `Bearer ${tokens.access_token}` });
          console.log('Requesting /auth/me with headers:', headers.get('Authorization'));
          return this.http.get<{ success: boolean; data: GoMeResponse }>(`${this.baseUrl}/auth/me`, { headers }).pipe(
            map((meResp) => {
              console.log('/auth/me response:', meResp);
              const me = meResp.data;
              const session: Session = {
                token: tokens.access_token,
                refreshToken: tokens.refresh_token,
                user: {
                  id: me.id,
                  fullName: `${me.first_name} ${me.last_name}`.trim(),
                  email: me.email,
                  role: mapGoRole(me.role),
                  avatarColor: avatarColor(me.id),
                },
              };
              return session;
            }),
          );
        }),
      );
  }

  // ── Role-based API path helpers ─────────────────────────────────────────────

  /**
   * Возвращает базовый путь API в зависимости от роли пользователя.
   * freshman → /api/v1/freshman/...
   * mentor   → /api/v1/mentor/...
   * head     → /api/v1/head/...
   */
  roleBasePath(role: Role): string {
    if (role === 'student') return 'freshman';
    if (role === 'admin') return 'head';
    return 'mentor';
  }

  // ── Private ─────────────────────────────────────────────────────────────────

  private authHeaders(): HttpHeaders {
    const raw = typeof localStorage === 'undefined' ? null : localStorage.getItem('mentorhub_session');
    if (!raw) return new HttpHeaders({ 'Content-Type': 'application/json' });
    try {
      const session = JSON.parse(raw) as Session;
      return new HttpHeaders({
        Authorization: `Bearer ${session.token}`,
        'Content-Type': 'application/json',
      });
    } catch {
      return new HttpHeaders({ 'Content-Type': 'application/json' });
    }
  }
}
