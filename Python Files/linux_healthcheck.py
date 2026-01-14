#!/usr/bin/env python3
import subprocess
import shutil
from datetime import datetime
from pathlib import Path

OUTPUT_LINES = []

def log(line=""):
    print(line)
    OUTPUT_LINES.append(line)

# This is the command execution engine:
def run(cmd):
    log(f"$ {cmd}")
    try:
        result = subprocess.run(
            cmd,
            shell=True,
            text=True,
            capture_output=True
        )
        if result.stdout:
            log(result.stdout.rstrip())
        if result.stderr:
            log(result.stderr.rstrip())
    except Exception as e:
        log(f"ERROR: {e}")
    log()

# The format section so the healthcheck can have separators and be readable:
def header(title):
    log("=" * 70)
    log(title)
    log("=" * 70)

# Commands to be executed:
def main():
    timestamp = datetime.now().strftime("%Y-%m-%d_%H-%M-%S")

    # Save path: ~/Documents/Linux System Healthchecks
    save_dir = Path.home() / "Documents" / "Linux System Healthchecks"
    report_file = save_dir / f"linux_healthcheck_{timestamp}.txt"

    log(f"Linux System Healthcheck — {datetime.now()}\n")  # Establishes the sequence of checks.

    header("System Uptime and Load Average")
    run("uptime")  # check for how long the system has been up

    header("Top CPU-Consuming Processes")
    run("ps aux --sort=-%cpu | head -n 10")  # list all processes sorted by CPU usage (top 10)

    header("Disk Space Usage")
    run("df -h")  # check disk space usage in human-readable format

    header("Disk Usage by Root Directories")
    run("du -sh /* 2>/dev/null | sort -hr")  # check disk usage of root directories, suppress errors, sort by size

    header("Memory and Swap Usage")
    run("free -h")  # check memory and swap usage in human-readable format

    header("Failed System Services")
    if shutil.which("systemctl"):
        run("systemctl --failed")  # check for failed system services
    else:
        log("systemctl not available on this system. This may be a minimal or non-systemd environment.\n")

    header("Listening Network Ports")
    run("ss -tulnp")  # list all listening network ports with associated processes

    # Create folders if they do not exist and save report
    save_dir.mkdir(parents=True, exist_ok=True)
    report_file.write_text("\n".join(OUTPUT_LINES) + "\n")

    log(f"Report saved to: {report_file}")

if __name__ == "__main__":
    main()  # If run directly: python3 linux_healthcheck.py


