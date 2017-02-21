FROM lewispeckover/base:3.5
COPY ./docker/ /
ENTRYPOINT ["/entrypoint.sh"]
ADD https://github.com/lewispeckover/consulator/releases/download/0.1.8/consulator_0.1.8_linux_amd64 /bin/consulator
RUN chmod +x /bin/consulator
