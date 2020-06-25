# USAGE
# python real_time_object_detection.py --prototxt MobileNetSSD_deploy.prototxt.txt --model MobileNetSSD_deploy.caffemodel

# import the necessary packages
#from imutils.video import VideoStream
from imutils.video import FPS
import numpy as np
import argparse
import imutils
import time
import cv2
import paho.mqtt.client as paho
import base64
import json 

broker="localhost" #host
port=1883
light = paho.Client("Cones") 
stop = False
def on_connect(light, obj, flags, rc):
    print("rc: {}".format(rc))
    light.subscribe("light/CONE/stop", 0)
    
def on_publish(light,userdata,result):   
	light.subscribe("light/CONE/stop", 0)          #create function for callback
       

def on_message(light, userdata, msg):
    print('New message recieved -> topic:{} - payload:{}'.format(msg.topic, msg.payload))
    if msg.payload == b'stop':
    	global stop 
    	stop = True
    	print('Execution interrupted')

light.on_publish = on_publish    
light.on_message = on_message    
light.on_connect = on_connect                  #assign function to callback
light.connect(broker,port)   

light.loop_start()
# initialize the list of class labels MobileNet SSD was trained to
# detect, then generate a set of bounding box colors for each class
CLASSES = ["cone"]
#CLASSES = ["person","bicycle","car","motorbike","aeroplane","bus","train","truck","boat","traffic light","fire hydrant","stop sign","parking meter","bench","bird","cat","dog","horse","sheep","cow","elephant","bear","zebra","giraffe","backpack","umbrella","handbag","tie","suitcase","frisbee","skis","snowboard","sports ball","kite","baseball bat","baseball glove","skateboard","surfboard","tennis racket","bottle","wine glass","cup","fork","knife","spoon","bowl","banana","apple","sandwich","orange","broccoli","carrot","hot dog","pizza","donut","cake","chair","sofa","pottedplant","bed","diningtable","toilet","tvmonitor","laptop","mouse","remote","keyboard","cell phone","microwave","oven","toaster","sink","refrigerator","book","clock","vase","scissors","teddy bear","hair drier","toothbrush"]
COLORS = np.random.uniform(0, 255, size=(len(CLASSES), 3))

# load our serialized model from disk
print("[INFO] loading model...")
net = cv2.dnn.readNetFromTensorflow("frozen_inference_graph.pb", "graph.pbtxt")

# initialize the video stream, allow the cammera sensor to warmup,
# and initialize the FPS counter
print("[INFO] starting video stream...")
vs = cv2.VideoCapture('video.mp4')
time.sleep(2.0)
fps = FPS().start()

max_confidence = ( 0, 0, 0)

# loop over the frames from the video stream
while True:
	# grab the frame from the threaded video stream and resize it
	# to have a maximum width of 400 pixels
	ret, frame = vs.read()
	frame = imutils.resize(frame, width=400)

	# grab the frame dimensions and convert it to a blob
	(h, w) = frame.shape[:2]
	blob = cv2.dnn.blobFromImage(frame,size=(300, 300), swapRB=True, crop=False)

	# pass the blob through the network and obtain the detections and
	# predictions
	net.setInput(blob)
	detections = net.forward()

	list_of_detections = [] #Restart the list of detections
	
	# loop over the detections
	for i in np.arange(0, detections.shape[2]):
		# extract the confidence (i.e., probability) associated with
		# the prediction
		confidence = detections[0, 0, i, 2]

		# filter out weak detections by ensuring the `confidence` is
		# greater than the minimum confidence
		if confidence > 0.4:
			# extract the index of the class label from the
			# `detections`, then compute the (x, y)-coordinates of
			# the bounding box for the object
			idx = int(detections[0, 0, i, 1])-1
			box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
			(startX, startY, endX, endY) = box.astype("int")
			distance=0.75*200/(endY-startY)
			# draw the prediction on the frame
			label = "{}: {:.2f} m".format("Distance" ,distance)
			cv2.rectangle(frame, (startX, startY), (endX, endY), COLORS[idx], 2)
			y = startY - 15 if startY - 15 > 15 else startY + 15
			cv2.putText(frame, label, (startX, y),	cv2.FONT_HERSHEY_SIMPLEX, 0.5, COLORS[idx], 2)
			#we add every new detections to the list
			list_of_detections.append([confidence, idx, startX, startY, endX, endY, distance])		
	
	count = 0 #restart count for new image
	d_min=100
	list_lows=[]
	rand=False
	line=False
	if list_of_detections != []:
		for i in list_of_detections:
			list_lows.append(i[5])
			if i[6]<d_min:
				d_min=i[6]
		for i in list_lows:
			if max(list_lows)-i<15:
				count+=1
		if max(list_lows)-min(list_lows)>40 or count>0:
			rand=True	
		if count>=4:
			line=True		
	
	#MQTT publish
	if line:
		light.publish("light/CONE/detection","0/"+str(d_min)) #publish
	elif rand:
		light.publish("light/CONE/detection","5/"+str(d_min))
	else: 
		light.publish("light/CONE/detection","30")
	# show the output frame
	cv2.imshow("Frame", frame)
	key = cv2.waitKey(1) & 0xFF

	if stop:
		break

	# if the `q` key was pressed, break from the loop
	if key == ord("q"):
		break

	# update the FPS counter
	fps.update()
# stop the timer and display FPS information
fps.stop()
print("[INFO] elapsed time: {:.2f}".format(fps.elapsed()))
print("[INFO] approx. FPS: {:.2f}".format(fps.fps()))

# do a bit of cleanup
cv2.destroyAllWindows()
