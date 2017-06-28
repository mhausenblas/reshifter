FROM centos:7
LABEL version="0.2" \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

RUN curl -s -L https://github.com/mhausenblas/reshifter/releases/download/v0.2-alpha/reshifter -o reshifter
COPY reshifter /app/reshifter
COPY ui/* /app/ui/

EXPOSE 8080
CMD ["/app/reshifter"]
