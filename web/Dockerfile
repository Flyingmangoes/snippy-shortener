# Set docker base image 
# Set root dir
# Install required dependency
# Copy all files to container
# Inject environment variable

FROM node:20-alpine AS builder
WORKDIR /app
COPY package*.json ./
RUN npm ci -x
COPY . .
ARG VITE_BACKEND_URL
ENV VITE_BACKEND_URL=$VITE_BACKEND_URL   
RUN npm run build

# Run the app
FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
EXPOSE 80
CMD [ "nginx", "-g", "daemon off;" ]
