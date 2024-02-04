FROM cosmtrek/air

RUN apt-get update && apt-get install libvips libvips-dev ffmpeg --yes && apt-get clean autoclean && apt-get autoremove --yes && rm -rf /var/lib/{apt,dpkg,cache,log}/

ENV GOCACHE /app/.cache
