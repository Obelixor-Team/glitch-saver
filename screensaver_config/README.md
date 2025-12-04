# Terminal Screensaver Configuration

This folder contains the files needed to set up the `glitch-saver` as a
terminal screensaver that activates after a period of inactivity.

## Files

- `saver.sh`: The main script that sets the inactivity timeout and runs the
`glitch-saver` application.
- `bash_config_snippet.sh`: A snippet to be added to your `~/.bashrc` file.
- `zsh_config_snippet.zsh`: A snippet to be added to your `~/.zshrc` file.
- `fish_config_snippet.fish`: A snippet to be added to your
`~/.config/fish/config.fish` file.

## Instructions

1. **Edit `saver.sh`**:
    Open the `saver.sh` file and change the `SAVER_APP` variable to the full path
of your `glitch-saver` executable. For example:

```bash
SAVER_APP="/home/user/projects/glitch-saver/glitch-saver"
```

    You can also change the `TMOUT` variable to your desired inactivity timeout
in seconds.

1. **Update Your Shell's Configuration**:
    Choose the appropriate snippet file for your shell and add its content to
your shell's configuration file.

- **For `bash`**:
        Copy the content of `bash_config_snippet.sh` and paste it at the end of
    your `~/.bashrc` file.
    **Important**: Remember to replace `/path/to/screensaver_config` with the
actual path to this `screensaver_config` folder.

- **For `zsh`**:
        Copy the content of `zsh_config_snippet.zsh` and paste it at the end of
    your `~/.zshrc` file.
    **Important**: Remember to replace `/path/to/screensaver_config` with the
actual path to this `screensaver_config` folder.

- **For `fish`**:
        Copy the content of `fish_config_snippet.fish` and paste it at the end of
    your `~/.config/fish/config.fish` file.
    **Important**: Remember to replace `/path/to/screensaver_config` with the
actual path to this `screensaver_config` folder.
    **Note**: The `fish` shell does not support the `TMOUT` variable, so the
inactivity timeout feature will not work. This will only allow you to manually
source the script to start the screensaver.

1. **Restart Your Terminal**:
    For the changes to take effect, you need to either restart your terminal or
run `source ~/.bashrc` (or the equivalent for your shell).

Now, any new terminal you open will have the screensaver enabled.
