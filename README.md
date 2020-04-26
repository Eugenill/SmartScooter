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

4. Together with multiprocessing scripts we can further optimize the model (Same steps for CONE model).
    - Download the Openvino 2020.1 Linux Package (install all needed dependences)
    - Use mo_tf.py, --input_model [TL_v5/frozen_infrence_graph.pb](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/frozen_infrence_graph.pb) --tensorflow_object_detection_api_pipeline_config [pipeline.config](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/pipeline.config) --tensorflow_use_custom_operations_config [Download this one](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/ssd_support_api_v1.15.json) --generate_deprecated_IR_V7 --data_type FP16 --batch 1  
    - Once we have the [.xml and .bin files](TRAFFIC_LIGHTS/TL_W_COLORS/new_models/TL_v5/OPENVINO)
    - Download the last version (2020.1) Openvino Raspberry, following this [steps](https://www.pyimagesearch.com/2019/04/08/openvino-opencv-and-movidius-ncs-on-the-raspberry-pi/)
    - Use this [prediction script](TRAFFIC_LIGHTS/Prediction/assync_prediction.py) to predict with multiprocessing with NCS conected, (BRG2RGB must be added!)
    
5. Modify [prediction script](TRAFFIC_LIGHTS/Prediction/assync_prediction.py) to send the detection messages throw MQTT to the Nodered App for the Scooter implementation
`python3 async_prediction.py -m ../TL_W_COLORS/new_models/TL_v5/OPENVINO/IR7,FP16/frozen_inference_graph.xml -i video.mp4 -d MIRIAD` OTHER `[-pt PROB_THRESHOLD] [--no_show] [--labels LABELS]`

## References:

- xml_to_csv.py and generate_tfrecord.py from: https://github.com/datitran/raccoon_dataset
- bosch_to_pascal.py from: https://github.com/bosch-ros-pkg/bstld
- TensorFlow Object Dtection API from: https://github.com/tensorflow/models
- PyImagesearch: https://www.pyimagesearch.com/
- OpenCV latest version
