# SmartScooter Object Detection files

## Create the model

To create the model, git clone this repository to SmartScooter1996 google account in Colab, and run the ipynb:
- [Cones](CONES/Create_object_detctor.ipynb)
- [TL](TRAFFIC_LIGHTS/TL_W_COLORS/Bosch_TL_DATASET.ipynb)

## Inference time
1. Once you have created the graph.pbtxt with the frozen_inference_graph.pb, REPLACE "AddV2" for "Add" in the graph.pbtxt. And remember to work with the latest version of opencv (4.2)
2. Use prediction files:
    - [Cones](CONES/Prediction)
    - [TL](TRAFFIC_LIGHTS/Prediction)
    - pred.py for photo inference and vide_object_detection.py for video inference

## References:

- xml_to_csv.py and generate_tfrecord.py from: https://github.com/datitran/raccoon_dataset
- bosch_to_pascal.py from: https://github.com/bosch-ros-pkg/bstld
- TensorFlow Object Dtection API from: https://github.com/tensorflow/models
- OpenCV latest version
