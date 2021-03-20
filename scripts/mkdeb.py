#!/usr/bin/python3 -u

import os
import pathlib
import re
import subprocess
import sys
import yaml
from typing import List

# The outer capture group will grab the whole variable including the dollar sign
# and brackets, the inner capture group will only grab the variable name itself
# e.g. '${VER}-${REL}' -> [('${VER}', 'VER'), ('${REL}', 'REL')]
ENV_RE = re.compile(r'(\$\{([A-Za-z0-9\_]+)\})')

def print_error(*args, **kwargs):
  print(*args, file=sys.stderr, **kwargs)

def read_config(file: str) -> dict:
  with open(file, 'r') as f:
    return yaml.safe_load(f)

# Get the total size of a directory in bytes
def dir_size(path: str) -> int:
  root = pathlib.Path(path)
  return sum(f.stat().st_size for f in root.glob('**/*') if f.is_file())

# Parse strings containing env variables and make replacements if applicable
# e.g. '${VER}-${REL}' -> '0.8.4-1'
def parse_env(text: str, env: dict) -> str:
  result = text
  for expr, name in ENV_RE.findall(text):
    if (env_val := env.get(name)) is not None:
      result = result.replace(expr, env_val)
  return result

def write_control_file(prefix: str, data: dict, env: dict) -> None:
  lines = ''
  for key, val in data.items():
    lines += f"{key}: {parse_env(val, env)}\n"
  with open(os.path.join(prefix, 'DEBIAN/control'), 'w+') as f:
    f.write(lines)

def write_script(prefix: str, filename: str, content: str) -> None:
  script = os.path.join(prefix, f"DEBIAN/{filename}")
  with open(script, 'w+') as f:
    f.write(content)
  os.chmod(script, 0o775)

def write_conffiles(prefix: str) -> None:
  root = pathlib.Path(os.path.join(prefix, 'etc'))
  files = (
    '/' + os.path.relpath(str(f), prefix)
    for f in root.glob('**/*')
    if f.is_file()
  )
  # As per dpkg-deb requirement we need a newline at the end of the file
  content = '\n'.join(files) + '\n'
  with open(os.path.join(prefix, 'DEBIAN/conffiles'), 'w+') as f:
    f.write(content)

def write_md5sums(prefix: str) -> None:
  cmd = r"""
    find . -mindepth 1 -type f -not -path './DEBIAN/*' |\
    sed 's|^./||' | sort | xargs md5sum
  """
  # https://docs.python.org/3/library/subprocess.html#subprocess.check_output
  output = subprocess.check_output(
    cmd,
    shell=True,
    cwd=prefix,
    text=True
  ).strip()
  with open(os.path.join(prefix, 'DEBIAN/md5sums'), 'w+') as f:
    f.write(output)

def make_deb(prefix: str, args: List[str]) -> int:
  final_args = ['dpkg-deb', *args, '--build', prefix]
  return subprocess.run(final_args).returncode

def main() -> int:
  if not (config_file := os.environ.get('CONFIG_FILE')):
    print_error('the CONFIG_FILE variable is not set')
    return 1
  if not (prefix := os.environ.get('PREFIX')):
    print_error('the PREFIX variable is not set')
    return 1

  pathlib.Path(os.path.join(prefix, 'DEBIAN')).mkdir(
    parents=True,
    exist_ok=True
  )

  config = read_config(config_file)

  install_cmd = (
    cmd
    if (cmd := config.get('build', {}).get('install', {}).get('cmd'))
    else ['make', 'install']
  )
  install_env = config.get('build', {}).get('install', {}).get('env', {})
  for key, val in install_env.items():
    install_env[key] = parse_env(val, os.environ)
  if (ret := subprocess.run(install_cmd, env=install_env).returncode) != 0:
    return ret

  env = {
    **os.environ,
    'SIZE_KB': os.environ.get('SIZE_KB') or str(int(dir_size(prefix) / 1024))
  }
  write_control_file(prefix, config.get('control', {}), env)

  scripts = config.get('scripts', {})
  for filename, content in scripts.items():
    write_script(prefix, filename, content)

  write_conffiles(prefix)

  write_md5sums(prefix)

  args = config.get('build', {}).get('args', [])
  if (ret := make_deb(prefix, args)) != 0:
    return ret

  if (target_name := config.get('build', {}).get('rename')):
    dir_name = os.path.dirname(prefix)
    final_name = parse_env(target_name, env)
    os.rename(f"{prefix}.deb", os.path.join(dir_name, final_name))

  return 0

if __name__ == '__main__':
  sys.exit(main())
