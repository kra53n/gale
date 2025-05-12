'''
The Python wrap around `gale.h`.

This script demonstrates the core design principle behind
integrating Python with the C library gale.h

Design Overview:

The central idea is to leverage an existing, well-designed
C library (gale.h) that provides clean, thoughtful function
interfaces. These functions can then be accessed from any
programming language capable of interfacing with C -- whether
by loading a shared library (DLL) or interpreting C code directly.

By adopting this approach, we gain several key advantages:

- Performance: Direct access to optimized C code ensures minimal overhead.
- Interoperability: The design allows seamless integration with a wide range
of programming languages.
- Maintainability: The core logic remains in one place (the C library),
simplifying maintenance and reducing duplication.

This Python wrapper serves as a practical example of how such integration can be
achieved, demonstrating both the simplicity and power of using C as a
"universal interface" for high-performance components.
'''

import ctypes as c
import pathlib
import platform
import typing
import subprocess
import shutil


__all__ = ['ImgFormat', 'Img']


_gale_c = 'gale.c'
_gale_h = 'gale.h'
_system_name = platform.system()

_dll_name = {
    'Windows': 'gale.dll',
    'Linux': 'libgale.so',
}.get(_system_name)
if _dll_name is None:
    print("You are not under Windows and Linux and it is unexpected behaviour.")
    exit(1)

_dll_path = pathlib.Path(_dll_name).absolute()
if not _dll_path.exists():
    shutil.copy(_gale_h, _gale_c)
    _compilation_cmd = f'gcc {_gale_c} -o {_dll_name} -shared'
    if _system_name == 'Linux':
        _compilation_cmd += ' -fpic'
    _res = subprocess.run(_compilation_cmd.split())
    if _res.returncode != 0:
        print(f'Was running `{" ".join(_res.args)}` and got error with status code {_res.returncode}')
        exit(1)
    pathlib.Path(_gale_c).unlink()


_dll = c.CDLL(_dll_path.absolute())


_ImgSize = c.c_int
_ImgData = c.POINTER(c.c_ubyte)
_ImgComp = c.c_int
_ImgFormat = c.c_int
_Filename = c.c_char_p


class ImgFormat:
    NotSupported = 0
    JPG = 1
    PNG = 2
    BMP = 3
    TGA = 4
    PSD = 5
    HDR = 6
    GIF = 7
    PNM = 8


class _ImgCStruct(c.Structure):
    _fields_ = [
        ("w", _ImgSize),
        ("h", _ImgSize),
        ("d", _ImgData),
        ("c", _ImgComp),
        ("f", _ImgFormat)
    ]


_gale_load_img = getattr(_dll, 'gale_load_img')
_gale_load_img.restype = _ImgCStruct
_gale_load_img.argtypes = (_Filename, )

_gale_free_img = getattr(_dll, 'gale_free_img')
_gale_free_img.argtypes = (_ImgCStruct, )

_gale_save_img = getattr(_dll, 'gale_save_img')
_gale_save_img.restype = c.c_int
_gale_save_img.argtypes = (_ImgCStruct, _Filename)

_gale_save_img_as = getattr(_dll, 'gale_save_img_as')
_gale_save_img_as.restype = c.c_int
_gale_save_img_as.argtypes = (_ImgCStruct, _Filename, _ImgFormat)


class Img:
    def __init__(self, path: typing.Union[pathlib.Path, str]):
        self.p: str = path
        if isinstance(self.p, pathlib.Path):
            self.p = str(self.p.absolute())
        self.img: typing.Optional[ImgCStruct] = None

    def __enter__(self):
        self.img = _gale_load_img(self.p.encode())
        return self

    def __exit__(self, type, val, traceback):
        _gale_free_img(self.img)

    def save(self, filename: typing.Union[pathlib.Path, str]):
        if isinstance(self.p, pathlib.Path):
            filename = str(self.p.absolute())
        res = _gale_save_img(self._mg, filename.encode())

    def save_as(self, filename: typing.Union[pathlib.Path, str], format: ImgFormat):
        if isinstance(self.p, pathlib.Path):
            filename = str(self.p.absolute())
        res = _gale_save_img_as(self.img, filename.encode(), format)

