FROM node:16 as builder
WORKDIR /opt
COPY . .
RUN yarn install && yarn build

FROM nginx:1.21.3-alpine
COPY --from=builder /opt/build /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]