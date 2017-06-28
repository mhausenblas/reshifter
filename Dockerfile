FROM centos:7
LABEL version="0.2" \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

COPY ui/* /app/ui/
WORKDIR /app
RUN curl -s -L https://github.com/mhausenblas/reshifter/releases/download/v0.2.3-alpha/reshifter -o reshifter && \
    chmod +x reshifter
EXPOSE 8080
CMD ["/app/reshifter"]
