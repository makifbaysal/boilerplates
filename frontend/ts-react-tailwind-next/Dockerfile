FROM node:24-alpine AS build
WORKDIR /src
COPY package.json package-lock.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM node:24-alpine
RUN addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=build /src/public ./public
COPY --from=build /src/.next/standalone ./
COPY --from=build /src/.next/static ./.next/static
# NEXT_PUBLIC_API_BASE_URL: see .ai/architecture.md#config — deployer
# sets it (baked in at build time — it's a NEXT_PUBLIC_ var, inlined
# into the client bundle, not read at runtime).
USER app:app
EXPOSE 3000
ENTRYPOINT ["node", "server.js"]
