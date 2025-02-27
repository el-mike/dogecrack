# Btcrecover-based password cracker

## Local

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

## Vast.ai
1. In `Instance Configuration`, choose custom image and paste: `michalhuras/pitbull:latest`
2. Rent a machine
3. Go to `Instances`
4. Click "connect"
5. Connect with instance via ssh

## Pitbull tool
Once inside the container, you can use pitbull tool from any directory (`pitbull` is added to PATH). Pitbull offers several commands that you can use to manage and monitor btcrecover process.
Usage:
```bash
pitbull <command> [<args>]
```
Sections below describes available commands.
### Running
```bash
pitbull run [-t <token-list>] [-f <passlist-file-path>] [-u <passlist-file-url>] [-w <wallet-string>] [--skip <skip-count>] [--length-min <number>] [--length-max <number>]
# Examples:
pitbull run -t '%a%warsaw%d\n%a%acracow%a' -w myExampleWalletExtractString --skip 10000
pitbull run -f ./my-passlist.txt -w myExampleWalletExtractString
pitbull run -u https://my-file-storage.com/myFileId -w myExampleWalletExtractString
```

When used with `-f` flag, runs btcrecover for `myExampleWalletExtractString` with `./my-passlist.txt`.

When used with `-u` flag, downloads given file (in directory that holds `pitbull.sh` file, under `passlist.txt` name), and runs btcrecover for `myExampleWalletExtractString`.

Using `--skip` provides an option to skip certain amount of passwords (helpful when a run was interrupted, and we don't want to start from scratch and loose progress on already tested passwords).

Using `--length-min` and `--length-max` allows to skip passwords that are shorter or longer than provided values. Please note that the passwords will still be counted, they will be simply skipped from testing.

Pitbull runs a new terminal session with tmux, under the name "pitbull". Because of that, you can safely close the terminal session you started the Pitbull in (including logging out from SSH), and the process will continue to run without interruption. 
You can easily re-attach to pitbull session with:
```bash
tmux a -t "pitbull"
```
To see the live progress. You can also use [Attach](#attach) command described below.

You can also use `--no-tmux` flag, to run it in your current terminal directly - this should be used for debugging only.

Btcrecover's output is continuously written to `progress_view.txt` file (including loading indicators). Refer to [Output](#output) to see how you can access it.

### Status
```bash
pitbull status
```
It prints current status based on pitbull output. It can return following statuses:
* `WAITING` - happens when container has been set up, but pitbull has not been run yet.
* `RUNNING` - happens when process is still running (i.e. pitbull process is still active).
* `SUCCESS` - happens when btcrecover stops and phrase 'Password found' is found in the output.
* `FINISHED` - any other case.

When status script returns `SUCCESS` or `FINISHED`, check the [Output](#output) to see the btcrecover's results.

Refer to the script itself for additional info.

### Output
```bash
pitbull output
# or
cat $pitbullDir/progress_view.txt
```
To get current output.

### Progress
Run:
```bash
pitbull progress
```
To get btcrecover progress, in form of `done of to_be_done`. It's the begginging of btcrecover's progress indicator output line. This command returns last recorded progress (last progress indicator line saved to file), so if the process is finished, `progress` command will output how many passwords were tried before btcrecover ended. If no progress is available (btcrecover has not started yet, or some error happened before it started) `progress` will return `NO_PROGRESS_AVAILABLE`.

This command should be use in tandem with `status` - for example, if status is `FINISHED`, but `done` is less than `to_be_done`, we know some error happened, before all passwords were checked. If they are equal, and status is still `FINISHED`, that simply means that all of the passwords were checked, but no correct one has been found. 

### Kill
Run:
```bash
pitbull kill
```
To kill the entire terminal Pitbull was run in.

### Attach
Run
```bash
pitbull attach
```
To attach to the tmux session running Pitbull. Use `Ctrl+b -> d` to detach.

### Errors
Run
```bash
pitbull errors
```
To see if there were any errors while running the Pitbull process. Errors are written to `$pitbullDir/err_log.txt` file.
