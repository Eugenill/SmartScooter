# Import and initialize the pygame library
import os

import cv2 as cv
import pygame

model_path = 'RaspberryPi/Others/inference_graph/'
cvNet = cv.dnn.readNetFromTensorflow(model_path + 'frozen_inference_graph.pb', model_path + 'graph.pbtxt')
images_path = 'RaspberryPi/Others/images/'
images = os.listdir(images_path)
print(images)
colors_list = []
vertex_list = []
class_names = {}

pygame.init()

# Set up the drawing window
screen = pygame.display.set_mode([1300, 1200])

font = pygame.font.SysFont("comicsansms", 20)


def main():
    setup()
    print(vertex_list)
    print(colors_list)
    print(class_names)
    running = True

    # Run until the user asks to quit
    for next_img in images:

        # Did the user click the window close button?
        for event in pygame.event.get():

            if event.type == pygame.QUIT:
                running = False

        # Fill the background with white
        screen.fill((255, 255, 255))

        show_and_wait(next_img)

        # Print rects and text
        print_rects_and_text()

        # Loop over folder images
        cv_img = next_image(next_img)
        rows = cv_img.shape[0]
        cols = cv_img.shape[1]
        cv_out = detect(cv_img)

        for detection in cv_out[0, 0, :, :]:
            score = float(detection[2])
            if score > 0.3:
                left = detection[3] * cols
                top = detection[4] * rows
                right = detection[5] * cols
                bottom = detection[6] * rows
                rect = (int(left), int(top), int(right), int(bottom))
                print(rect)
                # pygame.draw.rect(screen, (23, 230, 210), rect, width=2)
        # cv.waitKey()
        # Flip the display
        pygame.display.flip()

        if not running:
            pygame.quit()

    # Done! Time to quit.
    pygame.quit()


def define_color(total_of_classes):
    for n in range(total_of_classes):
        colors_list.append((n * 5, n * 15, 255 - n * 15))


def rect_vertexes(total_of_classes):
    next_top_left_vertex = (150, 800)
    vertex_list.append(next_top_left_vertex)
    for i in range(total_of_classes):

        if next_top_left_vertex[0] >= 650:
            next_top_left_vertex = (200, next_top_left_vertex[1] + 50)

        else:
            next_top_left_vertex = (next_top_left_vertex[0] + 200, next_top_left_vertex[1])

        vertex_list.append(next_top_left_vertex)


def print_rects_and_text():
    values = list(class_names.values())
    for i in range(len(class_names)):
        text = font.render(values[i], True, colors_list[i])
        text_rect = text.get_rect()
        text_rect.center = (vertex_list[i][0] + 50, vertex_list[i][1] + 20)
        screen.blit(text, text_rect)


def setup():
    # Add buttons for every class
    names = {
        "R-2": "STOP",
        "R-1": "CEDA",
        "R-101": "ENTRADA PROHIBIDA",
        "R-301": "VELOCIDAD MÁXIMA",
        "R-307": "PARA Y ESTAC. PROHIBIDO",
        "R-308": "ESTAC. PROHIBIDO",
        "S-13": "CEDA",
        "S-18": "TAXIS",
        "S-19": "BUS STOP",
        "S-28": "CALLE RESIDENCIAL",
        "S-30": "ZONA 30",
        "green": "SEMÁFORO VERDE",
        "yellow": "SEMÁFORO ÁMBAR",
        "red": "SEMÁFORO ROJO",
        "off": "SEMÁFORO OFF",
        "P-18": "OBRAS",
        "P-21": "NIÑOS"
    }
    for key in names:
        class_names[key] = names[key]

    define_color(len(class_names))
    rect_vertexes(len(class_names))


def detect(img):
    print("image readed")
    cvNet.setInput(cv.dnn.blobFromImage(img, size=(650, 400), swapRB=True, crop=False))
    return cvNet.forward()


def next_image(next_img):
    return cv.imread(images_path + next_img)


def show_and_wait(next_img):
    img = pygame.image.load(images_path + next_img)
    screen.blit(img, (0, 0))


main()
