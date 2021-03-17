#!/usr/bin/python3 -u

import itertools
import json
import os
import subprocess
import sys
import time
import urllib.parse
from http.client import HTTPResponse
from urllib.error import HTTPError, URLError
from urllib.request import Request, urlopen
from typing import Iterable, List, Optional, Tuple

Response = Tuple[Optional[Exception], Optional[dict]]
Content = Optional[object]

MAX_RETRIES = 2
DEFAULT_HEADERS = {
  'Content-Type': 'application/json',
  'Accept': 'application/vnd.github.v3+json'
}

def print_error(*args, **kwargs):
  print(*args, file=sys.stderr, **kwargs)

def flatten(iterable: Iterable) -> List:
  return list(itertools.chain(*iterable))

def request_json(url: str, headers: Optional[dict], body=None) -> Response:
  req = Request(url, method='POST')
  for key, value in {**DEFAULT_HEADERS, **(headers or {})}.items():
    req.add_header(key, value)
  body_bytes = json.dumps(body).encode('utf-8') if body else None
  response_body: bytes
  response_code: int
  try:
    res: HTTPResponse
    with urlopen(req, data=body_bytes) as res:
      response_body = res.read()
      response_code = res.getcode()
  except (HTTPError, URLError) as e:
    return (e, None)
  response_json: dict = json.loads(response_body)
  if response_code >= 400:
    message = {'status': response_code, 'body': response_body}
    return (Exception(json.dumps(message, indent=2)), None)
  return (None, response_json)

def request_safe(
  url: str,
  headers: Optional[dict],
  body=None,
  max_retries=MAX_RETRIES
) -> Content:
  i = 0
  while True:
    url_path = urllib.parse.urlparse(url).path
    print(f"[{i}] Connecting to \"{url_path}\"...")
    (e, data) = request_json(url, headers, body)
    if not e:
      return data
    print_error(e)
    if max_retries > 0 and (i == max_retries - 1):
      break
    time.sleep(2)
    i += 1
  return None

def create_release(owner: str, repo: str, token: str, body: dict) -> Content:
  url = f"https://api.github.com/repos/{owner}/{repo}/releases"
  headers = {'Authorization': f"token {token}"}
  return request_safe(url, headers, body)

# Easier than do it in pure Python and does not require additional libraries
def curl_upload(upload_url: str, headers: dict, file_path: str):
  header_list = flatten(['-H', f"{h}: {headers[h]}"] for h in headers)
  result = subprocess.run([
    'curl',
    '-sS',
    *header_list,
    '--data-binary',
    f"@{file_path}",
    upload_url
  ])
  return result.returncode

def upload_asset(
  upload_url: str,
  token: str,
  file_path: str,
  max_retries=MAX_RETRIES
) -> bool:
  headers = {'Authorization': f"token {token}"}
  for i in range(max_retries):
    if (ret := curl_upload(upload_url, headers, file_path)) == 0:
      break
    print_error(f"curl finished with exit code {ret}")
    if i == max_retries - 1:
      return False
  return True

def main() -> int:
  # Get all required variables from env, if we get an error that's good, it just
  # means that we forgot to pass a variable
  owner = os.environ['REPO_OWNER']
  repo = os.environ['REPO_NAME']
  token = os.environ['GH_RELEASE_TOKEN']
  release_tag = os.environ['CURRENT_TAG']

  # Get list of friendly architectures from the first argument, each item will
  # have a corresponding binary, tar and deb file under dist/
  friendly_arches = sys.argv[1].split(',')

  # Create a GitHub release (requires a valid tag to exist)
  release_data = create_release(owner, repo, token, {
    'tag_name': release_tag,
    'name': f"zramd {release_tag}",
    'body': ''
  })
  if not release_data:
    return 1

  # https://uploads.github.com/repos/OWNER/REPO/releases/ID/assets{?name,label}
  upload_url: str = release_data['upload_url']
  upload_url = upload_url.split('/assets')[0] + '/assets?name='

  # Upload assets
  assets = flatten(
    (f"zramd_{a}.deb", f"zramd_{a}.tar.gz")
    for a in friendly_arches
  )
  for asset in assets:
    current_file = os.path.join('dist', asset)
    current_url = upload_url + os.path.basename(current_file)
    if not upload_asset(current_url, token, current_file):
      return 1

  return 0

if __name__ == '__main__':
  sys.exit(main())
