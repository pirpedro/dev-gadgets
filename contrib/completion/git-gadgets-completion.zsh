_dev-gadgets() {
  local curcontext=$curcontext state line ret=1
  declare -A opt_args

  _arguments -C \
    ': :->command' \
    '*:: :->option-or-argument' && ret=0

  case $state in
  command)
    declare -a commands
    commands=(
      'init: configure current repository'
      'update:update git-extras'
    )
    _describe -t commands command commands && ret=0
    ;;
  esac

  _arguments \
    '(-v --version)'{-v,--version}'[show current version]'
}

zstyle -g existing_user_commands ':completion:*:*:git:*' user-commands

zstyle ':completion:*:*:git:*' user-commands $existing_user_commands \
  gadgets:'bootstrapping git projects'
