---
-- create_vm_from_source.applescript
-- This script creates a new VM with the specified properties and a source file as disk.
-- Usage: osascript create_vm_from_source.applescript --name <VM_NAME> --backend <BACKEND> --arch <ARCH> --source <SOURCE_DISK_PATH> 
-- Example: osascript create_vm_from_source.applescript --name "MyVM" --backend "QeMu" --arch "aarch64" --source "/path/to/source" 
on run argv
    -- Initialize variables
    set vmName to ""
    set vmBackend to ""
    set vmArch to ""

    -- Parse arguments
    repeat with i from 1 to (count argv)
        set currentArg to item i of argv
        if currentArg is "--name" then
            set vmName to item (i + 1) of argv
        else if currentArg is "--backend" then
            set vmBackend to item (i + 1) of argv as string
        else if currentArg is "--arch" then
            set vmArch to item (i + 1) of argv
        else if currentArg is "--source" then
            set sourcePath to POSIX file (POSIX path of item (i + 1) of argv)
        end if
    end repeat

    -- Create a new VM with the specified properties
    tell application "UTM"
        set vm to make new virtual machine with properties Â
          { backend:vmBackend, Â
            configuration:{ Â
              name:vmName, Â
              architecture:vmArch, Â
              drives:{  Â
                {source:sourcePath} Â
              } Â
            } Â
          }
    end tell
end run