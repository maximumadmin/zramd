#!/usr/bin/python3 -u

import os
import subprocess
import sys
from typing import Optional, Tuple

TARGETS = (
  ('arm', '6', 'armel'),
  ('arm', '7', 'armhf'),
  ('arm64', None, 'arm64'),
  ('amd64', None, 'amd64'),
)

# Parse tag names like v0.8.5 or v0.8.5-1
def parse_tag(tag: str) -> Tuple[str, str]:
  version, release, *_ = [*tag.split('-'), '']
  # Remove the leading 'v' from version
  return (version[1:], release or '1')

def build(goarch: str, goarm: Optional[str], friendly_arch: str) -> int:
  out_file = f"build/zramd_{friendly_arch}"
  prefix = f"build/zramd_{friendly_arch}_root"
  version, release = parse_tag(os.environ['CURRENT_TAG'])
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
      'VERSION': version,
      'RELEASE': release,
      'PREFIX': prefix,
      'BIN_FILE': out_file
    }
  )
  return proc.returncode

def clean() -> int:
  return subprocess.run(['make', 'clean'], env=os.environ).returncode

def main() -> int:
  ret = clean()
  if ret != 0:
    return ret

  # Build all targets sequentially, building in parallel will have minimal or no
  # benefit and would make logging messy
  for target in TARGETS:
    ret = build(*target)
    if ret != 0:
      return ret

  # Finally write the used architectures so we can use them at later steps
  with open('targets.txt', 'w') as f:
    f.write(','.join(row[2] for row in TARGETS))

  return 0

if __name__ == '__main__':
  sys.exit(main())
