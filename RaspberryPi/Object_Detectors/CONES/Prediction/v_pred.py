import numpy as np
import cv2 as cv
from imutils.video import FPS

cvNet = cv.dnn.readNetFromTensorflow('frozen_inference_graph.pb', 'graph.pbtxt')
vs=cv.VideoCapture('video.mp4')
fps = FPS().start()
while(vs.isOpened()):
	frame = vs.read()
	width  = vs.get(3) #float
	height = vs.get(4) # float
	cvNet.setInput(cv.dnn.blobFromImage(frame,size=(300, 300), swapRB=True, crop=False))
	cvOut = cvNet.forward()
	for detection in cvOut[0,0,:,:]:
	    score = float(detection[2])
	    if score > 0.3:
		left = detection[3] * cols
		top = detection[4] * rows
		right = detection[5] * cols
		bottom = detection[6] * rows
		cv.rectangle(img, (int(left), int(top)), (int(right), int(bottom)), (23, 230, 210), thickness=2)
	fps.update()
fps.stop()
print("[INFO] elapsed time: {:.2f}".format(fps.elapsed()))
print("[INFO] approx. FPS: {:.2f}".format(fps.fps()))
cv.imshow('img', frame)
cv.waitKey()
