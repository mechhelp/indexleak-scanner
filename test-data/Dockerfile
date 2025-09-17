# Use the official Nginx image from Docker Hub
FROM nginx:alpine

# Copy the test data to the Nginx web root directory
COPY ./zPeter /usr/share/nginx/html/zPeter

# Copy the custom Nginx configuration
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80
