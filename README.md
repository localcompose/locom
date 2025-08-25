# locom

**locom** (short for "local compose") is a command-line tool to help manage local Docker Compose stacks in a clean, self-hosted, offline-friendly way.

## Commands

- `locom new`: Create a minimal config file in the current folder#.

## Disclaimer

The installation, test, and cleanup steps described above have been verified against this release in a "happy flow."
However, we cannot guarantee that no issues will arise on your specific system.
Use the cleanup commands only if you installed this exact version, as future releases may change file names, locations, or behavior.

Proceed at your own discretion.

## Happy path

Version `0.0.3-poc`.

### Prerequisites

* `curl`
* `docker`, `docker compose`
* `openssl`

### Download and install



Prerequisite: `tar` 

<details>
<summary>linux</summary>

```sh
curl -LO https://github.com/localcompose/locom/releases/download/0.0.3-poc/locom_linux_amd64.tar.gz
tar -xvzf locom_linux_amd64.tar.gz
chmod +x locom
sudo mv -f locom /usr/local/bin/

# ---- Cleanup (use with caution)
# Run only if you installed this version. Future versions may differ.
sudo rm -f /usr/local/bin/locom
rm -f locom_linux_amd64.tar.gz
```

</details>

<details>
<summary>darwin</summary>

```sh
curl -LO https://github.com/localcompose/locom/releases/download/0.0.3-poc/locom_darwin_amd64.tar.gz
tar -xvzf locom_darwin_amd64.tar.gz
chmod +x locom
sudo mv -f locom /usr/local/bin/

# ---- Remove after testing / cleanup
sudo rm -f /usr/local/bin/locom
rm -f locom_darwin_amd64.tar.gz
```
</details>

<details>
<summary>windows</summary>


> ⚠️ Run the following commands in an **Administrator PowerShell** or **Administrator Git Bash** session,  
> since moving binaries into `%SystemRoot%\System32` requires elevated privileges.

> ⚠️ Precaution: You need to be an **Administrator** on your system to install into `%SystemRoot%\System32`.  
> The commands below use `runas /user:%USERNAME%` to ensure execution with your account.  
> Depending on your UAC settings, you may be prompted for elevation.


<details>
<summary>Git Bash for Windows</summary>

```zsh
curl -LO https://github.com/localcompose/locom/releases/download/0.0.3-poc/locom_windows_amd64.tar.gz
tar -xvzf locom_windows_amd64.tar.gz

# Move to System32 (always in PATH) via runas
winpty runas /user:$USERNAME "cmd /c move /Y locom.exe %SystemRoot%\System32\"

# ---- Cleanup (use with caution)
# Run only if you installed this version. Future versions may differ.
winpty runas /user:$USERNAME "cmd /c del /Q %SystemRoot%\System32\locom.exe"
rm -f locom_windows_amd64.tar.gz
```
</details>

<details>
<summary>PowerShell</summary>

```powershell
# Download and extract
curl -LO https://github.com/localcompose/locom/releases/download/0.0.3-poc/locom_windows_amd64.tar.gz
tar -xvzf locom_windows_amd64.tar.gz

# Move to System32 (always in PATH) via runas
runas /user:$env:USERNAME "powershell -Command Move-Item -Force .\locom.exe $env:SystemRoot\System32\"

# ---- Cleanup (use with caution)
# Run only if you installed this version. Future versions may differ.
runas /user:$env:USERNAME "powershell -Command Remove-Item -Force $env:SystemRoot\System32\locom.exe"
Remove-Item -Force .\locom_windows_amd64.tar.gz
```
</details>

</details>

### Setup

```sh
export LOCOM_STAGE=$(locom version | cut -d' ' -f3)
echo $LOCOM_STAGE
locom init $LOCOM_STAGE
cd $LOCOM_STAGE
```

<details>
<summary>linux (Ubuntu)</summary>

Tested on Ubuntu, but may work on other distros without change.

Some commands need sudo on id docker installed by snap.

```sh
sudo $(which locom) network
locom hosts --verify
locom proxy
locom cert selfsigned setup
locom cert selfsigned trust

cd proxy
sudo docker compose up -d
```

</details>

<details>
<summary>darwin</summary>

```sh
locom network
locom hosts --verify
locom proxy
locom cert selfsigned setup
locom cert selfsigned trust

cd proxy
docker compose up -d
```
</details>

<details>
<summary>windows</summary>

```sh
locom network
locom hosts --verify
locom proxy
locom cert selfsigned setup
locom cert selfsigned trust

cd proxy
docker compose up -d
```
</details>

## test

```sh
open https://proxy.locom.self
```

<details>
<summary>Advanced test and troublshooting</summary>

```sh
openssl s_client -connect proxy.locom.self:443 -servername proxy.locom.self </dev/null 2>/dev/null   | grep -E "subject=|issuer="

curl -I https://proxy.locom.self
curl https://proxy.locom.self
curl -L https://proxy.locom.self
curl -s -o /dev/null -w "%{http_code}\n" https://proxy.locom.self
```

</details>


## cleanup

<details>
<summary>linux (Ubuntu)</summary>

Tested on Ubuntu, but may work on other distros without change.

```sh
sudo docker compose down
cd ..

locom cert selfsigned untrust
```

</details>

<details>
<summary>darwin</summary>

```sh
docker compose down
cd ..

locom cert selfsigned untrust
```
</details>

<details>
<summary>windows</summary>

```sh
docker compose down
cd ..

locom cert selfsigned untrust
```
</details>
