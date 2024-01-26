from selenium import webdriver
from selenium.webdriver.common.keys import Keys
from selenium.webdriver.common.by import By
from selenium.webdriver.chrome.options import Options
import time
import undetected_chromedriver as uc
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
import logging
from selenium.common.exceptions import TimeoutException
import random


options = uc.ChromeOptions()
prefs = {"download.default_directory" : "C:\\Users\\iphon\\Desktop\\DEMOProject\\All_Rars"}
options.add_experimental_option("prefs", prefs )
driver = uc.Chrome(options=options)

def Make_links():
    f = open("matches.txt", "r")

    links = []

    for i in f:
        links.append(i[:-1])

    get_page(links)
    
def get_page(listThing):
    counter = 0
    for link in listThing: 
        first = random.randint(6, 30)
        second = random.randint(3, 8)
        third = random.randint(60, 140)
        print(link + "First:" +  str(first) + " Second: " + str(second) + " Third: " + str(third))
        get_demo(link, first, second, third)
        with open ("matches.txt", "r+") as f:
            d = f.readlines()
            f.seek(0)
            for i in d:
                if i != link:
                    f.write(i)
            f.truncate()
        

        
def get_demo(matchLink, first, second, third):
    driver.get(matchLink)
    time.sleep(first)
    WebDriverWait(driver, 10).until(EC.element_to_be_clickable((By.CLASS_NAME, "stream-box")))
    clickAble = driver.find_element(By.CLASS_NAME, "stream-box")
    print("Clicking")
    time.sleep(second)
    clickAble.click()
    time.sleep(third)


Make_links()



time.sleep(15)
driver.quit()
