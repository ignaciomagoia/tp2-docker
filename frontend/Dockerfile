# Etapa 1: build de React
FROM node:18 AS builder

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .
RUN npm run build

# Etapa 2: servir con Nginx
FROM nginx:alpine

# Copiar build al servidor web
COPY --from=builder /app/build /usr/share/nginx/html

# Exponer el puerto 80 (puerto por defecto de Nginx)
EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
