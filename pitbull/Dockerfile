FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

  # Setup and common scripts
COPY install.sh helpers.sh variables.sh \
  # Command scripts
  run_pitbull.sh status.sh progress.sh output.sh kill.sh errors.sh attach.sh \
  # Operation scripts
  capture_output.sh download_passlist.sh run_btcrecover.sh \
  # Entry point script
  pitbull.sh ./

RUN apt update && cp ./pitbull.sh pitbull && \
  chmod -R +x ./ && \
  # Adds pitbull executable to PATH, so client can run it from any directory.
  # Since pitbull can be run via all sorts of shells (interactive, non-interactive,
  # login, non-login), we add it to bash config files, as well as to /etc/environment.
  echo "export PATH=/app:${PATH}" >> ~/.profile && \
  echo "export PATH=/app:${PATH}" >> ~/.bashrc && \
  echo "export PATH=/app:${PATH}" >> /etc/environment && \
  ./install.sh

 
ENTRYPOINT ["/bin/bash"]
