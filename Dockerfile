FROM scratch
COPY awsAuth /usr/bin/awsAuth
ENTRYPOINT ["/usr/bin/awsAuth"]