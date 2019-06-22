FROM golang:1.12.4-alpine3.9
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN apk add --no-cache \ 
    git \
    tzdata \
    gcc \
    g++ && \
    cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
    echo Asia/Shanghai > /etc/timezone && \
    apk del tzdata
RUN go get -u github.com/beego/bee

# RUN git config --global credential.helper store

CMD [ "/bin/sh" ]