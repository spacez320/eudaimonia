# Linux Memory

This is a set of notes on understanding Linux memory usage--how the kernel manages physical memory
resources and how to observe memory usage in processes.

## Overview

There are two types of memory:

- **Physical**, usually RAM.
- **Swap**, disk space that the kernel will use to extend available memory.

In Linux, physical memory complexities are abstracted and simplified for applications into **virtual
memory**.

Properties of virtual memory include:

- Virtual memory is controlled by a Memory Management Unit (MMU).
- Virtual memory addresses are mapped to physical ones when memory is accessed by the CPU.
- Virtual memory **pages** map to one or more physical memory **pages** (or page frames--logical
  divisions of memory), which are recorded in page tables.

Linux also divides memory into **zones** according to possible usage and requires for direct memory
access (DMA). This is architecture specific.

**Non-Uniform Memory Access (NUMA)** systems also affect memory allocation by using multiple memory
manage subsystems across memory locations with different latencies.

## Terms

Below is a quick look-up of various terms used to describe memory in linux.

| Term              | Description                                                             |
| ----------------- | ----------------------------------------------------------------------- |
| Available memory  | Memory that is available for allocation (including reclaimable memory). |
| Disk buffer       | Memory which holds file writes to eventually write to disk.             |
| Disk cache        | Memory which holds files expected to be read (avoiding disk IO).        |
| Free memory       | Memory currently being completely unused.                               |
| Phyiscal memory   | Actual memory, supported by hardware.                                   |
| Resident set size | Amount of physical memory a process is using.                           |
| Shared memory     | Memory which is being shared between processes.                         |
| Swap              | Disk which is used to extend total memory.                              |
| Used memory       | Memory that is being used by something.                                 |
| Virtual memory    | Abstracted memory combining physical memory and swap.                   |

## Kernel Memory Usage

The kernel itself manages memory for processes but also must preserve some for its own data
structures and uses. This section tries to describe relevant concepts.

Free memory is divided between **highmem** and **lowmem**.

- Highmem consists of everything above approximately 1GB of physical memory and is generally for
  user-space programs and page caching.
- Lowmem is available for kernel data structures.

The kernel keeps track of pages which are "reclaimable", mostly involving anonymous or page cache
memory (explained below). "Unreclaimable" memory may include DMA or kernel-used memory (exceptions
exist). Memory reclaim behavior depends on current available memory and is controlled by the
`kswapd` daemon.

Memory thresholds exist which affect reclaim behavior, including:

- **Low watermark**, which activates kswapd to find memory it can free.
- **Min watermark**, which will block memory allocation until enough pages exist.

In the event that there is so little memory left that the kernel itself is threatened, the **OOM
Killer** is invoked, which will forcibly sacrifice processes to free memory.

## Process Memory Usage

Processes are usually the main consumers of memory and will allocate memory by making kernel system
calls.

- Processes accure physical memory with mapped virtual memory.
- The amount of physical memory that processes are using is called **real memory**.

Memory that a process allocates that is not backed by the file system is called **anonymous
memory**. This also includes implicit memory created for a program's stack and heap. Writes to
anonymous memory are considered "dirty" and will be swapped if the kernel takes it back.

Processes will allocate or commit memory (e.g. with `malloc`). **Committed memory** is often
accounted for separately from **used memory** (processes may have committed more memory than they
are actually using).

Relevant system calls include:

- `malloc`, `free`, `calloc`, `realloc`

Relevant commands include:

- `free`

### Shared Memory

Shared memory is just that--memory which can be accessed by multiple processes. It can be considered
as a type of interprocess communication and  employs a special set of system calls. Shared memory is
mapped to a process's own virtual memory space.

There are two APIs which control shared memory use: System V and POSIX.

Relevant system calls include:

- `shmget`, `shmat`, `shift`, `shmctl` (System V)
- `shm_open`, `mmap`, munmap` (POSIX)

Relevant commands include:

- `ipcs`, show information on IPC facilities.

## File Memory Usage

Memory is often populated by reading from disk (i.e. reading files). In order to avoid repeated disk
interactions, Linux provides a **page cache** which caches file reads and writes.

Pages which page cache data for writing to a file system are called **"dirty" pages**. The kernel
knows to synchronize (or flush) this data to disk before that memory may be reallocated.

## Observability

This section describes the ways an operator can observe memory related information.

### Proc

The proc file system provides information from the kernel on performance and process state,
including memory.

There is man page documentation on the proc fs which can expand on concepts and file paths described
below.

```sh
man 5 proc
```

### `/proc/meminfo`

`/proc/meminfo` shows global memory statistics, including total available and current usage.

Tools that use `/proc/meminfo` include:

- `free`

### `/proc/<pid>/map_files`

_TODO_

### `/proc/<pid>/maps`

This file path shows memory mappings within a processes' virtual memory.

### `/proc/<pid>/mem`

This file path provides direct access to the memory of a process outside of system calls.

### `/proc/<pid>/oom_adj`, `/proc/<pid>/oom_score`, and `/proc/<pid>/oom_score_adj`

_TODO_

### `/proc/<pid>/oom_adj`

_TODO_

### `/proc/<pid>/smaps`

_TODO_

### `/proc/<pid>/stat`

Shows general status information about a process, including memory usage.
Relevant information for a process includes:

- `minflt`, `cminflt`, `majflt`, and `cmajflt` which keeps track of page faults.
- `vsize` represents virtual memory size.
- `rss` is "resident set size"--the number of pages in real memory (note that `man 5 proc` claims
  this value is inaccurate and `/proc/<pid>/statm` should be used instead).

Tools that use `/proc/<pid>/stat` include:

- `ps`

### Commands

There are several commands useful for observing memory in interactive ways.

- `free` can show system-wide, human-readable statistics on memory usage.
- `top` is an interactive, dynamic, real-time process statistics table.
- `vmstat` reports on system activity, including memory (also other things like IO, disk, and CPU).

## TODO

This is a holding zone for topics on Linux memory that are mentioned in documentation or other
literature that this article doesn't yet try to explain.

- What are kernel buffers?
- What are memory mapped files and what is the `mmap` system call used for?
- What are anonymous pages?
- What is tmpfs?
- What is Slab and how does it relate to kernel memory use?
- What are bounce buffers?
- What is vmalloc?
- What are huge pages?
- How does memory compaction work?
- What is CMA?
- What are page faults?
- What is real memory?

## Reference

- <https://www.kernel.org/doc/html/next/admin-guide/mm/concepts.html>
- <https://www.baeldung.com/linux/read-process-memory>
