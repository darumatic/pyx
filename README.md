# pyx

Single command to run python3 script anywhere. pyx = install python + install git + checkout repository + run your script.


![Alt text](docs/example.gif?raw=true "pyx example")


## Example usage

1) Run git repository python scripts

   ```
   $ pyx https://github.com/darumatic/pyx scripts/hello.py
   ```
   
   or 
   
   ```
   $ pyx git@github.com:darumatic/pyx.git scripts/hello.py
   ```
   
   For github repositories, we could also use the short repository name.
   
   ```
   $ pyx darumatic/pyx scripts/hello.py
   ```
   
3) Run http python scripts

   ```
   $ pyx https://raw.githubusercontent.com/darumatic/pyx/master/scripts/hello.py
   ```
   
4) Run local python scripts
   
   ```
   $ pyx hello.py
   ```
