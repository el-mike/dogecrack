FROM nvidia/opencl:runtime-ubuntu18.04

WORKDIR /app

COPY . .

RUN chmod +x ./install.sh && cp ./pitbull.sh pitbull && chmod +x ./pitbull && \
  echo "export PATH=/app:${PATH}" >> /root/.bashrc

 
ENTRYPOINT ["/bin/bash"]
