#!/usr/bin/env python
"""
This script Converts Yaml annotations to Pascal .xml Files
of the Bosch Small Traffic Lights Dataset.
Example usage:
    python bosch_to_pascal.py input_yaml out_folder
"""

import os
import sys
import yaml
from lxml import etree
import os.path
import xml.etree.cElementTree as ET


def write_xml(savedir, image, imgWidth, imgHeight,
              depth=3, pose="Unspecified"):
	boxes = image['boxes']
	impath = image['path']
#-------------NEW
	onlypath=impath.split('./')[1]
	second_folder=onlypath.split('/')[2]
	imagename = impath.split('/')[-1]
	path_w_noimage=onlypath.split(imagename)[0]
	currentfolder = savedir.split("\\")[-1]
	if 'train' in onlypath:		
		if second_folder not in os.listdir("/content/drive/My Drive/SMART_SCOOTER/images/rgb_jpg/train"): 
			os.mkdir("/content/drive/My Drive/SMART_SCOOTER/images/rgb_jpg/train/"+second_folder)
	if imagename in os.listdir("/content/drive/My Drive/SMART_SCOOTER/images/"+path_w_noimage):
#--------------
		annotation = ET.Element("annotaion")
		ET.SubElement(annotation, 'folder').text = str(currentfolder)
		ET.SubElement(annotation, 'filename').text = str(onlypath) #onlypath
		imagename = imagename.split('.')[0]
		size = ET.SubElement(annotation, 'size')
		ET.SubElement(size, 'width').text = str(imgWidth)
		ET.SubElement(size, 'height').text = str(imgHeight)
		ET.SubElement(size, 'depth').text = str(depth)
		ET.SubElement(annotation, 'segmented').text = '0'
		for box in boxes:
			obj = ET.SubElement(annotation, 'object')
			ET.SubElement(obj, 'name').text = str(box['label'])
			ET.SubElement(obj, 'pose').text = str(pose)
			ET.SubElement(obj, 'occluded').text = str(box['occluded'])
			ET.SubElement(obj, 'difficult').text = '0'

			bbox = ET.SubElement(obj, 'bndbox')
			ET.SubElement(bbox, 'xmin').text = str(box['x_min'])
			ET.SubElement(bbox, 'ymin').text = str(box['y_min'])
			ET.SubElement(bbox, 'xmax').text = str(box['x_max'])
			ET.SubElement(bbox, 'ymax').text = str(box['y_max'])

		xml_str = ET.tostring(annotation)
		root = etree.fromstring(xml_str)
		xml_str = etree.tostring(root, pretty_print=True)
		save_path = os.path.join(savedir, imagename + ".xml")
		with open(save_path, 'wb') as temp_xml:
			temp_xml.write(xml_str)


if __name__ == '__main__':
    if len(sys.argv) < 3:
        print(__doc__)
        sys.exit(-1)
    yaml_path = sys.argv[1]
    out_dir = sys.argv[2]
    images = yaml.load(open(yaml_path, 'rb').read())

    for image in images:
        write_xml(out_dir, image, 1280, 720, depth=3, pose="Unspecified")
