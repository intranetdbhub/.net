FROM ubuntu:latest

# Install OpenSSH server
RUN apt-get update && \
    apt-get install -y openssh-server && \
    mkdir /var/run/sshd

# Set root password (lab use only)
RUN echo 'root:root' | chpasswd

# Allow root login
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config

#Mounted Volumes Permanent Visibility in Directory:
# Create symlink from /Documents to ~/Documents
RUN ln -s /Documents /root/Documents
# Create symlink from /Shared to ~/Shared
RUN ln -s /Shared /root/Shared

EXPOSE 22

CMD ["/usr/sbin/sshd", "-D"]
