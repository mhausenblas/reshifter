FROM centos:7
ARG release_version
LABEL version=$release_version \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

COPY ui/* /app/ui/
WORKDIR /app
RUN curl -s -L 'https://github.com/mhausenblas/reshifter/releases/download/v'$release_version'-alpha/reshifter' -o reshifter && \
    chmod +x reshifter
EXPOSE 8080
CMD ["/app/reshifter"]
