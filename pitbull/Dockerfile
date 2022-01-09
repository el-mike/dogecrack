FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

COPY install.sh pitbull.sh capture_output.sh download_passlist.sh run_btcrecover.sh status.sh output.sh ./

RUN apt update && cp ./pitbull.sh pitbull && \
  chmod +x ./install.sh ./capture_output.sh ./download_passlist.sh ./run_btcrecover.sh ./status.sh ./output.sh ./pitbull && \
  # Adds pitbull executable to PATH, so client can run it from any directory.
  # Since bash init depends on using login/non-login shell, we put it two files,
  # to handle both cases.
  echo "export PATH=/app:${PATH}" >> ~/.profile && \
  echo "export PATH=/app:${PATH}" >> ~/.bashrc && \
  ./install.sh

 
ENTRYPOINT ["/bin/bash"]
