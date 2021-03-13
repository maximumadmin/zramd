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

def read_config(file: str) -> str:
  with open(file, 'r') as f:
    return yaml.safe_load(f)

def write_control_file(prefix: str, data: dict) -> None:
  lines = ''
  for key, val in data.items():
    next_val = val
    for expr, name in ENV_RE.findall(val):
      if (env_val := os.environ.get(name)):
        next_val = next_val.replace(expr, env_val)
    lines += f"{key}: {next_val}\n"
  with open(os.path.join(prefix, 'DEBIAN/control'), 'w+') as f:
    f.write(lines)

def write_script(prefix: str, filename: str, content: str) -> None:
  script = os.path.join(prefix, f"DEBIAN/{filename}")
  with open(script, 'w+') as f:
    f.write(content)
  os.chmod(script, 0o775)

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

  config: dict = read_config(config_file)

  control = config.get('control', {})
  write_control_file(prefix, control)

  scripts = config.get('scripts', {})
  for filename, content in scripts.items():
    write_script(prefix, filename, content)

  write_md5sums(prefix)

  command = ['make', 'install']
  command = config.get('build', {'install': command}).get('install', command)
  if (ret := subprocess.run(command, env={'PREFIX': prefix}).returncode) != 0:
    return ret

  args = config.get('build', {'dpkg_deb': []}).get('dpkg_deb', [])
  if (ret := make_deb(prefix, args)) != 0:
    return ret

  return 0

if __name__ == '__main__':
  sys.exit(main())
