# pyx

Single command to run python3 script anywhere. pyx = install python + install git + checkout repository + run your script.


![Alt text](docs/example.gif?raw=true "pyx example")


## Example usage

1) Run git repository python scripts

   ```
   $ pyx --branch=master https://github.com/darumatic/pyx  docs/example/hello.py
   ```
   
   or 
   
   ```
   $ pyx git@github.com:darumatic/pyx.git docs/example/hello.py
   ```
   
   For github repositories, we could also use the short repository name.
   
   ```
   $ pyx darumatic/pyx docs/example/hello.py
   ```
   
   
2) Run local python scripts
   
   ```
   $ pyx ./darumatic/pyx docs/example/hello.py
   ```

## Install pyx

### Windows 10 

1. Download the binary archive file. [https://github.com/darumatic/pyx/releases/download/1.0.4/pyx_1.0.4_Windows_x86_64.zip](https://github.com/darumatic/pyx/releases/download/1.0.4/pyx_1.0.4_Windows_x86_64.zip)
2. Unzip the archive file.
3. Open terminal, run ```./pyx --version```

### Linux

1. Open terminal, and run the following command. 
```bash
sudo wget -c https://github.com/darumatic/pyx/releases/download/1.0.4/pyx_1.0.4_Linux_x86_64.tar.gz -O - | sudo tar -xz -C /usr/local/bin
```
3. Run ```./pyx --version```


### MacOS

1. Open terminal, run ```brew tap darumatic/tap && brew install pyx```

## How it works
1. pyx will first check if python3 installed in local.
2. if not, pyx will download the static build python and extract to ~/.pyx/python
3. pyx will check if the repository exists in local. 
4. if not, pyx clone the repository. 
5. If the repository exists, pyx will pull the repository to keep the script updated. 
6. pyx call python3 to run the script.

## Why we need pyx
pyx makes sharing python scripts among the teams easier. 

- Single command to run python scripts, easy to document.
- Keep the script updated.
