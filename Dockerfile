FROM avcosystems/golang-node:1.13.0

WORKDIR /src

COPY . .
RUN make deps
RUN make all

EXPOSE 9001

CMD bin/server
