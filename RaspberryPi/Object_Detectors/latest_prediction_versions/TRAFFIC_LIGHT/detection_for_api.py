# USAGE
# python real_time_object_detection.py --prototxt MobileNetSSD_deploy.prototxt.txt --model MobileNetSSD_deploy.caffemodel

import cv2
import imutils
import numpy as np
import paho.mqtt.client as paho
import time
# import the necessary packages
from imutils.video import FPS

broker = "maqiatto.com"  # host
port = 1883
username = "eugeni.llagostera@gmail.com"
password = "asdf1234"
light = paho.Client("Traffic_lights")
stop = False
slow = False
preTopic = "eugeni.llagostera@gmail.com/"


def on_connect(light, obj, flags, rc):
    print("rc: {}".format(rc))
    light.subscribe("light/TL/stop", 0)
    light.subscribe("light/TL/speed", 0)


def on_publish(light, userdata, result):
    light.subscribe("light/TL/stop", 0)  # create function for callback
    light.subscribe("light/TL/speed", 0)


def on_message(light, userdata, msg):
    print('New message recieved -> topic:{} - payload:{}'.format(msg.topic, msg.payload))
    global slow
    if msg.payload == b'stop':
        global stop
        stop = True
        print('Execution interrupted')
    elif msg.payload == b"slow":

        slow = True
        print("Riding slowly")
    elif msg.payload == b"fast":

        slow = False
        print("Riding fast")


light.on_publish = on_publish
light.on_message = on_message
light.on_connect = on_connect  # assign function to callback
light.username_pw_set(username, password)
light.connect(broker, port)

light.loop_start()
# initialize the list of class labels MobileNet SSD was trained to
# detect, then generate a set of bounding box colors for each class
CLASSES = ["off", "green", "yellow", "red"]
# CLASSES = ["person","bicycle","car","motorbike","aeroplane","bus","train","truck","boat","traffic light",
# "fire hydrant","stop sign","parking meter","bench","bird","cat","dog","horse","sheep","cow","elephant","bear",
# "zebra","giraffe","backpack","umbrella","handbag","tie","suitcase","frisbee","skis","snowboard","sports ball",
# "kite","baseball bat","baseball glove","skateboard","surfboard","tennis racket","bottle","wine glass","cup","fork",
# "knife","spoon","bowl","banana","apple","sandwich","orange","broccoli","carrot","hot dog","pizza","donut","cake",
# "chair","sofa","pottedplant","bed","diningtable","toilet","tvmonitor","laptop","mouse","remote","keyboard",
# "cell phone","microwave","oven","toaster","sink","refrigerator","book","clock","vase","scissors","teddy bear",
# "hair drier","toothbrush"]
COLORS = np.random.uniform(0, 255, size=(len(CLASSES), 3))

# load our serialized model from disk
print("[INFO] loading model...")
net = cv2.dnn.readNetFromTensorflow("frozen_inference_graph.pb", "graph.pbtxt")

# initialize the video stream, allow the cammera sensor to warmup,
# and initialize the FPS counter
print("[INFO] starting video stream...")
vs = cv2.VideoCapture('../../TRAFFIC_LIGHTS/Prediction/video.mp4')
time.sleep(2.0)
fps = FPS().start()
state = 1  # green
relevant_detection = []
max_confidence = (0.7, 0.2)
# loop over the frames from the video stream
while True:
    # grab the frame from the threaded video stream and resize it
    # to have a maximum width of 400 pixels
    ret, frame = vs.read()
    frame = imutils.resize(frame, width=400)

    # grab the frame dimensions and convert it to a blob
    (h, w) = frame.shape[:2]
    blob = cv2.dnn.blobFromImage(frame, size=(300, 300), swapRB=True, crop=False)

    # pass the blob through the network and obtain the detections and
    # predictions
    net.setInput(blob)
    detections = net.forward()

    list_of_detections = []  # Restart the list of detections
    count = True  # restart count for new image
    d_min = 10000
    p_max = 0
    k = 0
    l = 0
    idx = 0
    if slow:
        l = 1
    # loop over the detections
    for i in np.arange(0, detections.shape[2]):
        # extract the confidence (i.e., probability) associated with
        # the prediction
        confidence = detections[0, 0, i, 2]
        # filter out weak detections by ensuring the `confidence` is
        # greater than the minimum confidence
        if confidence > max_confidence[l]:
            # extract the index of the class label from the
            # `detections`, then compute the (x, y)-coordinates of
            # the bounding box for the object
            idx = int(detections[0, 0, i, 1]) - 1
            box = detections[0, 0, i, 3:7] * np.array([w, h, w, h])
            (startX, startY, endX, endY) = box.astype("int")
            distance = 15 * 10 / (endY - startY)

            # we add every new detections to the list
            list_of_detections.append([confidence, idx, startX, startY, endX, endY, distance])
    # now we have to compare the las detection chosen as relevant ot the ones detected now
    list_of_detections.sort()
    list_of_detections.reverse()  # in order of more confidence
    if list_of_detections:
        if not slow:
            if not relevant_detection:
                relevant_detection = list_of_detections[
                    0]  # if not finded any detection similar to the last one we take the one with most confidence
            for i in list_of_detections:
                if ((relevant_detection[2] - 10) <= i[2] <= (relevant_detection[2] + 10)) and (
                        (relevant_detection[3] - 10) <= i[3] <= (relevant_detection[3] + 10)) and (
                        (relevant_detection[4] - 10) <= i[4] <= (relevant_detection[4] + 10)) and (
                        (relevant_detection[5] - 10) <= i[5] <= (relevant_detection[5] + 10)):
                    relevant_detection = i
                    count = False
            if count:
                for i in list_of_detections:
                    if i[6] < d_min:
                        d_min = i[6]
                        l = i
                    if i[0] > p_max:
                        p_max = i[0]
                        k = i
                    if l[0] / l[6] > k[0] / k[6]:
                        relevant_detection = l
                    else:
                        relevant_detection = k
        elif slow:
            if not relevant_detection:
                relevant_detection = list_of_detections[
                    0]  # if not finded any detection similar to the last one we take the one with most confidence
            for i in list_of_detections:
                if abs(relevant_detection[6] - i[6]) <= 5 and relevant_detection[6] - i[6] < d_min:
                    d_min = i[6]
                    l = i
            if l != 1:
                relevant_detection = l

    # crop_img = frame[startY:endY, startX:endX]
    # if len(crop_img) != 0:
    # cv2.imwrite("crop_img.jpg", crop_img)

    # draw the prediction on the frame
    # label = "{}: {:.2f}%".format(CLASSES[relevant_detection[1]] ,relevant_detection[0] * 100)
    cv2.rectangle(frame, (relevant_detection[2], relevant_detection[3]), (relevant_detection[4], relevant_detection[5]),
                  COLORS[idx], 2)
    # y = startY - 15 if startY - 15 > 15 else startY + 15
    # cv2.putText(frame, label, (startX, y),	cv2.FONT_HERSHEY_SIMPLEX, 0.5, COLORS[relevant_detection[1]], 2)

    # image_text = base64.b64decode(crop_img)
    # create client object
    # establish connection
    if relevant_detection[1] != state:
        # MQTT publish
        light.publish(preTopic + "RP1_detection",
                      str(relevant_detection[0]) + '/' + CLASSES[relevant_detection[1]] + "/" + str(
                          relevant_detection[6]))  # publish
        # light.publish("light/TL/class",relevant_detection[1])                   #publish
        # light.publish("light/TL/box",str(relevant_detection[2])+'/'+str(relevant_detection[3])+'/
        # '+str(relevant_detection[4]) + '/' + str(relevant_detection[5])                
        # #publish
        # light.publish("light/TL/class",relevant_detection)                   #publish
        # light.publish("light/TL/image",crop_img)
        state = relevant_detection[1]
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
