FROM centos
ARG rversion
LABEL version=$rversion \
      description="Kubernetes admin tool for backup and restoring clusters" \
      maintainer="michael.hausenblas@gmail.com"

COPY app/ui/* /app/ui/
WORKDIR /app
## for online development, uncomment the following 2 lines:
RUN curl -sf -L 'https://github.com/mhausenblas/reshifter/releases/download/v'$rversion'-alpha/reshifter' -o reshifter && \
    chmod +x reshifter
## for pure offline development, uncoment the following 2 lines:
# COPY reshifter .
# RUN chmod +x reshifter
RUN chown -R 1001:1 /app
USER 1001
EXPOSE 8080
CMD ["/app/reshifter"]
