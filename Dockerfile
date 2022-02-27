FROM golang:1.17
ENV TZ Asia/Tokyo
ENV GOBIN=/go/bin
WORKDIR /go/src/app
RUN go install github.com/x-motemen/gore/cmd/gore@latest && \
    go install github.com/go-delve/delve/cmd/dlv@master && \
    cp ${GOBIN}/dlv ${GOBIN}/dlv-dap && \
    go install golang.org/x/tools/gopls@latest && \
    go install golang.org/x/tools/cmd/goimports@latest && \
    go install github.com/ramya-rao-a/go-outline@latest && \
    go install github.com/stamblerre/gocode@latest && \
    go install github.com/uudashr/gopkgs/v2/cmd/gopkgs@latest && \
    go install github.com/rogpeppe/godef@latest && \
    go install honnef.co/go/tools/cmd/staticcheck@latest && \
    # 単体テストをいい感じに生成してくれるツール
    go install github.com/cweill/gotests/gotests@latest
EXPOSE 8080
CMD ["go", "run", "main.go"]
