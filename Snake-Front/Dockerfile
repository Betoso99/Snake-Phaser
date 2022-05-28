# syntax=docker/dockerfile:1
FROM node:12-alpine

WORKDIR /src/app

COPY package.json .

RUN npm install

COPY . .

EXPOSE 8000

CMD ["npm", "start"]