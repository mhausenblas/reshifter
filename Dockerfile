FROM scratch
ARG rversion
FROM centos:7
LABEL version=$rversion \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

COPY ui/* /app/ui/
WORKDIR /app
RUN curl -s -L 'https://github.com/mhausenblas/reshifter/releases/download/v'$rversion'-alpha/reshifter' -o reshifter && \
    chmod +x reshifter
EXPOSE 8080
USER 1001
CMD ["/app/reshifter"]
