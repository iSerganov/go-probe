FROM golang:1.18-bullseye AS golang-ffprobe

# Install ffmpeg containing ffprobe
RUN apt-get update && apt-get install -y --no-install-recommends \
		ffmpeg \
	&& rm -rf /var/lib/apt/lists/*

