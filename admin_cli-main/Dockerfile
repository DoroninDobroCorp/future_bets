FROM node:20-alpine
WORKDIR /opt/app
ADD package.json package.json
RUN npm install
ADD . .
ENV NODE_ENV=production

RUN npm run build
RUN npm prune --production

EXPOSE 31562

CMD npm start
