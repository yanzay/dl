FROM alpine

RUN apk update && \
    apk upgrade && \
    apk add ffmpeg py-pip ca-certificates && \
    rm -rf /var/cache/apk/*

RUN pip install youtube-dl

COPY dl /dl

CMD [ "/dl" ]
