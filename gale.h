#ifndef GALE_INCLUDE_H
#define GALE_INCLUDE_H

#define STB_IMAGE_IMPLEMENTATION
#include "stb_image.h"
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h"

typedef unsigned int gale_ImgSize;
typedef unsigned char *  gale_ImgData;
typedef const char *  gale_Filename;
typedef int gale_Channels;

typedef struct {
    gale_ImgSize w;
    gale_ImgSize h;
    gale_ImgData d;
    gale_Channels c;
} gale_Img;

gale_Img gale_load_img(gale_Filename f) {
    gale_Img i;
    i.d = stbi_load(f, &i.w, &i.h, &i.c, 0);
    return i;
}

gale_Img gale_save_img(gale_Img i, gale_Filename f) {
    int stride = i.w * i.c;
    int comp = 0;
    stbi_write_png(f, i.w, i.h, i.c, i.d, stride);
}


#endif // GALE_INCLUDE_H


#ifdef GALE_DEBUG

#include <stdio.h>

void gale_print_img_info(gale_Img i) {
    printf(":: w: %d\n", i.w);
    printf(":: h: %d\n", i.h);
    printf(":: c: %d\n", i.c);
}

#endif // GALE_DEBUG
