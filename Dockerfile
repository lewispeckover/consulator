FROM lewispeckover/base:3.5
COPY ./docker/ /
ENTRYPOINT ["/entrypoint.sh"]
ENV VERSION=0.4.3
ADD https://github.com/lewispeckover/consulator/releases/download/${VERSION}/consulator_${VERSION}_linux_amd64 /bin/consulator
RUN chmod +x /bin/consulator
