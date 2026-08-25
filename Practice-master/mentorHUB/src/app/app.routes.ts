import { Routes } from '@angular/router';
import { PortalComponent } from './portal.component';

export const routes: Routes = [
  { path: '', component: PortalComponent, data: { view: 'home' }, title: 'MentorHub KBTU' },
  { path: 'login', component: PortalComponent, data: { view: 'login' }, title: 'Sign in · MentorHub KBTU' },
  { path: 'dashboard', component: PortalComponent, data: { view: 'dashboard' }, title: 'Dashboard · MentorHub KBTU' },
  { path: 'mentors', component: PortalComponent, data: { view: 'mentors' }, title: 'Find a mentor · MentorHub KBTU' },
  { path: 'meetings', component: PortalComponent, data: { view: 'meetings' }, title: 'Meetings · MentorHub KBTU' },
  { path: 'roadmap', component: PortalComponent, data: { view: 'roadmap' }, title: 'Roadmap · MentorHub KBTU' },
  { path: 'admin', component: PortalComponent, data: { view: 'admin' }, title: 'Administration · MentorHub KBTU' },
  { path: 'profile', component: PortalComponent, data: { view: 'profile' }, title: 'Profile · MentorHub KBTU' },
  { path: '**', redirectTo: '' },
];
