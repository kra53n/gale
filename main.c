#include <stdio.h>
#include "gale.h"

int main() {
    gale_Filename img_name = "sword_art_online.jpg";
    gale_Img img = gale_load_img(img_name);
    if (img.err) {
        printf("%s\n", img.err);
        return 1;
    }
    int x1 = img.w * 0.2;
    int x2 = img.w * 0.8;
    int y1 = img.h * 0.2;
    int y2 = img.h * 0.8;
    gale_crop_img(&img, x1, y1, x2, y2);
    if (img.err) {
        printf("%s\n", img.err);
        return 1;
    }
    /* gale_rotate_img(&img); */
    gale_save_img_as(&img, "sao.bmp", gale_ImgFormat_BMP);
    if (img.err) {
        printf("%s\n", img.err);
        return 1;
    }
    return 0;
}
