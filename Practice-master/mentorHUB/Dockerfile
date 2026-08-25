# ── Stage 1: Build ──────────────────────────────────────────────────────────
FROM node:20-alpine AS builder

WORKDIR /app

# Cache dependencies
COPY package*.json ./
RUN npm ci

# Copy sources and build
COPY . .
RUN npx ng build --configuration production --base-href /mentorhub/

# ── Stage 2: Serve ───────────────────────────────────────────────────────────
FROM nginx:alpine

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy build output to nginx html directory
COPY --from=builder /app/dist/mentorHUB/browser /usr/share/nginx/html

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
