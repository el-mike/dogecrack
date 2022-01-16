# Btcrecovery-based password cracker

## Local

Requires NVIDIA container runtime: [Installation](https://docs.nvidia.com/datacenter/cloud-native/container-toolkit/install-guide.html)

```
docker build -t pitbull .
docker run --runtime=nvidia -ti pitbull
```

## Vast.ai
1. In `Instance Configuration`, choose custom image and paste: `michalhuras/pitbull:4.0` (make sure you are using the newest version)
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
pitbull run [-f <file-url>] [-w <wallet-string>]
# Example:
pitbull run -f https://my-file-storage.com/myFileId -w myExampleWalletExtractString
```

Downloads a given file, and runs btcrecover for `myExampleWalletExtractString`.
Note that the file specfied in arguments will be saved in directory  that holds `pitbull.sh` file.

Pitbull runs a new terminal session with tmux, under the name "pitbull". Thanks to that, you can safely close the terminal session you started the Pitbull in (including logging out from SSH), and the process will continue to run without interruption. 
You can easily re-attach to pitbull session with:
```bash
tmux a -t "pitbull"
```
To see the live progress.

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
To get btcrecover progress, in form of `done of to_be_done`. It's the begginging of btcrecover's progress indicator output line. If the recovery process itself has not started yet, this command will return `Progress not found: $lastLine`, where `$lastLine` is the last line of `progress_view.txt` output file.

### Kill
Run:
```bash
pitbull kill
```
To kill the entire terminal Pitbull was run in.
