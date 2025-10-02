# Dockerfile (Go + user 1000 + sudo + devcontainer-friendly)
FROM golang:1.25.1-bookworm

# --- args personalizáveis ---
ARG USERNAME=dev
ARG UID=1000
ARG GID=1000
ARG WORKDIR=/home/${USERNAME}/.local/share/dev-gadgets

# --- pacotes mínimos + sudo ---
RUN apt-get update && apt-get install -y --no-install-recommends \
      sudo ca-certificates git bash-completion less \
    && rm -rf /var/lib/apt/lists/*

# --- cria usuário/grupo 1000, dá sudo e ajusta donos ---
RUN groupadd --gid ${GID} ${USERNAME} \
 && useradd --uid ${UID} --gid ${GID} -m ${USERNAME} \
 && echo "${USERNAME} ALL=(root) NOPASSWD:ALL" > /etc/sudoers.d/${USERNAME} \
 && chmod 0440 /etc/sudoers.d/${USERNAME} \
 # diretórios úteis p/ devcontainer
 && mkdir -p ${WORKDIR} \
 && chown -R ${USERNAME}:${USERNAME} ${WORKDIR} /go

# workspace padrão do devcontainer
WORKDIR ${WORKDIR}

# zera surpresas: shell de login e PATH com ~/.local/bin
SHELL ["/bin/bash", "-l", "-c"]
ENV PATH="/home/${USERNAME}/.local/bin:${PATH}"

# roda como usuário normal
USER ${USERNAME}
