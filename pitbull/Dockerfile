FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

COPY . .

RUN chmod +x ./install.sh && cp ./pitbull.sh pitbull && chmod +x ./pitbull && \
  # Adds pitbull executable to PATH, so client can run it from any directory.
  echo "export PATH=/app:${PATH}" >> /root/.bashrc && \
  ./install.sh

 
ENTRYPOINT ["/bin/bash"]
