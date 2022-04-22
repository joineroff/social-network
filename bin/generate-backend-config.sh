cat << EOF > ./backend/config/config.yml
debug: false
app_name: social-network
app_env: ${ENV}

http:
  port: ${BACKEND_INTERNAL_PORT}
  domain: ${APP_SUBDOMAIN}.${DOMAIN}

mysql:
  master:
      host: ${MYSQL_HOST}
      port: ${MYSQL_PORT}
      user: ${MYSQL_USER}
      password: ${MYSQL_PASSWORD}
      database: ${MYSQL_DATABASE}
  slaves:
      - host: ${MYSQL_SLAVE1_HOST}
        port: ${MYSQL_SLAVE1_PORT}
        user: ${MYSQL_SLAVE1_USER}
        password: ${MYSQL_SLAVE1_PASSWORD}
        database: ${MYSQL_SLAVE1_DATABASE}
      - host: ${MYSQL_SLAVE2_HOST}
        port: ${MYSQL_SLAVE2_PORT}
        user: ${MYSQL_SLAVE2_USER}
        password: ${MYSQL_SLAVE2_PASSWORD}
        database: ${MYSQL_SLAVE2_DATABASE}

auth:
  token-secret: ${AUTH_SECRET}
  accessTTL: 1h
  refresTTL: 30h
EOF
