FROM node:14-alpine

WORKDIR ./build

COPY package.json yarn.lock ./
RUN yarn install

COPY ./app.js .

CMD ["node", "app.js"]