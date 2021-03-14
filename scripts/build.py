#!/usr/bin/python3 -u

import os
import subprocess
import sys
from typing import Optional

TARGETS = (
  ('arm', '6', 'armel'),
  ('arm', '7', 'armhf'),
  ('arm64', None, 'arm64'),
  ('amd64', None, 'amd64'),
)

def build(goarch: str, goarm: Optional[str], friendly_arch: str) -> int:
  out_file = f"dist/zramd_{friendly_arch}"
  prefix = f"dist/zramd_root_{friendly_arch}"
  version, release, *_ = [*os.environ['CURRENT_TAG'].split('-'), '']
  proc = subprocess.run(
    ['make', f"output={out_file}", 'make_tgz=1', 'make_deb=1', 'skip_clean=1'],
    env={
      # Pass all environment variables, contains some Go variables
      **os.environ,
      # Set Go build-specific variables
      'GOOS': 'linux',
      'GOARCH': goarch,
      **({'GOARM': goarm} if goarch == 'arm' else {}),
      # Required to create a Debian package
      'DEB_ARCH': friendly_arch,
      'PREFIX': prefix,
      'VERSION': version,
      'RELEASE': release or '1',
      'BIN_FILE': out_file
    }
  )
  return proc.returncode

def clean() -> int:
  return subprocess.run(['make', 'clean'], env=os.environ).returncode

def main() -> int:
  if (ret := clean()) != 0:
    return ret
  for target in TARGETS:
    if (ret := build(*target)) != 0:
      return ret
  return 0

if __name__ == '__main__':
  sys.exit(main())
