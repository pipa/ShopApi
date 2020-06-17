# alpine Go...
# @date 06/2020
FROM golang

ARG app_env
ENV APP_ENV $app_env
ENV PORT 4200
ENV WEBROOT /go/src/github.com/pipa/ShopApi

WORKDIR ${WEBROOT}
ADD ./ ${WEBROOT}

# RUN go get github.com/jinzhu/gorm
RUN go get ./
RUN go build

CMD if [ ${APP_ENV} = production ]; \
  then \
  app; \
  else \
  go get github.com/pilu/fresh && \
  fresh; \
  fi

EXPOSE ${PORT}
