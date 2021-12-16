FROM michalhuras/pitbull:1.0

WORKDIR /app

RUN apt update && apt -y install openssh-server tmux && \
  # Allow login as root via SSH
  echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
  # Set root password for SSH to 12345
  echo 'root:12345' | chpasswd && \
 chsh -s /usr/bin/tmux root

# Start SSH service and wait on bash process.
ENTRYPOINT service ssh start && /bin/bash
