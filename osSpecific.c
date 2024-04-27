#ifdef LINUX
#include <X11/X.h>

#endif
#ifdef WINDOWS
#include <windows.h>
#include <winuser.h>

#endif

struct mousePos{
	int x;
	int y;
};

int init();

struct mousePos getMousePos();

int setMousePos(struct mousePos cordinates);

int createWindow();

int setWindowText();

int deleteWindow();
