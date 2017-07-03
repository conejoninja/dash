FROM debian

ADD ./dash /data/dash
ADD ./static /data/static
ADD ./views /data/views
RUN chmod +x /data/dash
EXPOSE 1313

ENTRYPOINT cd /data/ && /data/dash