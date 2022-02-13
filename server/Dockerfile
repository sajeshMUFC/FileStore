#
# HIGHLY CONFIDENTIAL
#
# Copyright (c) 2020, 2021 Wabtec Corporation. All rights reserved.
#
# NOTICE: This file contains material that is confidential and proprietary to
# Wabtec and/or other developers and may be protected by patents, copyright, and/or
# trade secrets. No license is granted under any intellectual property rights
# of Wabtec except as may be provided in a duly executed agreement with Wabtec.
# Any unauthorized reproduction or distribution of material from this file is
# expressly prohibited. This source code is and remains the property of Wabtec.
#


FROM golang:1.17-alpine3.14 AS builder
RUN apk update \
    && apk upgrade \
    && apk add --no-cache \
    git \
    build-base \
    openssh

WORKDIR /build
COPY . .
RUN go build -o app .

# Now build the production container
FROM alpine:3.14
RUN apk update \
    && apk add ca-certificates \
    && adduser -S -D -H -h /app appuser
USER appuser
COPY --from=builder /build/ /app/
COPY --from=builder /build/saved /app/saved

WORKDIR /app

CMD ["./app"]
