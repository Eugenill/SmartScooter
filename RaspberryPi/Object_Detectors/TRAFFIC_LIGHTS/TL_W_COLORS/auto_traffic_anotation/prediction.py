import cv2 as cv
import os
from shutil import copyfile
from xml.etree import ElementTree as et
import numpy as np
import copy

labels = ["person", "bicycle", "car", "motorcycle", "airplane", "bus",
    "train", "truck", "boat", "traffic light", "fire hydrant", 
   "street sign","stop sign", "parking meter", "bench", "bird", "cat", "dog", 
   "horse", "sheep", "cow", "elephant", "bear", "zebra", "giraffe","hat", 
   "backpack", "umbrella", "shoe","eye glasses","handbag", "tie", "suitcase", "frisbee", 
   "skis", "snowboard", "sports ball", "kite", "baseball bat", 
   "baseball glove", "skateboard", "surfboard", "tennis racket", 
   "bottle","plate", "wine glass", "cup", "fork", "knife", "spoon", "bowl", 
   "banana", "apple", "sandwich", "orange", "broccoli", "carrot", 
   "hot dog", "pizza", "donut", "cake", "chair", "couch", 
   "potted plant", "bed","mirror", "dining table","window","desk", "toilet", "door","tv", "laptop", 
   "mouse", "remote", "keyboard", "cell phone", "microwave", "oven", 
   "toaster", "sink", "refrigerator","blender", "book", "clock", "vase", 
   "scissors", "teddy bear", "hair drier", "toothbrush","hair brush"]
COLORS = np.random.uniform(0, 255, size=(len(labels), 3))
cvNet = cv.dnn.readNetFromTensorflow('frozen_inference_graph.pb', 'graph.pbtxt')
pose ='Unspecified'
i="train"
for filename in os.listdir("/home/eugeni/Desktop/SmartScooter/traffic_lights/train/"):
	if filename.endswith(".jpg"):
		fil=filename	
		xmlfile=fil.split(".jpg")[0]+".xml"	
		copyfile("/home/eugeni/Desktop/SmartScooter/traffic_lights/example.xml", "/home/eugeni/Desktop/SmartScooter/traffic_lights/train/"+xmlfile)	
		img = cv.imread(filename)
		h = img.shape[0]
		w = img.shape[1]
		print(xmlfile)
		for filename2 in os.listdir("/home/eugeni/Desktop/SmartScooter/traffic_lights/train/"):
			if filename2==xmlfile:		
				doc = et.parse(xmlfile)
				root = doc.getroot()
				name=root[1]
				name.text='train/'+filename	
				path = root[2]
				path.text="/content/drive/My Drive/SMART-SCOOTER/ML/models/research/object_detection/train_xml_annotations/"+xmlfile
				folder=root[0]
				folder.text="/traffic_lights/train/"
				size=root[4]
				size[0].text=str(w)
				size[1].text=str(h)
				size[2].text="3"
				doc.write(xmlfile)
	
				cvNet.setInput(cv.dnn.blobFromImage(img, size=(300, 300), swapRB=True, crop=False))
				detections = cvNet.forward()
				c=0
				for i in np.arange(0, detections.shape[2]):
				#for detection in cvOut[0,0,:,:]:
					#score = float(detection[2])
					confidence = detections[0, 0, i, 2]
					
					if confidence > 0.3:
						#idx = int(detection[1])
						#left = detection[3] * w
						#top = detection[4] * h
						#right = detection[5] * w
						#bottom = detection[6] * h
						idx = int(detections[0, 0, i, 1])-1
						box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
						(startX, startY, endX, endY) = box.astype("int")
						if idx==9:
							print('1')
							if c==0:
								obje=root[6]
								obje[0].text="traffic light"
								obje[4][0].text=str(startX)
								obje[4][1].text=str(startY)
								obje[4][2].text=str(endX)
								obje[4][3].text=str(endY)
								c+=1
							else:
	                                                        annotation=et.Element('annotation')
	                                                        obj = et.Element('object')
								root.append(obj)
	                                                        et.SubElement(obj, 'name').text = 'traffic light'
	                                                        et.SubElement(obj, 'pose').text = pose
	                                                        et.SubElement(obj, 'truncated').text = '0'
	                                                        et.SubElement(obj, 'difficult').text = '0'

	                                                        bbox = et.SubElement(obj, 'bndbox')
	                                                        et.SubElement(bbox, 'xmin').text = str(startX)
	                                                        et.SubElement(bbox, 'ymin').text = str(startY)
	                                                        et.SubElement(bbox, 'xmax').text = str(endX)
	                                                        et.SubElement(bbox, 'ymax').text = str(endY)
		
							#cv.rectangle(img, (startX, startY), (endX, endY),COLORS[idx], 2)
							doc.write(xmlfile)
		os.rename("/home/eugeni/Desktop/SmartScooter/traffic_lights/train/"+xmlfile, "/home/eugeni/Desktop/SmartScooter/traffic_lights/train_xml_annotations/"+xmlfile)					
							
#cv.imshow("TL", img)			
cv.waitKey()
