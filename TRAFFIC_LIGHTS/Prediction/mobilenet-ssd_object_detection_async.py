"""
 Copyright (c) 2018 Intel Corporation

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

      http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
"""

# Most of these codes are based on 
# OpenVINO python sample code "object_detection_demo_ssd_async.py" 
# https://software.intel.com/en-us/openvino-toolkit/
# and
# Object detection with deep learning and OpenCV
# https://www.pyimagesearch.com/2017/09/11/object-detection-with-deep-learning-and-opencv/
# Example Usage:
# MYRIAD: python mobilenet-ssd_object_detection_async.py -i cam -m IR\MobileNetSSD_FP16\MobileNetSSD_deploy.xml -d MYRIAD
# CPU:    python mobilenet-ssd_object_detection_async.py -i cam -m IR\MobileNetSSD_FP32\MobileNetSSD_deploy.xml -l Intel\OpenVINO\inference_engine_samples_2017\intel64\Release\cpu_extension.dll

# import the necessary packages
import numpy as np
from argparse import ArgumentParser
import cv2
import time
import sys
import os
import logging as log
from timeit import default_timer as timer
from openvino.inference_engine import IENetwork, IEPlugin

# initialize the list of class labels MobileNet SSD was trained to
# detect, then generate a set of bounding box colors for each class
CLASSES = [
    "background", "aeroplane", "bicycle", "bird", "boat", "bottle", "bus",
    "car", "cat", "chair", "cow", "diningtable", "dog", "horse", "motorbike",
    "person", "pottedplant", "sheep", "sofa", "train", "tvmonitor"
]
COLORS = np.random.uniform(0, 255, size=(len(CLASSES), 3))


# construct the argument parse and parse the arguments
def build_argparser():
    parser = ArgumentParser()
    parser.add_argument(
        "-m",
        "--model",
        help="Path to an .xml file with a trained model.",
        required=True,
        type=str)
    parser.add_argument(
        "-i",
        "--input",
        help="Path to video file or image. 'cam' for capturing video stream from camera",
        required=True,
        type=str)
    parser.add_argument(
        "-l",
        "--cpu_extension",
        help="MKLDNN (CPU)-targeted custom layers.Absolute path to a shared library with the kernels "
        "impl.",
        type=str,
        default=None)
    parser.add_argument(
        "-pp",
        "--plugin_dir",
        help="Path to a plugin folder",
        type=str,
        default=None)
    parser.add_argument(
        "-d",
        "--device",
        help="Specify the target device to infer on; CPU, GPU, FPGA or MYRIAD is acceptable. Demo "
        "will look for a suitable plugin for device specified (CPU by default)",
        default="CPU",
        type=str)
    parser.add_argument(
        "-pt",
        "--prob_threshold",
        help="Probability threshold for detections filtering",
        default=0.2,
        type=float)

    return parser


def main():
    log.basicConfig(
        level=log.INFO,
        format="%(asctime)s %(levelname)s %(name)s %(funcName)s(): %(message)s",
        stream=sys.stdout)
    args = build_argparser().parse_args()
    model_xml = args.model
    model_bin = os.path.splitext(model_xml)[0] + ".bin"
    # Plugin initialization for specified device and load extensions library if specified
    log.info("Initializing plugin for {} device...".format(args.device))
    plugin = IEPlugin(device=args.device, plugin_dirs=args.plugin_dir)
    if args.cpu_extension and 'CPU' in args.device:
        plugin.add_cpu_extension(args.cpu_extension)
    # Read IR
    log.info("Reading IR...")
    net = IENetwork(model=model_xml, weights=model_bin)

    if plugin.device == "CPU":
        supported_layers = plugin.get_supported_layers(net)
        not_supported_layers = [
            l for l in net.layers.keys() if l not in supported_layers
        ]
        if len(not_supported_layers) != 0:
            log.error(
                "Following layers are not supported by the plugin for specified device {}:\n {}".
                format(plugin.device, ', '.join(not_supported_layers)))
            log.error(
                "Please try to specify cpu extensions library path in demo's command line parameters using -l "
                "or --cpu_extension command line argument")
            sys.exit(1)
    assert len(
        net.inputs.keys()) == 1, "Demo supports only single input topologies"
    assert len(net.outputs) == 1, "Demo supports only single output topologies"
    input_blob = next(iter(net.inputs))
    out_blob = next(iter(net.outputs))
    log.info("Loading IR to the plugin...")
    exec_net = plugin.load(network=net, num_requests=2)
    # Read and pre-process input image
    n, c, h, w = net.inputs[input_blob].shape
    log.info(
        "net.inpute.shape(n, c, h, w):{}".format(net.inputs[input_blob].shape))
    del net
    if args.input == 'cam':
        input_stream = 0
    else:
        input_stream = args.input
        assert os.path.isfile(args.input), "Specified input file doesn't exist"

    cap = cv2.VideoCapture(input_stream)
    cur_request_id = 0
    next_request_id = 1

    log.info("Starting inference in async mode...")
    log.info("To switch between sync and async modes press Tab button")
    log.info("To stop the demo execution press Esc button")
    is_async_mode = True
    ##is_async_mode = False
    render_time = 0
    ret, frame = cap.read()
    initial_w = cap.get(3)
    initial_h = cap.get(4)
    resize_w = 640
    resize_h = 480

    ##
    accum_time = 0
    curr_fps = 0
    fps = "FPS: ??"
    prev_time = timer()

    while cap.isOpened():
        if is_async_mode:
            ret, next_frame = cap.read()
            next_frame = cv2.flip(next_frame, 1)
        else:
            ret, frame = cap.read()
            frame = cv2.flip(frame, 1)
        if not ret:
            break
        # Main sync point:
        # in the truly Async mode we start the NEXT infer request, while waiting for the CURRENT to complete
        # in the regular mode we start the CURRENT request and immediately wait for it's completion
        inf_start = timer()
        if is_async_mode:
            in_frame = cv2.dnn.blobFromImage(
                cv2.resize(next_frame, (300, 300)), 0.007843, (300, 300),
                127.5)
            exec_net.start_async(
                request_id=next_request_id, inputs={input_blob: in_frame})
        else:
            in_frame = cv2.dnn.blobFromImage(
                cv2.resize(frame, (300, 300)), 0.007843, (300, 300), 127.5)
            ## start async request if CURRENT request has been completed.
            if exec_net.requests[cur_request_id].wait(-1) == 0:
                exec_net.start_async(
                    request_id=cur_request_id, inputs={input_blob: in_frame})
            log.info("in_frame shape:{} cur_req_id:{} next_req_id:{}".format(
                in_frame.shape, cur_request_id, next_request_id))
        if exec_net.requests[cur_request_id].wait(-1) == 0:
            inf_end = timer()
            det_time = inf_end - inf_start

            # Parse detection results of the current request
            log.debug("computing object detections...")
            detections = exec_net.requests[cur_request_id].outputs[out_blob]
            log.debug("detections shape:{}".format(detections.shape))
            
            # ref: https://www.pyimagesearch.com/2017/09/11/object-detection-with-deep-learning-and-opencv/
            # loop over the detections
            for i in np.arange(0, detections.shape[2]):
                # extract the confidence (i.e., probability) associated with the
                # prediction
                confidence = detections[0, 0, i, 2]
                # filter out weak detections by ensuring the `confidence` is
                # greater than the minimum confidence
                if confidence > args.prob_threshold:
                    # extract the index of the class label from the `detections`,
                    # then compute the (x, y)-coordinates of the bounding box for
                    # the object
                    idx = int(detections[0, 0, i, 1])
                    box = detections[0, 0, i, 3:7] * np.array(
                        [initial_w, initial_h, initial_w, initial_h])
                    (startX, startY, endX, endY) = box.astype("int")
                    log.debug("startX, startY, endX, endY: {}".format(
                        box.astype("int")))

                    # display the prediction
                    label = "{}: {:.2f}%".format(CLASSES[idx],
                                                 confidence * 100)
                    log.info("{} {}".format(cur_request_id, label))
                    cv2.rectangle(frame, (startX, startY), (endX, endY),
                                  COLORS[idx], 2)
                    y = startY - 15 if startY - 15 > 15 else startY + 15
                    cv2.putText(frame, label, (startX, y),
                                cv2.FONT_HERSHEY_SIMPLEX, 0.5, COLORS[idx], 2)

            # Draw performance stats
            inf_time_message = "Inference time: N\A for async mode" if is_async_mode else \
                "Inference time: {:.3f} ms".format(det_time * 1000)
            render_time_message = "OpenCV rendering time: {:.3f} ms".format(
                render_time * 1000)
            async_mode_message = "Async mode is on. Processing request {}".format(cur_request_id) if is_async_mode else \
                "Async mode is off. Processing request {}".format(cur_request_id)

            ##
            frame = cv2.resize(frame, (resize_w, resize_h))
            cv2.putText(frame, inf_time_message, (15, 15),
                        cv2.FONT_HERSHEY_COMPLEX, 0.5, (200, 10, 10), 1)
            cv2.putText(frame, render_time_message, (15, 30),
                        cv2.FONT_HERSHEY_COMPLEX, 0.5, (10, 10, 200), 1)
            cv2.putText(frame, async_mode_message, (10, int(initial_h - 20)),
                        cv2.FONT_HERSHEY_COMPLEX, 0.5, (10, 10, 200), 1)

        ## ref. https://github.com/rykov8/ssd_keras/blob/master/testing_utils/videotest.py
        # Calculate FPS
        # This computes FPS for everything, not just the model's execution
        # which may or may not be what you want
        curr_time = timer()
        exec_time = curr_time - prev_time
        prev_time = curr_time
        accum_time = accum_time + exec_time
        curr_fps = curr_fps + 1
        if accum_time > 1:
            accum_time = accum_time - 1
            fps = "FPS: " + str(curr_fps)
            curr_fps = 0

        # Draw FPS in top left corner
        cv2.rectangle(frame, (resize_w - 50, 0), (resize_w, 17),
                      (255, 255, 255), -1)
        cv2.putText(frame, fps, (resize_w - 50 + 3, 10),
                    cv2.FONT_HERSHEY_SIMPLEX, 0.35, (0, 0, 0), 1)

        # 
        render_start = timer()
        cv2.imshow("Detection Results(MobileNetSSD)", frame)
        render_end = timer()
        render_time = render_end - render_start

        if is_async_mode:
            cur_request_id, next_request_id = next_request_id, cur_request_id
            frame = next_frame

        key = cv2.waitKey(1)
        if key == 27:
            break
        if (9 == key):
            is_async_mode = not is_async_mode
            log.info("Switched to {} mode".format("async" if is_async_mode else
                                                  "sync"))

    cv2.destroyAllWindows()
    del exec_net
    del plugin


if __name__ == '__main__':
    sys.exit(main() or 0)