FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

COPY . .

RUN apt update && chmod +x ./install.sh && cp ./pitbull.sh pitbull && chmod +x ./pitbull && \
  # Adds pitbull executable to PATH, so client can run it from any directory.
  # Since bash init depends on using login/non-login shell, we put it two files,
  # to handle both cases.
  echo "export PATH=/app:${PATH}" >> ~/.profile && \
  echo "export PATH=/app:${PATH}" >> ~/.bashrc && \
  ./install.sh

 
ENTRYPOINT ["/bin/bash"]
