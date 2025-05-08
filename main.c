#include <stdio.h>

#define GALE_DEBUG
#include "gale.h"

int main() {
    const char *  filename_to_load = "sword_art_online.jpg";
    const char *  filename_to_save = "sword_art_online.png";
    gale_Img img = gale_load_img(filename_to_load);
    gale_print_img_info(img);
    gale_save_img(img, filename_to_save);
}
