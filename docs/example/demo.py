import requests

if __name__ == '__main__':
    url = "https://api.github.com/repos/darumatic/pyx?page=1&per_page=100"
    resp = requests.get(url=url)
    data = resp.json()
    print("Pyx repository has {} stars. Why not give it a star? https://github.com/darumatic/pyx.git".format(data["stargazers_count"]))
