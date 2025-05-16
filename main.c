#include <stdio.h>
#include "gale.h"

int main() {
    gale_Filename img_name = "sword_art_online.jpg";
    gale_Img img = gale_load_img(img_name);
	int x1 = 0;
	int y1 = 0;
	int x2 = img.w;
	int y2 = img.h;
    gale_Err err = gale_resize_img(&img, x1, y1, x2 * 0.5, y2);
    if (err) {
        printf("Error occured while resizing the ");
        return 1;
    }
    printf(":: w: %d", img.w);
    gale_save_img_as(img, "sao.bmp", gale_ImgFormat_BMP);
    return 0;
}
