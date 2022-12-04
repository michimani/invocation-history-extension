FROM public.ecr.aws/lambda/provided:al2 as build
RUN yum install -y golang unzip
RUN go env -w GOPROXY=direct
ADD go.mod go.sum ./
RUN go mod download
ADD . .
RUN go build -o /main
RUN mkdir -p /opt
ADD ./bin/extension.zip ./
RUN unzip extension.zip -d /opt
RUN rm extension.zip

FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
COPY entry.sh /
RUN chmod 755 /entry.sh
RUN mkdir -p /opt/extensions
WORKDIR /opt/extensions
COPY --from=build /opt/extensions .
ENTRYPOINT [ "/entry.sh" ]
CMD ["/main"]