cat << EOF > ./backend/config/config.yml
debug: false
app_name: social-network
app_env: ${ENV}

http:
  port: ${BACKEND_INTERNAL_PORT}
  domain: ${APP_SUBDOMAIN}.${DOMAIN}

mysql:
  host: ${MYSQL_HOST}
  port: ${MYSQL_PORT}
  user: ${MYSQL_USER}
  password: ${MYSQL_PASSWORD}
  database: ${MYSQL_DATABASE}

auth:
  token-secret: ${AUTH_SECRET}
  accessTTL: 1h
  refresTTL: 30h
EOF
