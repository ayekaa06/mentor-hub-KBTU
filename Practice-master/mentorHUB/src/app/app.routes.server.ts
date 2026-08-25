import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  {
    path: '**',
    // Personal workspaces depend on the authenticated user, so render them on
    // request instead of generating static pages at build time.
    renderMode: RenderMode.Server
  }
];
