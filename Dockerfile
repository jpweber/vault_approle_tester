FROM alpine
RUN apk --update upgrade && \
    apk add curl ca-certificates && \
    update-ca-certificates && \
    rm -rf /var/cache/apk/*
    
COPY vault_approle_test_linux /vault_approle_test
EXPOSE 8080
RUN chmod +x /vault_approle_test
CMD sh