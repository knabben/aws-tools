FROM golang:1.7

ENV GOPATH /root/.go
ENV PROJECT ${GOPATH}/src/github.com/knabben/aws-tools/

WORKDIR ${PROJECT}

RUN mkdir -p ${PROJECT}
RUN mkdir -p ${GOPATH}/bin
ENV PATH=${PATH}:${GOPATH}/bin

ADD . ${PROJECT}

RUN curl https://glide.sh/get | sh
RUN glide install
RUN make compile

ENV REDIS_URL ${REDIS_URL:-localhost:6379}
ENV REDIS_PASS ${REDIS_PASS:-pass}

CMD ["sh", "-c", "${PROJECT}/main serve -r ${REDIS_URL} -p ${REDIS_PASS}"]
