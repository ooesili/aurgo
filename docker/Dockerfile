FROM ooesili/arch-base

RUN ["useradd", "-ms", "/bin/bash", "-G", "wheel", "aurgo"]
RUN ["bash", "-c", "chpasswd <<< aurgo:password1"]
COPY "sudoers-wheel" "/etc/sudoers.d/wheel"

USER "aurgo"
WORKDIR "/home/aurgo"
