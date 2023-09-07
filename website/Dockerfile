FROM node:20.5-alpine

COPY . /app
WORKDIR /app
RUN npm ci
RUN npm run build
CMD npm run serve -- --port 80 --host 0.0.0.0
