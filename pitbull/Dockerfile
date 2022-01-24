FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

COPY install.sh helpers.sh pitbull.sh run_pitbull.sh capture_output.sh download_passlist.sh \
  run_btcrecover.sh status.sh output.sh progress.sh kill.sh ./

RUN apt update && cp ./pitbull.sh pitbull && \
  chmod +x ./install.sh helpers.sh run_pitbull.sh ./capture_output.sh ./download_passlist.sh \
    ./run_btcrecover.sh ./status.sh ./output.sh ./progress.sh ./kill.sh ./pitbull && \
  # Adds pitbull executable to PATH, so client can run it from any directory.
  # Since pitbull can be run via all sorts of shells (interactive, non-interactive,
  # login, non-login), we add it to bash config files, as well as to /etc/environment.
  echo "export PATH=/app:${PATH}" >> ~/.profile && \
  echo "export PATH=/app:${PATH}" >> ~/.bashrc && \
  echo "export PATH=/app:${PATH}" >> /etc/environment && \
  ./install.sh

 
ENTRYPOINT ["/bin/bash"]
