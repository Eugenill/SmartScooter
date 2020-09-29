# Import and initialize the pygame library
import os
import cv2 as cv
import pandas as pd
import pygame


model_path = 'traffic_sign_model/inference_graph/'
cvNet = cv.dnn.readNetFromTensorflow(model_path + 'frozen_inference_graph.pb', model_path + 'graph.pbtxt')
images_base = 'imatges/'
images_folder = 'ts_frames'
images_path = images_base + images_folder+'/'
csv_name = images_folder
images = os.listdir(images_path)
colors_list = []
vertex_list = []
names = []
class_names = {
    "R-2": "STOP",
    "R-1": "CEDA",
    "R-101": "ENTRADA PROHIBIDA",
    "R-301": "VELOCIDAD MÁXIMA",
    "R-307": "PARA Y ESTAC. PROHIBIDO",
    "R-308": "ESTAC. PROHIBIDO",
    "S-13": "PASO PEATÓN",
    "S-18": "TAXIS",
    "S-19": "BUS STOP",
    "S-28": "CALLE RESIDENCIAL",
    "S-30": "ZONA 30",
    "green": "SEMÁFORO VERDE",
    "yellow": "SEMÁFORO ÁMBAR",
    "red": "SEMÁFORO ROJO",
    "off": "SEMÁFORO OFF",
    "P-18": "OBRAS",
    "P-21": "NIÑOS",
    "other":"OTHER"
}

data = []
column_name = ['filename', 'width', 'height', 'class', 'xmin', 'ymin', 'xmax', 'ymax']

pygame.init()
# Set up the drawing window
screen = pygame.display.set_mode([1300, 1200])
font = pygame.font.SysFont("comicsansms", 15)
font2 = pygame.font.SysFont("comicsansms", 40)


def detection_and_paint():
    for next_img in images:
        if next_img.split(".")[1] == "jpg":
            print("next image")
            # Fill the background with white
            screen.fill((0, 0, 0))

            cv_img = next_image(next_img)
            rows = cv_img.shape[0]
            cols = cv_img.shape[1]
            cv_out = detect(cv_img)
            i = 0
            for detection in cv_out[0, 0, :, :]:
                score = float(detection[2])
                if score > 0.1:
                    left = detection[3] * cols
                    top = detection[4] * rows
                    right = detection[5] * cols
                    bottom = detection[6] * rows
                    ext = False
                    while True:
                        display_img(next_img)
                        # Print rects and text
                        print_rects_and_text()
                        rect = pygame.Rect(int(left), int(top), int(right) - int(left), int(bottom) - int(top))
                        pygame.draw.rect(screen, (23, 230, 210), rect, 2)
                        text = font2.render(str(i) + "/" + str(len(cv_out[0, 0, :, :])), True, (0, 0, 0))
                        text_rect = text.get_rect()
                        text_rect.center = (1000, 20)
                        screen.blit(text, text_rect)

                        pos = pygame.mouse.get_pos()
                        events = pygame.event.get()
                        # print(events)
                        for event in events:
                            if event and event.type == pygame.QUIT:
                                pygame.quit()
                            if event and event.type == pygame.MOUSEBUTTONDOWN:
                                for enum_vertex in list(enumerate(vertex_list)):
                                    index = enum_vertex[0]
                                    point = enum_vertex[1]
                                    if point[0] < pos[0] < point[0] + 150 and point[1] < pos[1] < point[1] + 80:
                                        print(names[index])
                                        pygame.draw.rect(screen, (23, 230, 210), rect)
                                        data.append(
                                            [next_img, str(cols), str(rows), names[index], str(left), str(top), str(right),
                                             str(bottom)])
                                        ext = True
                                        break
                                if ext:
                                    break
                        if ext:
                            print("next detection")
                            break
                        pygame.display.flip()


def define_color(total_of_classes):
    for n in range(total_of_classes):
        colors_list.append((n * 5, n * 15, 255 - n * 15))


def rect_vertexes(total_of_classes):
    next_top_left_vertex = (0, 800)
    vertex_list.append(next_top_left_vertex)
    for i in range(total_of_classes):

        if next_top_left_vertex[0] >= 950:
            next_top_left_vertex = (0, next_top_left_vertex[1] + 40)

        else:
            next_top_left_vertex = (next_top_left_vertex[0] + 120, next_top_left_vertex[1])

        vertex_list.append(next_top_left_vertex)


def print_rects_and_text():
    values = list(class_names.values())
    for i in range(len(class_names)):
        rect = pygame.Rect(vertex_list[i][0], vertex_list[i][1], 120, 40)
        pygame.draw.rect(screen, colors_list[i], rect)
        text = font.render(values[i], True, (255, 255, 255), colors_list[i])
        text_rect = text.get_rect()
        text_rect.center = (vertex_list[i][0] + 60, vertex_list[i][1] + 20)
        screen.blit(text, text_rect)


def setup():
    # Add buttons for every class
    for key in class_names:
        names.append(class_names[key])

    define_color(len(class_names))
    rect_vertexes(len(class_names))


def detect(img):
    print("image readed")
    cvNet.setInput(cv.dnn.blobFromImage(img, size=(650, 400), swapRB=True, crop=False))
    return cvNet.forward()


def next_image(next_img):
    return cv.imread(images_path + next_img)


def display_img(next_img):
    img = pygame.image.load(images_path + next_img)
    screen.blit(img, (0, 0))


if __name__ == '__main__':
    setup()
    # print(vertex_list)
    # print(colors_list)
    # print(class_names)
    detection_and_paint()
    data = list(map(list, zip(*data)))
    print(data)
    dict = {}
    for i in range(len(column_name)):
        dict[column_name[i]] = data[i]
    annotations = pd.DataFrame(dict)
    annotations.to_csv(csv_name + '.csv')
