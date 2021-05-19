# coding:utf-8



#
# This script fetches all the public channels from thingspeak 
# we use the data in random device generator
# #

import requests
import json
import os
import sys
import re
import hashlib


apiURL = 'https://api.thingspeak.com/channels/public.json?page='


#-------------------#

def storeData(jsonFile, data):
    with open(jsonFile, 'w+') as f:
        f.seek(0)
        f.write(json.dumps(data, indent=2))
        f.truncate()

#-------------------#

if __name__ == '__main__':

    dataList    = []
    namesList   = []

    page = 0
    while True:

        page += 1
        url = apiURL + str( page)
        try:
            f = requests.get(url, headers={'User-Agent': 'Mozilla'})
            body = f.text
        except Exception as e:
            print(url)
            print("Error: ", __file__, "Line: ",
                sys._getframe().f_lineno, "\n\t", e)
            break


        data = json.loads(body)
        if len( data['channels']) == 0:
            print( "I hit the final page")
            break
        
        dataList += data['channels']

        for item in data['channels']:
            namesList.append( item['name'])

        storeData( './data.json', dataList)
        storeData( './names.json', namesList)

        print( "Page ", page, " done")
    
    print( "All done :)")

                

