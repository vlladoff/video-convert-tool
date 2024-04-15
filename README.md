# video-convert-tool

## Prerequisites

Before getting started, please make sure you have the following prerequisites installed on your system:

- **Docker:** Docker is a platform for developing, shipping, and running applications. You can install Docker by following the official installation instructions for your operating system:
   - [Install Docker for Windows](https://docs.docker.com/desktop/install/windows-install/)
   - [Install Docker for macOS](https://docs.docker.com/desktop/install/mac-install/)
   - [Install Docker for Linux](https://docs.docker.com/desktop/install/linux-install/)

## Installation
1. Clone the Repository
   
   ```bash
   git clone git@github.com:vlladoff/video-convert-tool.git
   cd video-convert-tool
   ```
2. You'll find an example environment file named `.env` in the project root. You need to change this file
3. To build and run the video-convert-tool using Docker run the following command:
   ```shell
   docker-compose -f compose-local.yaml up --build -d

## OpenAPI doc
<a href="https://html-preview.github.io/?url=https://github.com/vlladoff/video-convert-tool/blob/master/openapi.html">Example</a>

Generate documentation from `openapi.yaml` using `redocly/cli`:
```shell
docker run --rm -v $PWD:/spec redocly/cli build-docs openapi.yaml