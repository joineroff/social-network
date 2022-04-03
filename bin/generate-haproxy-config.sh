cat << EOF > ./haproxy/haproxy.cfg
global
    user haproxy
    group haproxy
    daemon

defaults
    log     global
    mode    http
    option  httplog
    option  dontlognull

frontend stats
    bind 0.0.0.0:8888
    stats enable
    stats uri /
    stats refresh 10s
    stats auth ${HAPROXY_STATS_USER}:${HAPROXY_STATS_PASSWORD}

frontend http-in
    bind 0.0.0.0:80
    acl sn hdr(host) -i ${APP_SUBDOMAIN}.${DOMAIN}
    acl default hdr_end(host) -i .${DOMAIN}

    ## figure out which one to use
    use_backend sn_cluster if sn
    use_backend direct_forward if default

backend sn_cluster
    option forwardfor
    server node1 backend:${BACKEND_INTERNAL_PORT}

# basically forwarding to the source itself
backend direct_forward
    option httpclose
    option http_proxy
EOF
