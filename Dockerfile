FROM nginx
COPY index.html /usr/share/nginx/html
COPY css /usr/share/nginx/html/css
COPY js /usr/share/nginx/html/js
COPY scss /usr/share/nginx/html/scss
COPY fonts /usr/share/nginx/html/fonts
COPY img /usr/share/nginx/html/img
