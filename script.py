# -*- coding: utf-8 -*-
from bs4 import BeautifulSoup
import csv
from urllib.request import urlopen
import lxml
import requests

class project:
    link = ''
    contributors = []

class contributor:
    login = ''
    groups = []

projects = []
contributors = []
baseURL = "https://api.github.com/"
url = "https://api.github.com/repos/kubernetes/kubernetes/contributors"

kubernetes = project()

kubernetes.link = url

projects.append(kubernetes)

for current in requests.get(url).json():
    print(current)
    newContributor = contributor()
    #newContributor.login = current["login"]
    #for contrib in requests.get( baseURL + "users/" + newContributor.login +  "/orgs").json():
    #    print( contrib )
    #print(current["contributions"], current["login"], current["url"] )

print(contributors)

#print(requests.get(url).json()[1])



print("Hello")
