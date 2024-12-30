# Fedora with QEMU and cloud-init

This covers the necessities for booting a [Fedora Cloud](https://fedoraproject.org/cloud) instance
with [QEMU](https://www.qemu.org) (directly) and having it execute
[cloud-init](https://cloudinit.readthedocs.io) after booting. This was done in service of a project
to emulate Raspberry Pi instances that could automatically configure themselves to join a network or
cluster.

## Things to Know

- **cloud-init** is an operating system initialization framework that can do system bootstrapping
  after an operating system starts (e.g. configuring users, configuring devices, mounting disks, and
  etc.).

- **QEMU** is a tool suite that allows for user-space CPU emulation, and can act as a simple
  hypervisor.

- **qcow**, specifically qcow2, is a copy-on-write disk image format used by QEMU. It has the
  ability to dynamically grow as data is added, or preserve the base image and store written data
  elsewhere.

- **System Management BIOS** (SMBIOS) is used to supply data from the BIOS to an operating system.
  In this case, it can be used to supply instructions to cloud-init.

## Steps (x86-64)

Assuming that we're working off of an x86-64 machine, these steps perform a basic POC of QEMU driven
Fedora Cloud before attempting to actually emulate another CPU architecture.

1.  Download the Fedora Cloud's QEMU qcow2 system for the x86-64 family of systems.

3.  Generate a random password for a user account being added in `user-data`.

    ```sh
    mkpasswd --method=SHA-512 <password>
    ```

2.  Create cloud-init files meta-data, user-data, and vendor-data.

    -   meta-data should contain:

        ```yaml
        instance-id: <something>
        ```

        where `something` is anything you like.

    -   user-data should contain:

        ```yaml
        #cloud-config
        users:
            - name: <username>
              groups: user,wheel
              hashed_passwd: <password>
              lock_passwd: false
              sudo: ALL=(ALL) NOPASSWD:ALL
        ```

        where `username` is anything you like and `password` is the output of the `mkpasswd`
        command.

    -   vendor-data can just be an empty file.

4.  Check the validity of the cloud-init user data.

    ```sh
    cloud-init schema --config-file user-data
    ```

5.  Create an ad-hoc IMDS webserver to serve the cloud-init files.

    ```sh
    python -m http.server --directory .
    ```

6.  Execute a QEMU start command against the x86-64 qcow2 image.

    ```sh
    qemu-system-x86_64 \
        -cpu host \
        -drive file=Fedora-Cloud-Base-39-1.5.x86_64.qcow2 \
        -m 512 \
        -machine accel=kvm \
        -nographic \
        -smbios type=1,serial=ds='nocloud;seedfrom=http://10.0.2.2:8000/'
    ```

7.  The console should show a system booting, some cloud-init messages, and a login prompt that
    accepts the username and password configured.

When you're ready to stop the system, you can just stop the QEMU process.

```sh
pkill qemu-system-x86_64
```

## QEMU Options

Here we try to break down the above QEMU command in order to more closely understand what we're
actually executing.

`-cpu` dictates the CPU model for QEMU to emulate. This can do different things depending on the
emulation type, and in this case, `host` dictates KVM to supply a CPU which has the same
capabilities as the host CPU.

`-drive` defines a drive for the emulated instance, which simultaneously defines a block device and
presents it as a device. In this case, we present the QCOW image itself as a drive.

`-m` defines the RAM size and can define various physical characteristics of memory, such as
expansion slots.`

`-machine` defines the emulated machine. This directive can also have QEMU use other hypervisors to
run virtual machines. `accel=kvm` enables KVM as an underlying hypervisor, known as "accelerators"
in QEMU parlance.

`-nographic` directs QEMU to display no guest graphics, and will instead direct output to the
console.

`-smbios` defines arguments for the SMBIOS.

## Steps (aarch64)

The above demonstrates the basics of engaging QEMU for an x86-64 machine on an x86-64 machine and
should successfully boot Fedora Cloud.

However, this isn't very representative of the Raspberry Pi architecture. We can also apply the same
cloud-init set-up for aarch64/arm64, which gets us closer (perhaps close enough) to what we want.

1.  Perform steps 1-7 above, downloading instead a qcow2 aarch64/arm64 image.

2.  Execute a QEMU start command against the aarch64 qcow2 image.

    ```sh
    qemu-system-aarch64 \
        -bios /usr/share/edk2/aarch64/QEMU_EFI.fd \
        -cpu cortex-a57 \
        -drive file=Fedora-Cloud-Base-39-1.5.aarch64.qcow2 \
        -m 1G \
        -machine virt \
        -nographic \
        -smbios type=1,serial=ds='nocloud;seedfrom=http://10.0.2.2:8000/' \
        -smp 2
    ```

    Note that the larger resources (virtual CPU and memory) provided are necessary for the image to
    successfully launch and present a login.

3.  The console should show a system booting, some cloud-init messages, and a login prompt that
    accepts the username and password configured.

## More QEMU Options

Arguments for running ARM machine in QEMU look a bit different, with one big difference being the
exclusion of KVM as a backing hypervisor.

`-bios` defines the BIOS file to load into memory and use for booting. In most cases, QEMU will load
a BIOS without being prompted, depending on the machine architecture, but in the case of ARM we
provide one that EDK2 happens to have.

`-machine virt` is described as a "general purpose platform" for running ARM based machines.

`-smp` (a reference to symmetric multi-processing) helps define specific CPU characteristics,
including cores, sockets, dies, threads, and other physical characteristics. `-smp 2` dictates the
machine should have two cores in one socket.

## Networking

Networking is easy to append onto the virtual instance using a user mode host network configuration.
We do this by adding an emulated network device and configuring it for internet access.

In this scheme, the guest will be able to access the host network and the public internet and the
guest will use the host's DNS resolution, but the host will have no route to the guest.

```sh
-netdev user,id=internet -device virtio-net-pci,netdev=internet
```

Configuring access from the host to the guest is a little more involved. It's possible to create a
bridge network on the host which the guest can take advantage of, but that can be difficult in
certain situations, like if the host only has wireless internet.

A simpler solution if you only need specific network paths is to allow QEMU to create a tunnel for a
particular protocol and pair of ports.

```sh
-netdev user,id=internet,hostfwd=tcp:127.0.0.1:8822-:22 -device virtio-net-pci,netdev=internet
```

This will allow the host to connect to the guest's SSH over port `8822`.
