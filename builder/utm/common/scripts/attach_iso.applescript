---
-- attach_iso.applescript
-- This script attaches an ISO file to a specified UTM virtual machine at index 1 (first drive).
-- Usage: osascript attach_iso.applescript <VM_ID> --iso <ISO_PATH>
-- Example: osascript attach_iso.applescript test --iso "ubuntu-24.04-live-server-arm64.iso"
on run argv
  set vmId to item 1 of argv # ID of the VM
  -- Parse the --iso argument
  set isoPath to item 3 of argv as string

  -- Attached drives to the VM
  set isoFile to POSIX file (POSIX path of isoPath)

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing drives
    set vmDrives to drives of config
    --- create a new drive
    set newDrive to {removable:true, source:isoFile}
    -- Add the new drive to the beginning of the list
    set vmDrives to {newDrive} & vmDrives
    --- set drives with new drive list
    set drives of config to vmDrives

    --- save the configuration (VM must be stopped)
    update configuration of vm with config
  end tell
end run