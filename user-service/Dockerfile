# Build stage
FROM node:20-alpine AS builder

WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

# Production stage
FROM node:20-alpine

WORKDIR /app
COPY --from=builder /app/dist ./dist
COPY package*.json ./
RUN npm install --only=production && \
    npm cache clean --force

ENV NODE_ENV=production
EXPOSE 3000
USER node

CMD ["node", "dist/src/main.js"]
