#ifndef GALE_INCLUDE_H
#define GALE_INCLUDE_H

#define STB_IMAGE_IMPLEMENTATION
#include "stb_image.h"
#define STB_IMAGE_WRITE_IMPLEMENTATION
#include "stb_image_write.h"

typedef unsigned int gale_ImgSize;
typedef unsigned char *  gale_ImgData;
typedef const char *  gale_Filename;
typedef int gale_ImgComp;
typedef char *  gale_Err;

#define gale_PREFIX_ERROR "ERROR  "

char gale_ERR_BUF[256] = {0};

typedef enum {
    gale_ImgFormat_NotSupported,
    gale_ImgFormat_JPG,
    gale_ImgFormat_PNG,
    gale_ImgFormat_BMP,
    gale_ImgFormat_TGA,
    gale_ImgFormat_PSD,
    gale_ImgFormat_HDR,
    gale_ImgFormat_GIF,
    gale_ImgFormat_PNM,
} gale_ImgFormat;

typedef struct {
    gale_ImgSize w;
    gale_ImgSize h;
    gale_ImgData d;
    gale_ImgComp c;
    gale_ImgFormat f;
    gale_Err err;
} gale_Img;

gale_ImgFormat gale__define_img_format(stbi__context *s) {
    if (stbi__jpeg_test(s)) return gale_ImgFormat_JPG;
    if (stbi__png_test(s)) return gale_ImgFormat_PNG;
    if (stbi__bmp_test(s)) return gale_ImgFormat_BMP;
    if (stbi__tga_test(s)) return gale_ImgFormat_TGA;
    if (stbi__psd_test(s)) return gale_ImgFormat_PSD;
    if (stbi__hdr_test(s)) return gale_ImgFormat_HDR;
    if (stbi__gif_test(s)) return gale_ImgFormat_GIF;
    if (stbi__pnm_test(s)) return gale_ImgFormat_PNM;
    return gale_ImgFormat_NotSupported;
}

gale_Img gale_load_img(gale_Filename filename) {
    FILE *f;
    gale_Img i = {0};
    int req_comp = 0;
    stbi__context s;

    f = stbi__fopen(filename, "rb");
    if (!f) {
        i.err = gale_ERR_BUF;
        strcat(i.err, gale_PREFIX_ERROR"no file: ");
        strcat(i.err, filename);
        strcat(i.err, "\n");
        return i;
    }

    stbi__start_file(&s, f);
    i.f = gale__define_img_format(&s);
    i.d = stbi__load_and_postprocess_8bit(&s, &i.w, &i.h, &i.c, req_comp);
    if (i.d) {
        fseek(f, - (int) (s.img_buffer_end - s.img_buffer), SEEK_CUR);
    }

    return i;
}

void gale_free_img(gale_Img i) {
    STBI_FREE(i.d);
}

int gale_save_img_as(gale_Img i, gale_Filename filename, gale_ImgFormat f) {
    int stride = i.w * i.c;
    int quality = 0;
    int comp = 0;
    switch (f) {
    case gale_ImgFormat_JPG: return stbi_write_jpg(filename, i.w, i.h, i.c, i.d, quality);
    case gale_ImgFormat_PNG: return stbi_write_png(filename, i.w, i.h, i.c, i.d, stride);
    case gale_ImgFormat_BMP: return stbi_write_bmp(filename, i.w, i.h, i.c, i.d);
    case gale_ImgFormat_TGA: return stbi_write_tga(filename, i.w, i.h, i.c, i.d);
    /* case gale_ImgFormat_HDR: return stbi_write_hdr(filename, i.w, i.h, i.c, i.d); */
    case gale_ImgFormat_PSD: return 0; // NOT SUPPORTED
    case gale_ImgFormat_GIF: return 0; // NOT SUPPORTED
    case gale_ImgFormat_PNM: return 0; // NOT SUPPORTED
    }
    return 0;
}

int gale_save_img(gale_Img i, gale_Filename f) {
    return gale_save_img_as(i, f, i.f);
}

void gale_crop_img(gale_Img *i, int x1, int y1, int x2, int y2) {
    int new_w = x2 - x1;
    int new_h = y2 - y1;

    int offset_y = i->w * i->c * y1;
    int shifted_cursor = offset_y + x1 * i->c;
    int cursor = 0;
    for (int y = 0; y < new_h; y++) {
        int src_offset = ((y1 + y) * i->w + x1) * i->c;
        int dst_offset = y * new_w * i->c;
        memcpy(&i->d[dst_offset], &i->d[src_offset], new_w * i->c);
    }
    i->w = new_w;
    i->h = new_h;
}


#endif // GALE_INCLUDE_H


#ifdef GALE_DEBUG

const char *  gale__img_format_to_str[9] = {
    "NotSupported",
    "JPG",
    "PNG",
    "BMP",
    "TGA",
    "PSD",
    "HDR",
    "GIF",
    "PNM",
};

#ifndef GALE_DEBUG_PREFIX
#define GALE_DEBUG_PREFIX "| "
#endif

#ifndef GALE_DEBUG_DELIMITER
#define GALE_DEBUG_DELIMITER " - "
#endif

void gale_print_img_info(gale_Img i) {
    printf(GALE_DEBUG_PREFIX"w"GALE_DEBUG_DELIMITER"%d\n", i.w);
    printf(GALE_DEBUG_PREFIX"h"GALE_DEBUG_DELIMITER"%d\n", i.h);
    printf(GALE_DEBUG_PREFIX"c"GALE_DEBUG_DELIMITER"%d\n", i.c);
    printf(GALE_DEBUG_PREFIX"f"GALE_DEBUG_DELIMITER"%s\n", gale__img_format_to_str[i.f]);
}

#endif // GALE_DEBUG
