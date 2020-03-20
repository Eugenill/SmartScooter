ADDING NECESSARY REPOSITORIES:

1. ADD the TensorFlow Object Dtection API with: git clone https://github.com/tensorflow/models.git
2. ADD this repository with: git clone https://github.com/Eugenill/SmartScooter.git

STEPS TO CREATE AN OBJECT DETECTOR:

1. Import Images of your dataset to models/research/object_detection/images/"name of the datset"/train and /test
2. Import xml (or yaml) annotations of your dataset in models/research/object_detection/annotations/"name of the dataset"
3. Convert to csv annotations:
    1. From xml: run xml_to_csv.py in CONES or TRAFFIC_LIGHT folders
    2. From yaml: run bosch_to_pascal.py to convert them into xml, and xml_to_csv.py to csv in TL_W_COLORS
4. Generate test.record and train.record running generate_tf_record.py. Will be found in CONES or TRAFFIC_LIGHT/SIMPLE or W_COLORS
5. Edit "model name".config, in config folders in this repo
6. Edit labelmap.pbtxt. Will be found in CONES/labelmap or TRAFFIC_LIGHT/SIMPLE/labelmap or W_COLORS/labelmap
7. Run the training with models/research/object_detection/model_main.py

Follow the instruccions on the *ipynb* on the folders to create the object detector correctly

References:

- xml_to_csv.py and generate_tfrecord.py from: https://github.com/datitran/raccoon_dataset
- bosch_to_pascal.py from: https://github.com/bosch-ros-pkg/bstld
