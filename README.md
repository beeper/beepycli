# Beepy CLI
Your one-stop-shop for all things Beepy, as long as all things means:
going from zero to messaging with Beeper on your device, or copying logs
from your Beepy to your computer.

## Download
You can find the latest binaries precompiled in [GitHub Actions].

After downloading, make sure to give `beepycli` permission to execute.
You can do this from the terminal by running `chmod +x beepycli`.

If you're on Linux, this is all you need!

### macOS
If you're running macOS, you've got one more step—you need a working
`libolm` installation. If you're using [Homebrew], the most popular
macOS package manager, this is as simple as running `brew install
libolm`. If you use [MacPorts], you can run `sudo port install olm`
instead.

You may also need to jump through some extra hoops when you first run
the CLI because it is not distributed through the App Store. On first
run, macOS will present you with a dialogue with two options: **Move to
Trash** or **Cancel**. Choose **Cancel**, and open **System Settings >
Privacy & Security** and scroll down until you see the option to approve
the binary. Then, run the CLI again and click **Open**.

Simple as 🥧.

## Build
Alternatively, you can build the CLI yourself by cloning the repo and
running `go build`. Building requires Go 1.20 or higher, and a `libolm`
installation.

## Usage
To log in to your Beeper account and install `gomuks`—the Beeper client
for Beepy—all you have to do is run `beepycli` from your terminal.

To copy logs from your Beepy device to your machine, run `beepycli
--logs`.

~~Made in collaboration with Shadow Wizard Money Gang.~~

[GitHub Actions]: https://github.com/beeper/beepycli/actions
[Homebrew]: https://brew.sh
[MacPorts]: https://www.macports.org
