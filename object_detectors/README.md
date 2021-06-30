# Smart Scooter Object Detection files

## Create the model

To create the model, git clone this repository to your Google Colab, and follow the instructions of one of this two ipynb:
- [CONE](CONES/CONES_object_detector.ipynb)
- [TL](TRAFFIC_LIGHTS/TL_object_detector.ipynb)

## Inference time
1. Prepare models for OpenCV
    Once you have created the graph.pbtxt with the frozen_inference_graph.pb, and replaced "AddV2" for "Add" in the graph.pbtxt. We have to export the frozen_inference_graph.pb and the graph.pbtxt (i.e. export all the new_model in zip format). 
    - Here are the new_models ready for inference:
        [CONE](CONES/new_models)
        [Traffic Ligts](TRAFFIC_LIGHTS/TL_W_COLORS/new_models)
2. Prediction scripts with OpenCV
   Use prediction scripts with OpenCv newtwork, remember to use the latest OpenCV version (here used 4.4.0, https://opencv.org/releases/):
    - [Cones](Prediction/CONES)
    - [Traffic Lights](Prediction/TRAFFIC_LIGHTS)
  
3. Install OPENVINO
   Together with the Intel Movidius NCS we can obtain more speed in the RaspberryPi 4B (4GB RAM). Follow this installation instructions: https://www.pyimagesearch.com/2019/04/08/openvino-opencv-and-movidius-ncs-on-the-raspberry-pi/

4. Optimize model
   Once you have completed the OPENVINO installation, we can further optimize the model (Same steps for CONE model).
    - Use `mo_tf.py --input_model [TL_v5/frozen_infrence_graph.pb](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/frozen_infrence_graph.pb) --tensorflow_object_detection_api_pipeline_config [pipeline.config](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/pipeline.config) --tensorflow_use_custom_operations_config [Download this one](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/ssd_support_api_v1.15.json) --generate_deprecated_IR_V7 --data_type FP16 --batch 1`
    - Once we have the [.xml and .bin files](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/OPENVINO/IR7,FP16)
    - Download the last version (2020.1) Openvino Raspberry, following this [steps](https://www.pyimagesearch.com/2019/04/08/openvino-opencv-and-movidius-ncs-on-the-raspberry-pi/)
    - Use this [prediction script](Prediction/TRAFFIC_LIGHTS/OPENVINO/async_pred.py) to predict with multiprocessing with NCS conected, (BRG2RGB must be added!)
    
5. Adding MQTT
   Finally use this script [Object detection + MQTT](Prediction/TRAFFIC_LIGHT/video_detection_mqtt.py) to send the detection messages throw MQTT to the Node-Red dashboard for the Scooter implementation.

## Utils
This [script](utils/auto_annotation_program.py) aims to save time when labelling images for object detection. Using a pre-trained model and the pygame library, with this script you can label your images quicker and get a csv with the annotations.

## References:
- xml_to_csv.py and generate_tfrecord.py from: https://github.com/datitran/raccoon_dataset
- bosch_to_pascal.py from: https://github.com/bosch-ros-pkg/bstld
- TensorFlow Object Dtection API from: https://github.com/tensorflow/models
- PyImagesearch: https://www.pyimagesearch.com/