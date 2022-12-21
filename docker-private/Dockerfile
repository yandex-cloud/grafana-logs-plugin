FROM bitnami/git AS magesrc

RUN git clone https://github.com/magefile/mage /mage && \
    ls -ld /mage

FROM golang:1.19-alpine AS gobuild

COPY --from=magesrc /mage/ /mage/

RUN cd /mage && \
    go run bootstrap.go && \
    which mage

RUN mkdir -p /yc-logging

ADD go.mod go.sum /yc-logging/
ADD Magefile.go /yc-logging/
ADD pkg/ /yc-logging/pkg
ADD src/plugin.json /yc-logging/src/plugin.json

RUN cd /yc-logging && \
    ls -la && \
    go mod vendor && \
    mage -v

FROM node:16 as nodebuild

RUN mkdir -p /yc-logging

ADD README.md /yc-logging/
ADD CHANGELOG.md /yc-logging/
ADD jest.config.js /yc-logging/
ADD LICENSE /yc-logging/
ADD package.json /yc-logging/
ADD src/ /yc-logging/src
ADD .config/ /yc-logging/.config
ADD .eslintrc /yc-logging/
ADD tsconfig.json /yc-logging/
ADD yarn.lock /yc-logging/


RUN cd /yc-logging && \
    ls -la && \
    yarn install --pure-lockfile &&\
    yarn build

COPY --from=gobuild /yc-logging/dist/ /yc-logging/dist/

FROM grafana/grafana
COPY --from=nodebuild /yc-logging/dist/ /data/plugins/yandexcloud-logging-datasource-unsigned
