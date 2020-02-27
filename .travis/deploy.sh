ssh -o UserKnownHostsFile=/dev/null -o StrictHostKeyChecking=no ${USER}@${HOST} "\
    docker pull steden/gosocks5:0.1 && \
    docker restart gosocks5 || docker run -d -p 1080:1080 --name gosocks5 steden/gosocks5:0.1 -username=${SOCKS_USER} -password=${SOCKS_PASS}"