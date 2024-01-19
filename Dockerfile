FROM alpine:3.19

# 国内源
RUN echo http://mirrors.aliyun.com/alpine/v3.19/main/ > /etc/apk/repositories

# +8时区
RUN apk update && apk add tzdata ca-certificates bash
RUN rm -rf /etc/localtime && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo "Asia/Shanghai" > /etc/timezone

# git 信息
ARG branch=0
ARG commit=0
LABEL branch=$branch commit=$commit

ENV WORKDIR /app
ADD ./tmp/main $WORKDIR/main
RUN chmod +x $WORKDIR/main

WORKDIR $WORKDIR
CMD ./main -config/config.yaml
