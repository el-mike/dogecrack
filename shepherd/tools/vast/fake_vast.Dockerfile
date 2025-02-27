# syntax=docker/dockerfile:1

FROM michalhuras/pitbull:dev_local

WORKDIR /app

RUN apt update && apt -y install openssh-server && \
  # Allow login as root via SSH
  echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
  # Set root password for SSH to 12345
  echo "root:12345" | chpasswd

# Start SSH service, remove pitbull's view file (if exists) to get clean state and wait on bash process.
ENTRYPOINT service ssh start && rm -f /app/progress_view.txt && /bin/bash
