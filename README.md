# zramd

Automatically setup swap on zram ✨

## Why swap on zram?

* Significantly improves system responsiveness, especially when swap is under pressure.
* More secure, user data leaks into swap are on volatile media.
* Without swap-on-drive, there's better utilization of a limited resource: benefit of swap without the drive space consumption.
* Further reduces the time to out-of-memory kill, when workloads exceed limits.

See also https://fedoraproject.org/wiki/Changes/SwapOnZRAM#Benefit_to_Fedora

## Installation

### Install on Arch Linux from the AUR

* Install the `zramd` package form the [AUR](https://aur.archlinux.org/packages/zramd/).
* Enable and start the service:
  ```bash
  sudo systemctl enable --now zramd
  ```

### Install on Ubuntu / Debian / Raspberry Pi OS

* Head to the releases section and download the `.deb` file corresponding to your system architecture
* Install using `sudo dpkg -i DEB_FILE`

### Manual installation on any distribution without systemd

* Head to the releases section and download the `.tar.gz` file corresponding to your system architecture
* Extract the downloaded file e.g. `tar xf TAR_FILE`
* Copy the `zramd` binary to `/usr/local/bin`.
* There are various ways to setup autostart depending on your init system, for example you can add a line to `/etc/rc.local` e.g.
  ```bash
  /usr/local/bin/zramd start
  ```

## Usage

* zramd --help
  ```
  Usage: zramd <command> [<args>]

  Options:
    --help, -h             display this help and exit
    --version              display version and exit

  Commands:
    start                  load zram module and setup swap devices
    stop                   stop swap devices and unload zram module
  ```

* zramd start --help
  ```
  Usage: zramd start [--algorithm ALGORITHM] [--fraction FRACTION] [--max-size MAX_SIZE] [--num-devices NUM_DEVICES] [--priority PRIORITY] [--skip-vm]

  Options:
    --algorithm ALGORITHM, -a ALGORITHM
                           zram compression algorithm
                           [default: zstd, env: ALGORITHM]
    --fraction FRACTION, -f FRACTION
                           maximum percentage of RAM allowed to use
                           [default: 1.0, env: FRACTION]
    --max-size MAX_SIZE, -m MAX_SIZE
                           maximum total MB of swap to allocate
                           [default: 8192, env: MAX_SIZE]
    --num-devices NUM_DEVICES, -n NUM_DEVICES
                           maximum number of zram devices to create
                           [default: 1, env: NUM_DEVICES]
    --priority PRIORITY, -p PRIORITY
                           swap priority
                           [default: 100, env: PRIORITY]
    --skip-vm, -s          skip initialization if running on a VM
                           [default: false, env: SKIP_VM]
    --help, -h             display this help and exit
    --version              display version and exit
  ```

## Compilation

### With Docker

* Choose a valid git tag and run the `make docker` command e.g.
  ```bash
  # This command will create builds (.tar.gz and .deb) for all supported architectures
  CURRENT_TAG=v0.8.5 make docker
  ```

### Manual Compilation

* Install `go` (at least version 1.16), the command may be different depending on the distribution:
  ```bash
  # ArchLinux
  sudo pacman -S go

  # Ubuntu
  sudo apt-get install golang
  ```
* You can run `make` to create a build with the same architecture as the current system:
  ```bash
  # If you cloned this repository you can just run make
  make

  # To create a Raspberry Pi build you need to specify the arch e.g.
  GOOS=linux GOARCH=arm GOARM=7 make

  # If you downloaded a .tar.gz or .zip (instead of cloning this repo) you need to specify additional info
  CURRENT_DATE=$(date --iso-8601=seconds) VERSION=Unknown make

  # So, to target the Raspberry Pi (without a repo) the command would look like
  CURRENT_DATE=$(date --iso-8601=seconds) VERSION=Unknown GOOS=linux GOARCH=arm GOARM=7 make
  ```
* A new executable called `zramd.bin` will be created under the `dist/` directory, now you can uninstall `go` if you like.

## Configuration

### With systemd

* The default configuration file is located at `/etc/default/zramd`, just edit the variables as you like and restart the `zramd` service i.e. `sudo systemctl restart zramd`

### Without systemd

* Just change the arguments as you like, e.g. `zramd start --max-size 1024` or `zramd start --fraction 0.5 --priority 0`

## Troubleshooting

* **modprobe: FATAL: Module zram not found in directory /lib/modules/...**  
  It can happen if you try to start the `zramd` service after a kernel upgrade, you just need to restart your computer.
* **error: swapon: /dev/zramX: swapon failed: Operation not permitted**  
  First make sure that you are running as root (or at least that you have the required capabilities), also keep in mind that Linux only supports up to 32 swap devices (although it can start throwing the error from above when using a high value like 24).

## Notes

* **Avoid** using other zram-related packages along this one, `zramd` loads and unloads the zram kernel module assuming that the system is not using zram for other stuff (like mounting `/tmp` over zram).
* Do **not** use zswap with zram, it would unnecessarily cause data to be [compressed and decompressed back and forth](https://www.phoronix.com/forums/forum/software/distributions/1231542-fedora-34-looking-to-tweak-default-zram-configuration/page5#post1232327).
* When dealing with virtual machines, zram should be used on the **host** OS so guest memory can be compressed transparently, see also comments on original zram [implementation](https://code.google.com/archive/p/compcache/).
  * If you boot the same system on a real computer as well as on a virtual machine, you can use the `--skip-vm` parameter to avoid initialization when running inside a virtual machine.
* For best results install `systemd-oomd` or `earlyoom` (they may not be available on all distributions).
* You can use `swapon -show` or `zramctl` to see all swap devices currently in use, this is useful if you want to confirm that all of the zram devices were setup correctly.
* To quickly fill the memory, you can use `tail /dev/zero` but keep in mind that your system may become unresponsive if you do not have an application like `earlyoom` to kill `tail` just before it reaches the memory limit.
* To test some zramd commands under the same conditions as the systemd unit you can use `systemd-run` e.g.
  ```bash
  sudo systemd-run -t \
    -p ProtectHostname=yes \
    -p PrivateNetwork=yes \
    -p IPAddressDeny=any \
    -p NoNewPrivileges=yes \
    -p RestrictNamespaces=yes \
    -p RestrictRealtime=yes \
    -p RestrictSUIDSGID=yes \
    -p MemoryDenyWriteExecute=yes \
    -p LockPersonality=yes \
    -p 'CapabilityBoundingSet=CAP_SYS_ADMIN CAP_SYS_MODULE' \
    -p 'SystemCallFilter=@module @swap @system-service' \
    -p SystemCallArchitectures=native \
    -p SystemCallErrorNumber=EPERM \
    -p 'DeviceAllow=block-* rw' \
    -p DevicePolicy=closed \
    -p RestrictAddressFamilies=AF_UNIX \
    -p RestrictAddressFamilies=~AF_UNIX \
    /usr/bin/zramd --version
  ```
