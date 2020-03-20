import os
from xml.etree import ElementTree as et

for filename in os.listdir("/home/smartscooter/Desktop/CONES_V1/test"):
        if filename.endswith(".xml"): 
                print(type(filename),filename)  
                doc = et.parse(filename)
                root = doc.getroot()
                path = root[2]
                print(path.text)
                print(filename.split(".xml"))        
                path.text="/content/drive/My Drive/SMART-SCOOTER/ML/models/research/object_detection/images/CONES_V1/test/"+filename.split(".xml")[0]+".jpg"
                print(path.text)
                doc.write(filename)
