FROM michalhuras/pitbull:3.0
# FROM michalhuras/pitbull:dev_local

WORKDIR /app

RUN apt update && apt -y install openssh-server && \
  # Allow login as root via SSH
  echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
  # Set root password for SSH to 12345
  echo "root:12345" | chpasswd

# Start SSH service and wait on bash process.
ENTRYPOINT service ssh start && /bin/bash
