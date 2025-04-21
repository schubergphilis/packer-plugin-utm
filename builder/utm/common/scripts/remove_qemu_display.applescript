-- remove_display.applescript
-- This script removes a display from a specified UTM virtual machine with the given display ID.
-- Usage: osascript remove_display.applescript <VM_UUID> <DISPLAY_ID>
-- Example: osascript remove_display.applescript A1B2C3 7FB247A3-DC9F-4A61-A123-0AEE1BEEC636

on run argv
  set vmId to item 1 of argv # UUID of the VM
  set displayId to item 2 of argv # ID of the display to remove

  tell application "UTM"
    -- Get the VM and its configuration
    set vm to virtual machine id vmId -- Id is assumed to be valid
    set config to configuration of vm

    -- Existing displays
    set vmDisplays to displays of config

    -- Find and remove the display with the given ID
    set updatedDisplays to {}
    repeat with display in vmDisplays
      if id of display is not displayId then
        set end of updatedDisplays to display
      end if
    end repeat

    -- Set the updated displays list
    set displays of config to updatedDisplays

    -- Save the configuration (VM must be stopped)
    update configuration of vm with config
  end tell
end run