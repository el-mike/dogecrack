FROM michalhuras/pitbull:1.0

WORKDIR /app

RUN apt update && apt -y install openssh-server tmux && \
  # Allow login as root via SSH
  echo "PermitRootLogin yes" >> /etc/ssh/sshd_config && \
  # Set root password for SSH to 12345
  echo 'root:12345' | chpasswd && \
  # Start SSH service.
  /etc/init.d/ssh start && \
  chsh -s /usr/bin/tmux root

ENTRYPOINT [ "/bin/bash", "-l", "-c" ]
