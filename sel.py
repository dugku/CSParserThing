from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
import time

driver = webdriver.Chrome()


time.sleep(5)
driver.get("https://www.hltv.org/results?event=7391")
time.sleep(10)
elemTop = driver.find_elements(By.CLASS_NAME, "a-reset")

links_List = []

for i in elemTop:
    no = i.get_attribute("href")
    if "pgl-cs2-major-copenhagen-2024-europe-rmr-decider-qualifier" in no:
        links_List.append(no)

f = open("matches.txt", "a")

lenght = len(links_List)
for i,ele in enumerate(links_List,lenght-2):
    f.writelines(ele + "\n")
print(links_List)
driver.close()