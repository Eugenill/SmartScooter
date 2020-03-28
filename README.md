# SmartScooter Object Detection files

## Create the model

To create the model, git clone this repository to SmartScooter1996 google account in Colab, and follow the instructions in one of this two ipynb:
- [CONE](CONES/CONE_object_detctor.ipynb)
- [TL](TRAFFIC_LIGHTS/TL_object_detector.ipynb)

## Inference time
1. Once you have created the graph.pbtxt with the frozen_inference_graph.pb, and replaced "AddV2" for "Add" in the graph.pbtxt. We have to export the frozen_inference_graph.pb and the graph.pbtxt (i.e. export all the new_model in zip format). 
    - We will be posting the new_models ready for inference here:
    
        [CONE](CONES/new_models)
        
        [TL](TRAFFIC_LIGHTS/TL_W_COLORS/new_models)
2. Use prediction scripts woth OpenCv newtwork, remember to use the latest OpenCv version:
    - [Cones](CONES/Prediction)
    - [TL](TRAFFIC_LIGHTS/Prediction)
    1. pred.py for photo inference 
    2. vide_object_detection.py for video inference
3. Together with the Intel Movidius NCS we can obtain more speed in the RaspberryPi. Follow this installation instructions: https://www.pyimagesearch.com/2019/04/08/openvino-opencv-and-movidius-ncs-on-the-raspberry-pi/

4. Together with multiprocessing scripts we can further optimize the model. (Soon)

5. Modify prediction script to send the detection messages throw MQTT to the Nodered App for the Scooter implementation (Soon)

## References:

- xml_to_csv.py and generate_tfrecord.py from: https://github.com/datitran/raccoon_dataset
- bosch_to_pascal.py from: https://github.com/bosch-ros-pkg/bstld
- TensorFlow Object Dtection API from: https://github.com/tensorflow/models
- PyImagesearch: https://www.pyimagesearch.com/
- OpenCV latest version
