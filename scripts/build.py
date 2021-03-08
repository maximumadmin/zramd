#!/usr/bin/python3 -u

import multiprocessing
import os
import subprocess
import sys
from typing import Optional

TARGETS = (
  ('arm', '6'),
  ('arm', '7'),
  ('arm64', None),
  ('amd64', None),
)

def build(goarch: Optional[str], goarm: Optional[str]) -> int:
  name = f"output=dist/zramd_{goarch}{goarm or ''}"
  proc = subprocess.run(
    ['make', 'release', name, 'compress=1', 'skip_clean=1'],
    env={
      **os.environ,
      'GOOS': 'linux',
      'GOARCH': goarch,
      **({'GOARM': goarm} if goarch == 'arm' else {})
    }
  )
  return proc.returncode

def clean() -> int:
  return subprocess.run(['make', 'clean'], env=os.environ).returncode

def main() -> int:
  if (ret := clean()) != 0:
    return ret
  processes = int(os.environ.get('PROCESSES', '1'))
  with multiprocessing.Pool(processes) as pool:
    codes = pool.starmap(build, TARGETS)
  return 0 if all(map(lambda x: x == 0, codes)) else 1

if __name__ == '__main__':
  sys.exit(main())
